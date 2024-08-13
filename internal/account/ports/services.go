package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
)

type AccountService interface {
	CreateAccount(ctx context.Context, nickname string) (*domain.Account, error)
	GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	ListAccounts(ctx context.Context, limit, offset int32) ([]*domain.Account, error)
	SearchAccounts(ctx context.Context, query string) ([]*domain.Account, error)
}
