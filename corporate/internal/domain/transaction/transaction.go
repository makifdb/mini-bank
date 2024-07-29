package transaction

import (
	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/base"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	base.Base
	FromAccountID uuid.UUID       `json:"from_account_id" gorm:"index"`
	ToAccountID   uuid.UUID       `json:"to_account_id" gorm:"index"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
}

func NewTransaction(fromAccountID, toAccountID uuid.UUID, amount, fee decimal.Decimal) (*Transaction, error) {
	return &Transaction{
		Base:          base.NewBase(),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Fee:           fee,
	}, nil
}
