package web

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/olivere/elastic/v7"

	account "github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/files"
	internal "github.com/AguilaMike/Stori_Challenge_Go/internal/common/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Puedes añadir lógica de seguridad aquí
	},
}

// HomeHandler maneja la ruta raíz
func HomeHandler(accountService *account.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := accountService.ListAccounts(r.Context(), 100, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		RenderTemplate(w, "index.html", users)
	}
}

// UploadFileHandler maneja la carga de archivos de transacciones
func UploadFileHandler(fileUploadService *files.FileUploadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Verificar si el formulario está multiparte
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Printf("Error parsing multipart form: %v", err)
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("transactionFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error al obtener el archivo de la solicitud: %+v", err)
			return
		}
		defer file.Close()

		userID, err := uuid.Parse(r.FormValue("userID"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			log.Printf("Error al obtener el ID del usuario: %+v", err)
			return
		}

		err = fileUploadService.UploadTransactionFile(r.Context(), file, header, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error al cargar el archivo: %+v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded and queued for processing"))
	}
}

// WebSocketHandler maneja las conexiones WebSocket
func WebSocketHandler(wsService *internal.WebSocketService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userID")
		if userID == "" {
			http.Error(w, "Missing userID", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error upgrading to WebSocket: %v", err)
			return
		}

		wsService.AddClient(userID, conn)

		// Manejar la desconexión
		defer func() {
			wsService.RemoveClient(userID)
			conn.Close()
		}()

		// Mantener la conexión abierta
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func HealthCheckHandler(db *sql.DB, esClient *elastic.Client, natsConn *nats.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		status := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Services:  make(map[string]string),
		}

		// Verificar la conexión a la base de datos
		if err := db.PingContext(ctx); err != nil {
			status.Status = "unhealthy"
			status.Services["database"] = "unhealthy"
		} else {
			status.Services["database"] = "healthy"
		}

		// Verificar la conexión a Elasticsearch
		_, err := esClient.ClusterHealth().Do(ctx)
		if err != nil {
			status.Status = "unhealthy"
			status.Services["elasticsearch"] = "unhealthy"
		} else {
			status.Services["elasticsearch"] = "healthy"
		}

		// Verificar la conexión a NATS
		if natsConn.Status() != nats.CONNECTED {
			status.Status = "unhealthy"
			status.Services["nats"] = "unhealthy"
		} else {
			status.Services["nats"] = "healthy"
		}

		w.Header().Set("Content-Type", "application/json")
		if status.Status == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(status)
	}
}

// SetupWebRoutes configura todas las rutas web y devuelve un http.Handler
func SetupWebRoutes(
	accountService *account.AccountService,
	fileUploadService *files.FileUploadService,
	wsService *internal.WebSocketService,
	db *sql.DB,
	esClient *elastic.Client,
	natsConn *nats.Conn,
	templateDir string,
	staticDir string,
) (http.Handler, error) {
	// Inicializar plantillas
	err := InitTemplates(templateDir)
	if err != nil {
		return nil, err
	}

	// Crear un nuevo mux para manejar las rutas
	mux := http.NewServeMux()

	// Configurar rutas
	mux.HandleFunc("/", HomeHandler(accountService))
	mux.HandleFunc("/upload", UploadFileHandler(fileUploadService))
	mux.HandleFunc("/ws", WebSocketHandler(wsService))
	mux.HandleFunc("/health", HealthCheckHandler(db, esClient, natsConn))

	// Servir archivos estáticos
	staticHandler, handler := GetStaticHandler(staticDir)
	mux.Handle(staticHandler, handler)

	return mux, nil
}
