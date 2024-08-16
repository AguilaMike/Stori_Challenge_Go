package rest

import (
	"encoding/json"
	"net/http"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/application"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/google/uuid"
)

type AccountHandler struct {
	service *application.AccountService
}

func NewAccountHandler(service *application.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) Manager(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ListAccounts(w, r)
	} else if r.Method == http.MethodPost {
		h.CreateAccount(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func convertAccountToDTO(account *domain.Account) AccountDTO {
	return AccountDTO{
		ID:       account.ID.String(),
		Nickname: account.Nickname,
		Email:    account.Email,
	}
}

func (h *AccountHandler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.ListAccounts(r.Context(), 1000, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir la estructura de dominio a una estructura DTO
	response := make([]AccountDTO, 0, len(accounts))
	for _, user := range accounts {
		response = append(response, convertAccountToDTO(user))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	// Convertir la estructura de dominio a una estructura DTO
	accountDTO := convertAccountToDTO(account)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountDTO)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")
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

	// Convertir la estructura de dominio a una estructura DTO
	accountDTO := convertAccountToDTO(account)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountDTO)
}

// Implement other handler methods (UpdateAccount, DeleteAccount) similarly
