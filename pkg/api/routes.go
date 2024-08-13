package api

import (
	"net/http"

	appAccount "github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	appTran "github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
	"github.com/AguilaMike/Stori_Challenge_Go/pkg/api/api_grpc"
	"github.com/AguilaMike/Stori_Challenge_Go/pkg/api/rest"
	pbAccount "github.com/AguilaMike/Stori_Challenge_Go/pkg/proto"
	"google.golang.org/grpc"
)

func SetupHTTPRoutes(accountService *appAccount.AccountService, transactionService *appTran.TransactionService) *http.ServeMux {
	router := http.NewServeMux()

	accountHandler := rest.NewAccountHandler(accountService)
	transactionHandler := rest.NewTransactionHandler(transactionService)

	router.HandleFunc("/accounts", accountHandler.CreateAccount)
	router.HandleFunc("/accounts/{id}", accountHandler.GetAccount)
	router.HandleFunc("/transactions", transactionHandler.CreateTransaction)
	router.HandleFunc("/transactions/summary", transactionHandler.GetTransactionSummary)

	return router
}

func SetupGRPCServer(accountService *appAccount.AccountService, transactionService *appTran.TransactionService) *grpc.Server {
	grpcServer := grpc.NewServer()

	pbAccount.RegisterAccountServiceServer(grpcServer, api_grpc.NewAccountServer(accountService))
	pbAccount.RegisterTransactionServiceServer(grpcServer, api_grpc.NewTransactionServer(transactionService))

	return grpcServer
}
