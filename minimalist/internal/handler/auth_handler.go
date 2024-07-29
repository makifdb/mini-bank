package handler

import (
	"encoding/json"
	"net/http"

	"github.com/makifdb/mini-bank/minimalist/internal/service"
)

type AuthHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AuthHandler {
	return &AuthHandler{
		adminService: adminService,
	}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", h.signUp)
	mux.HandleFunc("POST /login", h.login)
	mux.HandleFunc("POST /verify", h.verifyAdmin)
	mux.HandleFunc("POST /refresh", h.refreshToken)
	mux.HandleFunc("POST /logout", h.logout)
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {

	type SignUpRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	admin, err := h.adminService.SignUp(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {

	type LoginRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.adminService.Login(r.Context(), req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) verifyAdmin(w http.ResponseWriter, r *http.Request) {
	type VerifyRequest struct {
		Email string `json:"email" validate:"required,email"`
		Code  string `json:"code" validate:"required"`
	}

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.adminService.VerifyAdmin(r.Context(), req.Email, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var RefreshRequest struct {
		Token string `json:"token" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&RefreshRequest); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.adminService.RefreshToken(r.Context(), RefreshRequest.Token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
