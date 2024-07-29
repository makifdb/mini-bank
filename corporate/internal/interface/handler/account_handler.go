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

type AccountHandler struct {
	accountService *service.AccountService
	logger         *zap.Logger
}

func NewAccountHandler(accountService *service.AccountService, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{accountService: accountService, logger: logger}
}

// RegisterAccountRoutes registers account routes.
func (h *AccountHandler) RegisterRoutes(r *gin.RouterGroup) {
	api := r.Group("/accounts")
	{
		api.POST("", h.CreateAccount)
		api.GET("", h.GetAccounts)
		api.GET("/:id", h.GetAccount)
		api.PATCH("/:id", h.UpdateAccount)
		api.DELETE("/:id", h.DeleteAccount)
		api.POST("/:id/deposit", h.Deposit)
		api.POST("/:id/withdraw", h.Withdraw)
	}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account for a user
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   account  body  dto.CreateAccountRequest  true  "Account info"
// @Success 201 {object} account.Account
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	acc, err := h.accountService.CreateAccount(c, req.UserID, req.Currency, req.Amount)
	if err != nil {
		h.logger.Error("Failed to create account", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError(err.Error())))
		return
	}

	c.JSON(http.StatusCreated, acc)
}

// GetAccount godoc
// @Summary Get account by ID
// @Description Get account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Account ID"
// @Success 200 {object} account.Account
// @Failure 400 {object} error.APIError
// @Failure 404 {object} error.APIError
// @Router /accounts/{id} [get]
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("Failed to parse ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
		return
	}

	acc, err := h.accountService.GetAccount(c, id)
	if err != nil {
		h.logger.Error("Failed to find account", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusNotFound, error.NewErrorResponse(*error.NotFound("Account not found")))
		return
	}

	c.JSON(http.StatusOK, acc)
}

// UpdateAccount godoc
// @Summary Update account by ID
// @Description Update account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Account ID"
// @Param   account  body  dto.UpdateAccountRequest  true  "Account info"
// @Success 200 {object} account.Account
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts/{id} [patch]
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	var req dto.UpdateAccountRequest
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

	acc, err := h.accountService.UpdateAccount(c, id, req.Currency, req.Balance)
	if err != nil {
		h.logger.Error("Failed to update account", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to update account")))
		return
	}

	c.JSON(http.StatusOK, acc)
}

// DeleteAccount godoc
// @Summary Delete account by ID
// @Description Delete account by ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Account ID"
// @Success 200 {object} string
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts/{id} [delete]
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Error("Failed to parse ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
		return
	}

	if err := h.accountService.DeleteAccount(c, id); err != nil {
		h.logger.Error("Failed to delete account", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to delete account")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// GetAccounts godoc
// @Summary Get accounts by user ID
// @Description Get accounts by user ID
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   user_id  query  string  true  "User ID"
// @Param   limit  query  int  false  "Limit"
// @Param   offset  query  int  false  "Offset"
// @Success 200 {object} []account.Account
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts [get]
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	limit, offset := 10, 0
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		h.logger.Error("Failed to parse user_id", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid user_id")))
		return
	}

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

	accounts, err := h.accountService.GetAccounts(c, userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to find accounts", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to find accounts")))
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// Deposit godoc
// @Summary Deposit amount to account
// @Description Deposit amount to account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Account ID"
// @Param   amount  body  dto.DepositRequest  true  "Deposit amount"
// @Success 200 {object} account.Account
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts/{id}/deposit [post]
func (h *AccountHandler) Deposit(c *gin.Context) {
	var req dto.DepositRequest
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

	acc, err := h.accountService.Deposit(c, id, req.Amount)
	if err != nil {
		h.logger.Error("Failed to deposit amount", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to deposit amount")))
		return
	}

	c.JSON(http.StatusOK, acc)
}

// Withdraw godoc
// @Summary Withdraw amount from account
// @Description Withdraw amount from account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "Account ID"
// @Param   amount  body  dto.WithdrawRequest  true  "Withdraw amount"
// @Success 200 {object} account.Account
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /accounts/{id}/withdraw [post]
func (h *AccountHandler) Withdraw(c *gin.Context) {
	var req dto.WithdrawRequest
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

	acc, err := h.accountService.Withdraw(c, id, req.Amount)
	if err != nil {
		h.logger.Error("Failed to withdraw amount", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to withdraw amount")))
		return
	}

	c.JSON(http.StatusOK, acc)
}
