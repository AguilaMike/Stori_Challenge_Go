package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/ports"
)

type AccountService struct {
	repo  ports.AccountRepository
	query ports.AccountQueryRepository
}

func NewAccountService(repo ports.AccountRepository, query ports.AccountQueryRepository) *AccountService {
	return &AccountService{
		repo:  repo,
		query: query,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, nickname, email string) (*domain.Account, error) {
	account := domain.NewAccount(nickname, email)
	err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return s.query.GetByID(ctx, id)
}

func (s *AccountService) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return s.repo.Update(ctx, account)
}

func (s *AccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *AccountService) ListAccounts(ctx context.Context, limit, offset int64) ([]*domain.Account, error) {
	return s.query.List(ctx, limit, offset)
}

func (s *AccountService) SearchAccounts(ctx context.Context, query string) ([]*domain.Account, error) {
	return s.query.Search(ctx, query)
}
