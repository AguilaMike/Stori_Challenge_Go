package account

import (
	"database/sql"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/infrastructure"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/olivere/elastic/v7"
)

func SetupAccountDomain(db *sql.DB, esClient *elastic.Client, nc *nats.NatsClient) *application.AccountService {
	repo := infrastructure.NewPostgresAccountRepository(db, nc)
	queryRepo := infrastructure.NewElasticsearchAccountRepository(esClient, nc, "accounts")
	return application.NewAccountService(repo, queryRepo)
}
