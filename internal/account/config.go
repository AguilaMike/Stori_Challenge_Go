package account

import (
	"database/sql"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/infrastructure"
	"github.com/olivere/elastic/v7"
)

func SetupAccountDomain(db *sql.DB, esClient *elastic.Client) *application.AccountService {
	repo := infrastructure.NewPostgresAccountRepository(db)
	queryRepo := infrastructure.NewElasticsearchAccountRepository(esClient, "accounts")
	return application.NewAccountService(repo, queryRepo)
}
