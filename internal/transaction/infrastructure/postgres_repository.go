package infrastructure

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/db/sqlc"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/ports"
)

type PostgresTransactionRepository struct {
	queries *sqlc.Queries
	db      *sql.DB
	nats    *nats.NatsClient
}

func NewPostgresTransactionRepository(db *sql.DB, nc *nats.NatsClient) ports.TransactionRepository {
	return &PostgresTransactionRepository{
		queries: sqlc.New(db),
		db:      db,
		nats:    nc,
	}
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	_, err := r.queries.CreateTransaction(ctx, sqlc.CreateTransactionParams{
		ID:          transaction.ID,
		AccountID:   transaction.AccountID,
		Amount:      strconv.FormatFloat(transaction.Amount, 'f', -1, 64),
		Type:        transaction.Type,
		InputFileID: transaction.InputFileID,
		InputDate:   transaction.InputDate,
		CreatedAt:   transaction.CreatedAt,
	})
	if err != nil {
		return err
	}

	// Publish a message to NATS
	return r.publishEvent("transaction.created", transaction)
}

func (r *PostgresTransactionRepository) CreateBulk(ctx context.Context, transactions []*domain.Transaction) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Use the WithTx method to run queries within the transaction
	//qtx := r.queries.WithTx(tx)

	for _, t := range transactions {
		err = r.Create(ctx, t)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresTransactionRepository) publishEvent(subject string, payload interface{}) error {
	return r.nats.Publish(subject, payload)
}
