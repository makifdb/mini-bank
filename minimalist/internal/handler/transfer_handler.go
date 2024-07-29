package handler

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/makifdb/mini-bank/minimalist/internal/middleware"
	"github.com/makifdb/mini-bank/minimalist/internal/service"
)

type TransferHandler struct {
	transactionService *service.TransactionService
}

func NewTransferHandler(transactionService *service.TransactionService) *TransferHandler {
	return &TransferHandler{
		transactionService: transactionService,
	}
}

func (h *TransferHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /transfers", middleware.WrapHandlerFunc(h.Transfer))
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	type TransferRequest struct {
		FromAccountID int64   `json:"from_account_id" validate:"required"`
		ToAccountID   int64   `json:"to_account_id" validate:"required"`
		Amount        float64 `json:"amount" validate:"required"`
		Fee           float64 `json:"fee"`
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	amt := new(big.Float).SetPrec(32).SetFloat64(req.Amount)
	fee := new(big.Float).SetPrec(32).SetFloat64(req.Fee)

	err := h.transactionService.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, amt, fee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
