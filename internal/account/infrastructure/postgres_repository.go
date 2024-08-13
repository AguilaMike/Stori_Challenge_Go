package infrastructure

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/ports"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/db/sqlc"
)

type PostgresAccountRepository struct {
	queries *sqlc.Queries
}

func NewPostgresAccountRepository(db *sql.DB) ports.AccountRepository {
	return &PostgresAccountRepository{
		queries: sqlc.New(db),
	}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	_, err := r.queries.CreateAccount(ctx, sqlc.CreateAccountParams{
		ID:        account.ID,
		Nickname:  account.NickName,
		Balance:   strconv.FormatFloat(account.Balance, 'f', -1, 64),
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		Active:    account.Active,
	})
	return err
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	dbAccount, err := r.queries.GetAccount(ctx, id)
	if err != nil {
		return nil, err
	}
	balance, _ := strconv.ParseFloat(dbAccount.Balance, 64)
	return &domain.Account{
		ID:        dbAccount.ID,
		NickName:  dbAccount.Nickname,
		Balance:   balance,
		CreatedAt: dbAccount.CreatedAt,
		UpdatedAt: dbAccount.UpdatedAt,
		Active:    dbAccount.Active,
	}, nil
}

func (r *PostgresAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	_, err := r.queries.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:        account.ID,
		Nickname:  account.NickName,
		Balance:   strconv.FormatFloat(account.Balance, 'f', -1, 64),
		UpdatedAt: account.UpdatedAt,
		Active:    account.Active,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteAccount(ctx, id)
}

func (r *PostgresAccountRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Account, error) {
	dbAccounts, err := r.queries.ListAccounts(ctx, sqlc.ListAccountsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	accounts := make([]*domain.Account, len(dbAccounts))
	for i, dbAccount := range dbAccounts {
		balance, _ := strconv.ParseFloat(dbAccount.Balance, 64)
		accounts[i] = &domain.Account{
			ID:        dbAccount.ID,
			NickName:  dbAccount.Nickname,
			Balance:   balance,
			CreatedAt: dbAccount.CreatedAt,
			UpdatedAt: dbAccount.UpdatedAt,
			Active:    dbAccount.Active,
		}
	}
	return accounts, nil
}
