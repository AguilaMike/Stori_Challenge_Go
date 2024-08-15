package infrastructure

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/ports"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/db/sqlc"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
)

type PostgresAccountRepository struct {
	queries *sqlc.Queries
	nats    *nats.NatsClient
}

func NewPostgresAccountRepository(db *sql.DB, nc *nats.NatsClient) ports.AccountRepository {
	return &PostgresAccountRepository{
		queries: sqlc.New(db),
		nats:    nc,
	}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	_, err := r.queries.CreateAccount(ctx, sqlc.CreateAccountParams{
		ID:        account.ID,
		Nickname:  account.Nickname,
		Email:     account.Email,
		Balance:   strconv.FormatFloat(account.Balance, 'f', -1, 64),
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		Active:    account.Active,
	})
	if err != nil {
		return err
	}

	// Publish event to NATS
	return r.publishEvent("account.created", account)
}

func (r *PostgresAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	_, err := r.queries.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:        account.ID,
		Nickname:  account.Nickname,
		Email:     account.Email,
		Balance:   strconv.FormatFloat(account.Balance, 'f', -1, 64),
		UpdatedAt: account.UpdatedAt,
		Active:    account.Active,
	})
	if err != nil {
		return err
	}

	// Publish event to NATS
	return r.publishEvent("account.updated", account)
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteAccount(ctx, id)
	if err != nil {
		return err
	}

	// Publish event to NATS
	return r.publishEvent("account.deleted", map[string]string{"id": id.String()})
}

func (r *PostgresAccountRepository) publishEvent(subject string, payload interface{}) error {
	return r.nats.Publish(subject, payload)
}
