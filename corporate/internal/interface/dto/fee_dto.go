package dto

import (
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/fee"
	"github.com/shopspring/decimal"
)

type CreateFeeRequest struct {
	Amount   decimal.Decimal      `json:"amount" binding:"required"`
	Type     fee.Type             `json:"type" binding:"required"`
	Currency account.CurrencyCode `json:"currency" binding:"required"`
}

type UpdateFeeRequest struct {
	Amount   decimal.Decimal      `json:"amount" binding:"required"`
	Type     fee.Type             `json:"type" binding:"required"`
	Currency account.CurrencyCode `json:"currency" binding:"required"`
}
