package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makifdb/mini-bank/speedster/internal/core/service"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

type TransferHandler struct {
	transactionService *service.TransactionService
}

func NewTransferHandler(transactionService *service.TransactionService) *TransferHandler {
	return &TransferHandler{
		transactionService: transactionService,
	}
}

func (h *TransferHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/transfers", h.Transfer)
}

// Transfer godoc
// @Summary Transfer amount between accounts
// @Description Transfer amount from one account to another
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body TransferRequest true "Transfer request"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /transfers [post]
func (h *TransferHandler) Transfer(c *fiber.Ctx) error {
	type TransferRequest struct {
		FromAccountID string `json:"from_account_id" validate:"required"`
		ToAccountID   string `json:"to_account_id" validate:"required"`
		Amount        string `json:"amount" validate:"required"`
		Fee           string `json:"fee"`
	}

	var req TransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	amount, err := utils.NewBigDecimal(req.Amount)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid amount"})
	}

	fee, err := utils.NewBigDecimal(req.Fee)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid fee"})
	}

	err = h.transactionService.Transfer(c.Context(), req.FromAccountID, req.ToAccountID, amount, fee)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Transfer successful"})
}
