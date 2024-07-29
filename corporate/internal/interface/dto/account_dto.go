package dto

import (
	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	UserID   uuid.UUID            `json:"user_id" binding:"required"`
	Currency account.CurrencyCode `json:"currency" binding:"required"`
	Amount   decimal.Decimal      `json:"amount" binding:"required"`
}

type UpdateAccountRequest struct {
	Currency account.CurrencyCode `json:"currency" binding:"required"`
	Balance  decimal.Decimal      `json:"balance" binding:"required"`
}

type DepositRequest struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
}

type WithdrawRequest struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
}
