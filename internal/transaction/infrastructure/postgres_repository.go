package infrastructure

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/google/uuid"

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

func (r *PostgresTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	dbTransaction, err := r.queries.GetTransaction(ctx, id)
	if err != nil {
		return nil, err
	}
	amount, _ := strconv.ParseFloat(dbTransaction.Amount, 64)
	return &domain.Transaction{
		ID:          dbTransaction.ID,
		AccountID:   dbTransaction.AccountID,
		Amount:      amount,
		Type:        dbTransaction.Type,
		InputFileID: dbTransaction.InputFileID,
		InputDate:   dbTransaction.InputDate,
		CreatedAt:   dbTransaction.CreatedAt,
	}, nil
}

func (r *PostgresTransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]*domain.Transaction, error) {
	dbTransactions, err := r.queries.ListTransactionsByAccount(ctx, sqlc.ListTransactionsByAccountParams{
		AccountID: accountID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	transactions := make([]*domain.Transaction, len(dbTransactions))
	for i, dbTx := range dbTransactions {
		amount, _ := strconv.ParseFloat(dbTx.Amount, 64)
		transactions[i] = &domain.Transaction{
			ID:          dbTx.ID,
			AccountID:   dbTx.AccountID,
			Amount:      amount,
			Type:        dbTx.Type,
			InputFileID: dbTx.InputFileID,
			InputDate:   dbTx.InputDate,
			CreatedAt:   dbTx.CreatedAt,
		}
	}
	return transactions, nil
}

func (r *PostgresTransactionRepository) GetSummary(ctx context.Context, accountID uuid.UUID) (*domain.TransactionSummary, error) {
	summary, err := r.queries.GetTransactionSummary(ctx, accountID)
	if err != nil {
		return nil, err
	}
	totalBalance, _ := strconv.ParseFloat(summary.TotalBalance, 64)
	averageCredit, _ := strconv.ParseFloat(summary.AverageCredit, 64)
	averageDebit, _ := strconv.ParseFloat(summary.AverageDebit, 64)
	return &domain.TransactionSummary{
		TotalBalance:  totalBalance,
		TotalCount:    int(summary.TransactionCount),
		AverageCredit: averageCredit,
		AverageDebit:  averageDebit,
	}, nil
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
	err = r.nats.Publish(domain.TransactionCreatedEvent, transaction)

	return err
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
		// _, err := qtx.CreateTransaction(ctx, sqlc.CreateTransactionParams{
		// 	ID:        t.ID,
		// 	AccountID: t.AccountID,
		// 	Amount:    strconv.FormatFloat(t.Amount, 'f', -1, 64),
		// 	Type:      t.Type,
		// 	InputDate: t.InputDate,
		// 	CreatedAt: t.CreatedAt,
		// })
		// if err != nil {
		// 	return err
		// }
		err = r.Create(ctx, t)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
