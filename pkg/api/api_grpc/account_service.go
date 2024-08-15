package api_grpc

import (
	"context"
	"time"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	pb "github.com/AguilaMike/Stori_Challenge_Go/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountServer struct {
	pb.UnimplementedAccountServiceServer
	service *application.AccountService
}

func NewAccountServer(service *application.AccountService) *AccountServer {
	return &AccountServer{service: service}
}

func (s *AccountServer) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	account, err := s.service.CreateAccount(ctx, req.Nickname, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	return &pb.Account{
		Id:        account.ID.String(),
		Nickname:  account.Nickname,
		Balance:   account.Balance,
		CreatedAt: timestamppb.New(time.Unix(int64(account.CreatedAt), 0)),
		UpdatedAt: timestamppb.New(time.Unix(int64(account.UpdatedAt), 0)),
		Active:    account.Active,
	}, nil
}

func (s *AccountServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account ID: %v", err)
	}

	account, err := s.service.GetAccount(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get account: %v", err)
	}

	return &pb.Account{
		Id:        account.ID.String(),
		Nickname:  account.Nickname,
		Balance:   account.Balance,
		CreatedAt: timestamppb.New(time.Unix(int64(account.CreatedAt), 0)),
		UpdatedAt: timestamppb.New(time.Unix(int64(account.UpdatedAt), 0)),
		Active:    account.Active,
	}, nil
}

// Implement other gRPC methods (UpdateAccount, DeleteAccount, ListAccounts) similarly
