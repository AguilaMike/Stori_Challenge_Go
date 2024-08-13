package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	service *application.TransactionService
}

func NewTransactionHandler(service *application.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AccountID   string    `json:"account_id"`
		Amount      float64   `json:"amount"`
		Type        string    `json:"type"`
		InputFileID string    `json:"input_file_id"`
		InputDate   time.Time `json:"input_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountID, err := uuid.Parse(input.AccountID)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.CreateTransaction(r.Context(), accountID, input.Amount, input.InputFileID, input.InputDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

func (h *TransactionHandler) GetTransactionSummary(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Query().Get("account_id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	summary, err := h.service.GetTransactionSummary(r.Context(), accountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// Implement other handler methods (GetTransaction, ListTransactions) similarly
