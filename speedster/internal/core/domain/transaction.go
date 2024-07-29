package domain

import "github.com/makifdb/mini-bank/speedster/pkg/utils"

type Transaction struct {
	Base
	FromAccountID string `json:"from_account_id" pg:",notnull"`
	ToAccountID   string `json:"to_account_id" pg:",notnull"`
	Amount        string `json:"amount" pg:",notnull"`
	Fee           string `json:"fee" pg:",notnull"`
}

func NewTransaction(fromAccountID, toAccountID string, amount, fee *utils.BigDecimal) (*Transaction, error) {
	return &Transaction{
		Base:          NewBase(),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount.String(),
		Fee:           fee.String(),
	}, nil
}
