package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makifdb/mini-bank/speedster/internal/service"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
)

type FeeHandler struct {
	feeService *service.FeeService
}

func NewFeeHandler(feeService *service.FeeService) *FeeHandler {
	return &FeeHandler{
		feeService: feeService,
	}
}

func (h *FeeHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/fees", h.CreateFee)
	api.Put("/fees/:id", h.UpdateFee)
	api.Delete("/fees/:id", h.DeleteFee)
	api.Get("/fees", h.ListFees)
}

// CreateFee godoc
// @Summary Create a new fee
// @Description Create a new fee
// @Tags fees
// @Accept json
// @Produce json
// @Param fee body CreateFeeRequest true "Fee request"
// @Success 200 {object} models.Fee
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /fees [post]
func (h *FeeHandler) CreateFee(c *fiber.Ctx) error {
	type CreateFeeRequest struct {
		Amount   string              `json:"amount" validate:"required"`
		Type     models.FeeType      `json:"type" validate:"required"`
		Currency models.CurrencyCode `json:"currency" validate:"required"`
	}

	var req CreateFeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	fee, err := h.feeService.CreateFee(c.Context(), req.Amount, req.Type, req.Currency)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fee)
}

// UpdateFee godoc
// @Summary Update a fee
// @Description Update a fee
// @Tags fees
// @Accept json
// @Produce json
// @Param id path string true "Fee ID"
// @Param fee body UpdateFeeRequest true "Fee request"
// @Success 200 {object} models.Fee
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /fees/{id} [put]
func (h *FeeHandler) UpdateFee(c *fiber.Ctx) error {
	type UpdateFeeRequest struct {
		Amount   string              `json:"amount" validate:"required"`
		Type     models.FeeType      `json:"type" validate:"required"`
		Currency models.CurrencyCode `json:"currency" validate:"required"`
	}

	id := c.Params("id")

	var req UpdateFeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	fee, err := h.feeService.UpdateFee(c.Context(), id, req.Amount, req.Type, req.Currency)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fee)
}

// DeleteFee godoc
// @Summary Delete a fee
// @Description Delete a fee
// @Tags fees
// @Accept json
// @Produce json
// @Param id path string true "Fee ID"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /fees/{id} [delete]
func (h *FeeHandler) DeleteFee(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.feeService.DeleteFee(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Fee deleted"})
}

// ListFees godoc
// @Summary List all fees
// @Description List all fees
// @Tags fees
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Fee
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /fees [get]
func (h *FeeHandler) ListFees(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)
	fees, err := h.feeService.GetFees(c.Context(), limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusOK).JSON(fees)
}
