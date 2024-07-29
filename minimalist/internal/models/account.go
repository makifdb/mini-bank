package models

import (
	"fmt"
	"math/big"
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
	Balance  big.Float    `json:"balance"`
	Currency CurrencyCode `json:"currency_code"`
	UserID   int64        `json:"user_id"`
}

func NewAccount(currency CurrencyCode, balance big.Float, userID int64) (*Account, error) {
	if !isValidCurrencyCode(currency) {
		return nil, fmt.Errorf("invalid currency code")
	}

	return &Account{
		Base:     NewBase(),
		Currency: currency,
		Balance:  balance,
		UserID:   userID,
	}, nil
}

func (a *Account) Deposit(amount string) error {

	amt, _, err := big.ParseFloat(amount, 10, 2, big.ToZero)
	if err != nil {
		return fmt.Errorf("invalid amount")
	}

	a.Balance.Add(&a.Balance, amt)
	return nil
}

func (a *Account) Withdraw(amount string) error {

	amt, _, err := big.ParseFloat(amount, 10, 2, big.ToZero)
	if err != nil {
		return fmt.Errorf("invalid amount")
	}

	if a.Balance.Cmp(amt) == -1 {
		return fmt.Errorf("insufficient balance")
	}

	a.Balance.Sub(&a.Balance, amt)
	return nil
}

func (a *Account) Transfer(to *Account, amount string) error {
	if a.Currency != to.Currency {
		return fmt.Errorf("cannot transfer between different currencies")
	}

	if err := a.Withdraw(amount); err != nil {
		return err
	}

	if err := to.Deposit(amount); err != nil {
		return err
	}

	return nil
}

func isValidCurrencyCode(code CurrencyCode) bool {
	_, ok := validCurrencies[code]
	return ok
}
