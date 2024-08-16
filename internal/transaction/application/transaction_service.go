package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/email"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/ports"
	pb "github.com/AguilaMike/Stori_Challenge_Go/pkg/proto"
)

type TransactionService struct {
	repo    ports.TransactionRepository
	query   ports.TransactionQueryRepository
	account pb.AccountServiceClient
	sender  *email.Sender
}

func NewTransactionService(repo ports.TransactionRepository, query ports.TransactionQueryRepository, conn *grpc.ClientConn, sender *email.Sender) *TransactionService {
	return &TransactionService{
		repo:    repo,
		query:   query,
		account: pb.NewAccountServiceClient(conn),
		sender:  sender,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, accountID uuid.UUID, amount float64, inputFileID string, inputDate time.Time) (*domain.Transaction, error) {
	transaction := domain.NewTransaction(accountID, amount, inputFileID, inputDate)
	err := s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return transaction, nil
}

func (s *TransactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	return s.query.GetByID(ctx, id)
}

func (s *TransactionService) GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]*domain.Transaction, error) {
	return s.query.GetByAccountID(ctx, accountID, limit, offset)
}

func (s *TransactionService) GetTransactionSummary(ctx context.Context, accountID uuid.UUID) (*domain.TransactionSummary, error) {
	transactions, err := s.query.GetByAccountID(ctx, accountID, 1000, 0)
	if err != nil {
		return nil, err
	}

	summary := &domain.TransactionSummary{
		Monthly: make(map[string]*domain.TransactionMonthly),
	}

	for _, t := range transactions {
		key := t.InputDate.Format("2006-01")
		summary.TotalCount++
		summary.TotalBalance += t.Amount

		if _, ok := summary.Monthly[key]; !ok {
			summary.Monthly[key] = &domain.TransactionMonthly{
				Year:         t.InputDate.Year(),
				Month:        int(t.InputDate.Month()),
				Transactions: make([]domain.Transaction, 0),
			}
		}

		summary.Monthly[key].Transactions = append(summary.Monthly[key].Transactions, *t)
		summary.Monthly[key].Total++

		if t.Amount > 0 {
			summary.CreditCount++
			summary.TotalCredit += t.Amount
			summary.Monthly[key].AverageCredit += t.Amount
			summary.Monthly[key].CreditCount++
			summary.Monthly[key].Balance += t.Amount
		} else {
			summary.DebitCount++
			summary.TotalDebit += t.Amount
			summary.Monthly[key].AverageDebit += t.Amount
			summary.Monthly[key].DebitCount++
			summary.Monthly[key].Balance -= t.Amount
		}
	}

	// Recorremos los Montlhy del summary
	for k, v := range summary.Monthly {
		if v.CreditCount > 0 {
			v.AverageCredit = v.AverageCredit / float64(v.CreditCount)
		}

		if v.DebitCount > 0 {
			v.AverageDebit = v.AverageDebit / float64(v.DebitCount)
		}
		summary.Monthly[k] = v
	}

	if summary.CreditCount > 0 {
		summary.AverageCredit = summary.TotalCredit / float64(summary.CreditCount)
	}

	if summary.DebitCount > 0 {
		summary.AverageDebit = summary.TotalDebit / float64(summary.DebitCount)
	}

	return summary, nil
}

func (s *TransactionService) CreateBulkTransactions(ctx context.Context, transactions []*domain.Transaction) error {
	return s.repo.CreateBulk(ctx, transactions)
}

func (s *TransactionService) SendSummaryEmail(ctx context.Context, summary *domain.TransactionSummary, userID uuid.UUID) error {
	// Send email to user
	request := &pb.GetAccountRequest{Id: userID.String()}
	user, err := s.account.GetAccount(ctx, request)
	if err != nil {
		return err
	}

	s.sender.SendWithTemplate(user.Email, "summary.gohtml", summary)

	return nil
}
