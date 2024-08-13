package rest

import (
	"encoding/json"
	"net/http"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	"github.com/google/uuid"
)

type AccountHandler struct {
	service *application.AccountService
}

func NewAccountHandler(service *application.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account, err := h.service.CreateAccount(r.Context(), input.Nickname, input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

// Implement other handler methods (UpdateAccount, DeleteAccount, ListAccounts) similarly
