package ports

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, accountID uuid.UUID, amount float64, transactionType, inputFileID string, inputDate time.Time) (*domain.Transaction, error)
	GetTransaction(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int32) ([]*domain.Transaction, error)
	GetTransactionSummary(ctx context.Context, accountID uuid.UUID) (*domain.TransactionSummary, error)
}
