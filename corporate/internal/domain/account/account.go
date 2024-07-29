package account

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/base"
	"github.com/shopspring/decimal"
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
	base.Base
	Balance  decimal.Decimal `json:"balance"`
	Currency CurrencyCode    `json:"code" gorm:"index"`
	UserID   uuid.UUID       `json:"user_id" gorm:"index"`
}

func NewAccount(currency CurrencyCode, balance decimal.Decimal, userID uuid.UUID) (*Account, error) {

	if !IsCurrencyCodeValid(currency) {
		return nil, fmt.Errorf("invalid currency code")
	}

	return &Account{
		Base:     base.NewBase(),
		Currency: currency,
		Balance:  balance,
		UserID:   userID,
	}, nil
}

func (a *Account) Deposit(balance decimal.Decimal) error {
	if balance.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("deposit balance must be greater than zero")
	}
	a.Balance = a.Balance.Add(balance)
	return nil
}

func (a *Account) Withdraw(balance decimal.Decimal) error {
	if balance.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("withdrawal balance must be greater than zero")
	}
	if a.Balance.LessThan(balance) {
		return fmt.Errorf("insufficient balance")
	}
	a.Balance = a.Balance.Sub(balance)
	return nil
}

func (a *Account) Transfer(to *Account, balance decimal.Decimal) error {
	if a.UserID == to.UserID {
		return fmt.Errorf("transfer to the same account is not allowed")
	}
	if err := a.Withdraw(balance); err != nil {
		return err
	}
	if err := to.Deposit(balance); err != nil {
		return err
	}
	return nil
}

func IsCurrencyCodeValid(code CurrencyCode) bool {
	_, ok := validCurrencies[code]
	return ok
}
