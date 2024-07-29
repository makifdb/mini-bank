package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/application/service"
	"github.com/makifdb/mini-bank/corporate/internal/interface/dto"
	"github.com/makifdb/mini-bank/corporate/pkg/error"
	"go.uber.org/zap"
)

type FeeHandler struct {
	feeService *service.FeeService
	logger     *zap.Logger
}

func NewFeeHandler(feeService *service.FeeService, logger *zap.Logger) *FeeHandler {
	return &FeeHandler{
		feeService: feeService,
		logger:     logger,
	}
}

func (h *FeeHandler) RegisterRoutes(app *gin.RouterGroup) {
	app.POST("/fees", h.createFee())
	app.GET("/fees", h.getFees())
	app.PATCH("/fees/:id", h.updateFee())
	app.DELETE("/fees/:id", h.deleteFee())
}

// CreateFee godoc
// @Summary Create a new fee
// @Description Create a new fee
// @Tags fees
// @Accept  json
// @Produce  json
// @Param   fee  body  dto.CreateFeeRequest  true  "Fee info"
// @Success 201 {object} fee.Fee
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
func (h *FeeHandler) createFee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateFeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
			return
		}

		f, err := h.feeService.CreateFee(c, req.Amount, req.Type, req.Currency)
		if err != nil {
			h.logger.Error("Failed to create fee", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to create fee")))
			return
		}

		c.JSON(http.StatusCreated, f)
	}
}

// GetFees godoc
// @Summary Get all fees
// @Description Get all fees
// @Tags fees
// @Accept  json
// @Produce  json
// @Param   limit  query  int  false  "Limit"
// @Param   offset  query  int  false  "Offset"
// @Success 200 {object} []fee.Fee
// @Failure 500 {object} error.APIError
// @Router /fees [get]
func (h *FeeHandler) getFees() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset := 10, 0
		if l, ok := c.GetQuery("limit"); ok {
			if lInt, err := strconv.Atoi(l); err == nil {
				limit = lInt
			}
		}

		if o, ok := c.GetQuery("offset"); ok {
			if oInt, err := strconv.Atoi(o); err == nil {
				offset = oInt
			}
		}

		fees, err := h.feeService.GetFees(c, limit, offset)
		if err != nil {
			h.logger.Error("Failed to find fees", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to find fees")))
			return
		}

		c.JSON(http.StatusOK, fees)
	}
}

// UpdateFee godoc
// @Summary Update fee
// @Description Update fee
// @Tags fees
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Fee ID"
// @Param   fee  body  dto.UpdateFeeRequest  true  "Fee info"
// @Success 200 {object} fee.Fee
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /fees/{id} [patch]
func (h *FeeHandler) updateFee() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UpdateFeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
			return
		}

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.Error("Failed to parse ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
			return
		}

		f, err := h.feeService.UpdateFee(c, id, req.Amount, req.Type, req.Currency)
		if err != nil {
			h.logger.Error("Failed to update fee", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to update fee")))
			return
		}

		c.JSON(http.StatusOK, f)
	}
}

// DeleteFee godoc
// @Summary Delete fee
// @Description Delete fee
// @Tags fees
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Fee ID"
// @Success 204
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /fees/{id} [delete]
func (h *FeeHandler) deleteFee() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.Error("Failed to parse ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
			return
		}

		if err := h.feeService.DeleteFee(c, id); err != nil {
			h.logger.Error("Failed to delete fee", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to delete fee")))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
