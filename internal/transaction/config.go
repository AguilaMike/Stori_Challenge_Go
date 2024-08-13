package transaction

import (
	"database/sql"

	"google.golang.org/grpc"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/email"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/infrastructure"
	"github.com/olivere/elastic/v7"
)

func SetupTransactionDomain(db *sql.DB, esClient *elastic.Client, nc *nats.NatsClient, conn *grpc.ClientConn, sender *email.Sender) *application.TransactionService {
	repo := infrastructure.NewPostgresTransactionRepository(db, nc)
	queryRepo := infrastructure.NewElasticsearchTransactionRepository(esClient, nc, "transactions")
	return application.NewTransactionService(repo, queryRepo, conn, sender)
}
