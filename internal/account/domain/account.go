package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID
	Nickname  string
	Email     string
	Balance   float64
	CreatedAt int64
	UpdatedAt int64
	Active    bool
}

func NewAccount(nickname, email string) *Account {
	now := time.Now().UTC().Unix()
	return &Account{
		ID:        uuid.New(),
		Nickname:  nickname,
		Email:     email,
		Balance:   0,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
}

func (a *Account) UpdateBalance(amount float64) {
	a.Balance += amount
	a.UpdatedAt = time.Now().UTC().Unix()
}

func (a *Account) Deactivate() {
	a.Active = false
	a.UpdatedAt = time.Now().UTC().Unix()
}
