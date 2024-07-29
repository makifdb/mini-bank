package dto

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransferRequest struct {
	FromAccountID uuid.UUID       `json:"from_account_id" binding:"required"`
	ToAccountID   uuid.UUID       `json:"to_account_id" binding:"required"`
	Amount        decimal.Decimal `json:"amount" binding:"required"`
}
