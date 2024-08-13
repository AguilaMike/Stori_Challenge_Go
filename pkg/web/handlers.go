package web

import (
	"net/http"

	"github.com/google/uuid"

	account "github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	transaction "github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
)

// HomeHandler maneja la ruta raíz
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UsersHandler maneja la lista de usuarios
func UsersHandler(accountService *account.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := accountService.ListAccounts(r.Context(), 100, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		RenderTemplate(w, "users.html", users)
	}
}

// CreateUserHandler maneja la creación de usuarios
func CreateUserHandler(accountService *account.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			nickname := r.FormValue("nickname")
			email := r.FormValue("email")

			_, err := accountService.CreateAccount(r.Context(), nickname, email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/users", http.StatusSeeOther)
			return
		}
		RenderTemplate(w, "create_user.html", nil)
	}
}

// UserDetailHandler maneja los detalles de un usuario
func UserDetailHandler(accountService *account.AccountService, transactionService *transaction.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := uuid.MustParse(r.URL.Query().Get("id"))

		user, err := accountService.GetAccount(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactions, err := transactionService.GetTransactionsByAccount(r.Context(), userID, 100, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			User         interface{}
			Transactions interface{}
		}{
			User:         user,
			Transactions: transactions,
		}
		RenderTemplate(w, "user_detail.html", data)
	}
}

// UploadFileHandler maneja la carga de archivos de transacciones
func UploadFileHandler(transactionService *transaction.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("transactionFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		//userID := uuid.MustParse(r.FormValue("userID"))

		// Aquí deberías implementar la lógica para procesar el archivo
		// y enviarlo al worker a través de NATS

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
	}
}

// SetupWebRoutes configura todas las rutas web y devuelve un http.Handler
func SetupWebRoutes(
	accountService *account.AccountService,
	transactionService *transaction.TransactionService,
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
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/users", UsersHandler(accountService))
	mux.HandleFunc("/users/create", CreateUserHandler(accountService))
	mux.HandleFunc("/users/detail", UserDetailHandler(accountService, transactionService))
	mux.HandleFunc("/upload", UploadFileHandler(transactionService))

	// Servir archivos estáticos
	staticHandler, handler := GetStaticHandler(staticDir)
	mux.Handle(staticHandler, handler)

	return mux, nil
}
