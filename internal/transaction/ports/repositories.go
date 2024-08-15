package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	CreateBulk(ctx context.Context, transactions []*domain.Transaction) error
}

type TransactionQueryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]*domain.Transaction, error)
	GetSummary(ctx context.Context, accountID uuid.UUID) (*domain.TransactionSummary, error)
}
