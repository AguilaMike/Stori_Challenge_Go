package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Amount      float64
	Type        string // "credit" or "debit"
	InputFileID string
	InputDate   time.Time
	CreatedAt   int64
}

type TransactionSummary struct {
	AverageCredit float64
	AverageDebit  float64
	CreditCount   int
	DebitCount    int
	TotalBalance  float64
	TotalCount    int
	TotalCredit   float64
	TotalDebit    float64
	Monthly       map[string]*TransactionMonthly
}

type TransactionMonthly struct {
	Year          int
	Month         int
	AverageCredit float64
	AverageDebit  float64
	Balance       float64
	CreditCount   int
	DebitCount    int
	Total         int
	Transactions  []Transaction
}

const (
	TransactionCreatedEvent = "transaction.created"
)

func NewTransaction(accountID uuid.UUID, amount float64, inputFileID string, inputDate time.Time) *Transaction {
	return &Transaction{
		ID:          uuid.New(),
		AccountID:   accountID,
		Amount:      amount,
		Type:        getTransactionType(amount),
		InputFileID: inputFileID,
		InputDate:   inputDate,
		CreatedAt:   time.Now().UTC().Unix(),
	}
}

func getTransactionType(amount float64) string {
	if amount > 0 {
		return "credit"
	}
	return "debit"
}
