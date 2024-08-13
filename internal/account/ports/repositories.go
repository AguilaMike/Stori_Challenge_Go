package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	Update(ctx context.Context, account *domain.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int64) ([]*domain.Account, error)
}

type AccountQueryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	List(ctx context.Context, limit, offset int64) ([]*domain.Account, error)
	Search(ctx context.Context, query string) ([]*domain.Account, error)
}
