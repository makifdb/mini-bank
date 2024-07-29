package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/makifdb/mini-bank/corporate/internal/application/service"
	"github.com/makifdb/mini-bank/corporate/internal/interface/dto"
	"github.com/makifdb/mini-bank/corporate/pkg/error"
	"go.uber.org/zap"
)

type TransferHandler struct {
	service *service.TransferService
	logger  *zap.Logger
}

func NewTransferHandler(service *service.TransferService, logger *zap.Logger) *TransferHandler {
	return &TransferHandler{service: service, logger: logger}
}

func RegisterTransferRoutes(app *gin.RouterGroup, handler *TransferHandler) {
	app.POST("/transfers", handler.CreateTransfer)
}

// CreateTransfer godoc
// @Summary Create a new transfer
// @Description Create a new transfer between two accounts
// @Tags transfers
// @Accept  json
// @Produce  json
// @Param   transfer  body  dto.TransferRequest  true  "Transfer info"
// @Success 201 {string} message "Transfer completed successfully"
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /transfers [post]
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	err := h.service.Transfer(c.Request.Context(), req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		h.logger.Error("Transfer failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Transfer failed")))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transfer completed successfully"})
}
