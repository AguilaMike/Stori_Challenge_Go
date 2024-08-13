package api_grpc

import (
	"context"
	"time"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
	pb "github.com/AguilaMike/Stori_Challenge_Go/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionServer struct {
	pb.UnimplementedTransactionServiceServer
	service *application.TransactionService
}

func NewTransactionServer(service *application.TransactionService) *TransactionServer {
	return &TransactionServer{service: service}
}

func (s *TransactionServer) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.Transaction, error) {
	accountID, err := uuid.Parse(req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account ID: %v", err)
	}

	inputDate := req.InputDate.AsTime()

	transaction, err := s.service.CreateTransaction(ctx, accountID, req.Amount, req.InputFileId, inputDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	return &pb.Transaction{
		Id:          transaction.ID.String(),
		AccountId:   transaction.AccountID.String(),
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		InputFileId: transaction.InputFileID,
		InputDate:   timestamppb.New(transaction.InputDate),
		CreatedAt:   timestamppb.New(time.Unix(int64(transaction.CreatedAt), 0)),
	}, nil
}

func (s *TransactionServer) GetTransactionSummary(ctx context.Context, req *pb.GetTransactionSummaryRequest) (*pb.TransactionSummary, error) {
	accountID, err := uuid.Parse(req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account ID: %v", err)
	}

	summary, err := s.service.GetTransactionSummary(ctx, accountID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transaction summary: %v", err)
	}

	return &pb.TransactionSummary{
		TotalBalance:  summary.TotalBalance,
		TotalCount:    int32(summary.TotalCount),
		AverageCredit: summary.AverageCredit,
		AverageDebit:  summary.AverageDebit,
	}, nil
}

// Implement other gRPC methods (GetTransaction, ListTransactions) similarly
