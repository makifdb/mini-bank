package models

import (
	"fmt"

	"github.com/cockroachdb/apd/v3"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

// CurrencyCode represents the code of a currency.
type CurrencyCode string

const (
	USD CurrencyCode = "USD"
	EUR CurrencyCode = "EUR"
	TRY CurrencyCode = "TRY"
)

var validCurrencies = map[CurrencyCode]struct{}{
	USD: {},
	EUR: {},
	TRY: {},
}

// Account represents a user account.
type Account struct {
	Base
	Balance  *utils.BigDecimal `json:"balance"`
	Currency CurrencyCode      `json:"code" gorm:"index"`
	UserID   string            `json:"user_id" gorm:"index"`
}

func NewAccount(currency CurrencyCode, balance string, userID string) (*Account, error) {
	if !isValidCurrencyCode(currency) {
		return nil, fmt.Errorf("invalid currency code")
	}

	bal, err := utils.NewBigDecimal(balance)
	if err != nil {
		return nil, err
	}

	return &Account{
		Base:     NewBase(),
		Currency: currency,
		Balance:  bal,
		UserID:   userID,
	}, nil
}

func (a *Account) Deposit(amount string) error {
	amt, err := utils.NewBigDecimal(amount)
	if err != nil {
		return err
	}
	a.Balance = a.Balance.Add(amt)
	return nil
}

func (a *Account) Withdraw(amount string) error {
	amt, err := utils.NewBigDecimal(amount)
	if err != nil {
		return err
	}
	a.Balance = a.Balance.Sub(amt)
	if a.Balance.Value().Cmp(apd.New(0, 0)) < 0 {
		return fmt.Errorf("insufficient balance")
	}
	return nil
}

func (a *Account) Transfer(to *Account, amount string) error {
	if a.UserID == to.UserID {
		return fmt.Errorf("transfer to the same account is not allowed")
	}
	if err := a.Withdraw(amount); err != nil {
		return err
	}
	return to.Deposit(amount)
}

func isValidCurrencyCode(code CurrencyCode) bool {
	_, ok := validCurrencies[code]
	return ok
}
