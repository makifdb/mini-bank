package models

import "math/big"

type Transaction struct {
	Base
	FromAccountID int64      `json:"from_account_id"`
	ToAccountID   int64      `json:"to_account_id"`
	Amount        *big.Float `json:"amount"`
	Fee           *big.Float `json:"fee"`
}

func NewTransaction(fromAccountID, toAccountID int64, amount, fee *big.Float) (*Transaction, error) {
	return &Transaction{
		Base:          NewBase(),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Fee:           fee,
	}, nil
}
