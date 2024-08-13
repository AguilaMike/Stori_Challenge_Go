package grpc

import (
	"context"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	pb "github.com/AguilaMike/Stori_Challenge_Go/pkg/proto"
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
		return nil, err
	}

	return &pb.Account{
		Id:       account.ID.String(),
		Nickname: account.NickName,
		Balance:  account.Balance,
		// Set other fields
	}, nil
}

// Implementa otros métodos del servicio gRPC aquí
