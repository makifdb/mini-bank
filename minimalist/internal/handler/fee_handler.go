package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/makifdb/mini-bank/minimalist/internal/middleware"
	"github.com/makifdb/mini-bank/minimalist/internal/models"
	"github.com/makifdb/mini-bank/minimalist/internal/service"
)

type FeeHandler struct {
	feeService *service.FeeService
}

func NewFeeHandler(feeService *service.FeeService) *FeeHandler {
	return &FeeHandler{
		feeService: feeService,
	}
}

func (h *FeeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /fees", middleware.WrapHandlerFunc(h.CreateFee))
	mux.HandleFunc("PATCH /fees/{id}", middleware.WrapHandlerFunc(h.UpdateFee))
	mux.HandleFunc("DELETE /fees/{id}", middleware.WrapHandlerFunc(h.DeleteFee))
	mux.HandleFunc("GET /fees", middleware.WrapHandlerFunc(h.ListFees))
}

func (h *FeeHandler) CreateFee(w http.ResponseWriter, r *http.Request) {
	type CreateFeeRequest struct {
		Amount   float64             `json:"amount" validate:"required"`
		Type     models.FeeType      `json:"type" validate:"required"`
		Currency models.CurrencyCode `json:"currency" validate:"required"`
	}

	var req CreateFeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	fee, err := h.feeService.CreateFee(r.Context(), req.Amount, req.Type, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fee)
}

func (h *FeeHandler) UpdateFee(w http.ResponseWriter, r *http.Request) {
	type UpdateFeeRequest struct {
		Amount   float64             `json:"amount" validate:"required"`
		Type     models.FeeType      `json:"type" validate:"required"`
		Currency models.CurrencyCode `json:"currency" validate:"required"`
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	var req UpdateFeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	fee, err := h.feeService.UpdateFee(r.Context(), id, req.Amount, req.Type, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fee)
}

func (h *FeeHandler) DeleteFee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.feeService.DeleteFee(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *FeeHandler) ListFees(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	fees, err := h.feeService.GetFees(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fees)
}
