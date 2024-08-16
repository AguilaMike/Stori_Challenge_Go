package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/config"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/db"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/elasticsearch"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/email"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/websocket"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/infrastructure"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	pgDB, err := db.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgDB.Close()

	// Ejecutar migraciones
	if err := db.RunMigrations(pgDB, "/root/migrations"); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

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

	// Initialize repositories
	hostGrpc := fmt.Sprintf("%s:%s", cfg.DOMAIN, cfg.GRPCPort)
	connGrpc, err := grpc.NewClient(hostGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	emailSender := email.NewSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	transactionRepo := infrastructure.NewPostgresTransactionRepository(pgDB, nc)
	transactionQueryRepo := infrastructure.NewElasticsearchTransactionRepository(esClient, nc, "transactions")

	// Initialize service
	transactionService := application.NewTransactionService(transactionRepo, transactionQueryRepo, connGrpc, emailSender)

	// Initialize NATS client
	natsClient, err := nats.NewNatsClient(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	// Initialize WebSocket service
	wsService := websocket.NewWebSocketService()

	// Set up your worker logic here
	err = setupWorkerTasks(natsClient, transactionService, wsService)
	if err != nil {
		log.Fatalf("Failed to set up worker tasks: %v", err)
	}

	log.Println("Worker started successfully")

	// Wait for interrupt signal to gracefully shut down the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
}

func setupWorkerTasks(
	natsClient *nats.NatsClient,
	transactionService *application.TransactionService,
	wsService *websocket.WebSocketService) error {

	_, err := natsClient.Subscribe("transaction.file.uploaded", func(data []byte) {
		var fileInfo struct {
			FilePath string    `json:"file_path"`
			UserID   uuid.UUID `json:"user_id"`
		}

		log.Printf("Processing file: %s", data)
		err := json.Unmarshal(data, &fileInfo)
		if err != nil {
			log.Printf("Error unmarshaling file info: %v", err)
			return
		}

		transactions, err := processTransactionFile(fileInfo.FilePath, fileInfo.UserID)
		if err != nil {
			log.Printf("Error processing transaction file: %v", err)
			return
		}

		ctx := context.Background()
		err = transactionService.CreateBulkTransactions(ctx, transactions)
		if err != nil {
			log.Printf("Error saving transactions: %v", err)
			return
		}

		summary, err := transactionService.GetTransactionSummary(ctx, fileInfo.UserID)
		if err != nil {
			log.Printf("Error getting transaction summary: %v", err)
			return
		}

		err = transactionService.SendSummaryEmail(ctx, summary, fileInfo.UserID)
		if err != nil {
			log.Printf("Error sending summary email: %v", err)
			return
		}

		// Enviar actualización a través de WebSocket
		updateMessage, _ := json.Marshal(map[string]interface{}{
			"type":    "transaction_update",
			"summary": summary,
		})
		wsService.SendUpdate(fileInfo.UserID.String(), updateMessage)

		log.Printf("Successfully processed file %s for user %s", fileInfo.FilePath, fileInfo.UserID)
	})

	return err
}

func processTransactionFile(filePath string, userID uuid.UUID) ([]*domain.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []*domain.Transaction
	for _, record := range records {
		if len(record) != 2 {
			log.Printf("Skipping invalid record: %v", record)
			continue // Skip invalid records
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			log.Printf("Invalid date format: %v", err)
			continue // Skip invalid dates
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Invalid amount: %v", err)
			continue // Skip invalid amounts
		}

		if record[1][0] == '-' {
			amount *= -1
		}
		transaction := domain.NewTransaction(userID, amount, filePath, date)
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
