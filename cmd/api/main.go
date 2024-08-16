package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/config"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/db"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/elasticsearch"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/email"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/files"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/websocket"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction"
	"github.com/AguilaMike/Stori_Challenge_Go/pkg/api"
	"github.com/AguilaMike/Stori_Challenge_Go/pkg/web"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize pgDB connection
	pgDB, err := db.NewPostgresConnection(cfg.GetConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgDB.Close()

	// Ejecutar migraciones
	if err := db.RunMigrations(pgDB, "/root/migrations"); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize Elasticsearch client
	esClient, err := elasticsearch.NewElasticsearchClient(cfg.ElasticsearchURL)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Initialize NATS connection
	nc, err := nats.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Initialize services
	hostGrpc := fmt.Sprintf("%s:%s", cfg.DOMAIN, cfg.GRPCPort)
	connGrpc, err := grpc.NewClient(hostGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	accountService := account.SetupAccountDomain(pgDB, esClient, nc)
	transactionService := transaction.SetupTransactionDomain(pgDB, esClient, nc, connGrpc, emailSender)

	// Set up API HTTP router
	apiMux := api.SetupHTTPRoutes(accountService, transactionService)

	// Configurar rutas web
	templateDir := filepath.Join("web", "templates")
	staticDir := filepath.Join("web", "static")
	// Inicializar el servicio de carga de archivos
	fileUploadService, err := files.NewFileUploadService(nc)
	if err != nil {
		log.Fatalf("Failed to create file upload service: %v", err)
	}

	// Inicializar el servicio de WebSocket
	wsService := websocket.NewWebSocketService()

	webMux, err := web.SetupWebRoutes(accountService, fileUploadService, wsService, pgDB, esClient, nc.GetConnection(), templateDir, staticDir)
	if err != nil {
		log.Fatalf("Failed to set up web routes: %v", err)
	}

	// Combinar API y web mux
	mainMux := http.NewServeMux()
	mainMux.Handle("/api/", http.StripPrefix("/api", apiMux))
	mainMux.Handle("/", webMux)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: mainMux,
	}

	go func() {
		log.Printf("Starting HTTP server on port %s", cfg.APIPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Set up gRPC server
	grpcServer := api.SetupGRPCServer(accountService, transactionService)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("Servers shutdown successfully")
}
