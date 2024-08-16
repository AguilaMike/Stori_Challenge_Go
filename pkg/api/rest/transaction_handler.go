package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	transaction "github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/application"
)

type TransactionHandler struct {
	service *transaction.TransactionService
}

func NewTransactionHandler(service *transaction.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) GetTransactionSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//accountIDStr := r.URL.Query().Get("account_id")
	accountIDStr := r.PathValue("account_id")
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

	// Convert domain structure to DTO structure
	data := TransactionDTO{
		AccounID: accountID.String(),
		Summary: TransactionSummaryDTO{
			AverageCredit: summary.AverageCredit,
			AverageDebit:  summary.AverageDebit,
			CreditCount:   summary.CreditCount,
			DebitCount:    summary.DebitCount,
			TotalBalance:  summary.TotalBalance,
			TotalCount:    summary.TotalCount,
			TotalCredit:   summary.TotalCredit,
			TotalDebit:    summary.TotalDebit,
		},
		Monthly: make(map[string]*TransactionMonthlyDTO),
	}

	for _, v := range summary.Monthly {
		key := fmt.Sprintf("%d-%d", v.Year, v.Month)
		if _, ok := data.Monthly[key]; !ok {
			data.Monthly[key] = &TransactionMonthlyDTO{
				Year:          v.Year,
				Month:         v.Month,
				Total:         v.Total,
				Balance:       v.Balance,
				AverageCredit: v.AverageCredit,
				AverageDebit:  v.AverageDebit,
				Transactions:  make([]TransactionDetailDTO, 0, len(v.Transactions)),
			}

			for _, t := range v.Transactions {
				data.Monthly[key].Transactions = append(data.Monthly[key].Transactions, TransactionDetailDTO{
					ID:        t.ID.String(),
					Amount:    t.Amount,
					Type:      t.Type,
					InputDate: t.InputDate.Format("2006-01-02"),
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Implement other handler methods (GetTransaction, ListTransactions) similarly
