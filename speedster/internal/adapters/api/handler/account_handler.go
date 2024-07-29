package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makifdb/mini-bank/speedster/internal/core/domain"
	"github.com/makifdb/mini-bank/speedster/internal/core/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/accounts", h.CreateAccount)
	api.Get("/accounts/:id", h.GetAccountByID)
	api.Put("/accounts/:id", h.UpdateAccount)
	api.Delete("/accounts/:id", h.DeleteAccount)
	api.Post("/accounts/:id/deposit", h.Deposit)
	api.Post("/accounts/:id/withdraw", h.Withdraw)
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body CreateAccountRequest true "Account request"
// @Success 200 {object} domain.Account
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	type CreateAccountRequest struct {
		UserID   string `json:"user_id" validate:"required"`
		Currency string `json:"currency" validate:"required"`
		Amount   string `json:"amount" validate:"required"`
	}

	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	acc, err := h.accountService.CreateAccount(c.Context(), req.UserID, domain.CurrencyCode(req.Currency), req.Amount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(acc)
}

// GetAccountByID godoc
// @Summary Get account by ID
// @Description Get account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} domain.Account
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *fiber.Ctx) error {
	id := c.Params("id")
	acc, err := h.accountService.GetAccount(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(acc)
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param account body UpdateAccountRequest true "Account request"
// @Success 200 {object} domain.Account
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts/{id} [put]
func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	type UpdateAccountRequest struct {
		Currency string `json:"currency" validate:"required"`
		Balance  string `json:"balance" validate:"required"`
	}

	id := c.Params("id")
	var req UpdateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	acc, err := h.accountService.UpdateAccount(c.Context(), id, domain.CurrencyCode(req.Currency), req.Balance)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(acc)
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.accountService.DeleteAccount(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Account deleted"})
}

// Deposit godoc
// @Summary Deposit money into an account
// @Description Deposit money into an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param deposit body DepositRequest true "Deposit request"
// @Success 200 {object} domain.Account
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts/{id}/deposit [post]
func (h *AccountHandler) Deposit(c *fiber.Ctx) error {
	type DepositRequest struct {
		Amount string `json:"amount" validate:"required"`
	}

	id := c.Params("id")
	var req DepositRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	acc, err := h.accountService.Deposit(c.Context(), id, req.Amount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(acc)
}

// Withdraw godoc
// @Summary Withdraw money from an account
// @Description Withdraw money from an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param withdraw body WithdrawRequest true "Withdraw request"
// @Success 200 {object} domain.Account
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /accounts/{id}/withdraw [post]
func (h *AccountHandler) Withdraw(c *fiber.Ctx) error {
	type WithdrawRequest struct {
		Amount string `json:"amount" validate:"required"`
	}

	id := c.Params("id")
	var req WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	acc, err := h.accountService.Withdraw(c.Context(), id, req.Amount)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(acc)
}
