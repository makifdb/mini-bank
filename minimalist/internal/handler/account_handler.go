package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/makifdb/mini-bank/minimalist/internal/middleware"
	"github.com/makifdb/mini-bank/minimalist/internal/models"
	"github.com/makifdb/mini-bank/minimalist/internal/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /accounts", middleware.WrapHandlerFunc(h.CreateAccount))
	mux.HandleFunc("GET /accounts/{id}", middleware.WrapHandlerFunc(h.GetAccountByID))
	mux.HandleFunc("PATCH /accounts/{id}", middleware.WrapHandlerFunc(h.UpdateAccount))
	mux.HandleFunc("DELETE /accounts/{id}", middleware.WrapHandlerFunc(h.DeleteAccount))
	mux.HandleFunc("POST /accounts/{id}/deposit", middleware.WrapHandlerFunc(h.Deposit))
	mux.HandleFunc("POST /accounts/{id}/withdraw", middleware.WrapHandlerFunc(h.Withdraw))
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	type CreateAccountRequest struct {
		UserID   int64   `json:"user_id" validate:"required"`
		Currency string  `json:"currency" validate:"required"`
		Amount   float64 `json:"amount" validate:"required"`
	}

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.CreateAccount(r.Context(), req.UserID, req.Amount, models.CurrencyCode(req.Currency))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}
	acc, err := h.accountService.GetAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	type UpdateAccountRequest struct {
		Currency string  `json:"currency" validate:"required"`
		Balance  float64 `json:"balance" validate:"required"`
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.UpdateAccount(r.Context(), id, models.CurrencyCode(req.Currency), req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.accountService.DeleteAccount(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {

	type DepositRequest struct {
		Amount float64 `json:"amount" validate:"required"`
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.Deposit(r.Context(), id, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func (h *AccountHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	type WithdrawRequest struct {
		Amount float64 `json:"amount" validate:"required"`
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := h.accountService.Withdraw(r.Context(), id, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}
