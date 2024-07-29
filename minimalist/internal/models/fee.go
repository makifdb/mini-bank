package models

import (
	"fmt"
)

// FeeType represents the type of a fee.
type FeeType string

const (
	Transfer FeeType = "transfer"
	Deposit  FeeType = "deposit"
	Withdraw FeeType = "withdraw"
)

// Fee represents a fee applied to transactions.
type Fee struct {
	Base
	Amount   float64      `json:"amount"`
	Type     FeeType      `json:"type"`
	Currency CurrencyCode `json:"currency"`
}

func NewFee(amount float64, feeType FeeType, currency CurrencyCode) (*Fee, error) {
	if !isValidCurrencyCode(currency) {
		return nil, fmt.Errorf("invalid currency code")
	}
	if !isValidFeeType(feeType) {
		return nil, fmt.Errorf("invalid fee type")
	}

	return &Fee{
		Base:     NewBase(),
		Amount:   amount,
		Type:     feeType,
		Currency: currency,
	}, nil
}

func isValidFeeType(feeType FeeType) bool {
	switch feeType {
	case Transfer, Deposit, Withdraw:
		return true
	}
	return false
}
