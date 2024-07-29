package fee

import (
	"fmt"

	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/base"
	"github.com/shopspring/decimal"
)

// FeeType represents the type of a fee.
type Type string

const (
	Transfer Type = "transfer"
	Deposit  Type = "deposit"
	Withdraw Type = "withdraw"
)

var validFeeTypes = map[Type]struct{}{
	Transfer: {},
	Deposit:  {},
	Withdraw: {},
}

type Fee struct {
	base.Base
	Amount       decimal.Decimal      `json:"amount"`
	Type         Type                 `json:"type"`
	CurrencyCode account.CurrencyCode `json:"currency_code"`
}

// NewFee creates a new fee.
func NewFee(amount decimal.Decimal, feeType Type, currencyCode account.CurrencyCode) (*Fee, error) {

	// check if the currency code is valid
	if !account.IsCurrencyCodeValid(currencyCode) {
		return nil, fmt.Errorf("invalid currency code")
	}

	// check if the fee type is valid
	if !IsFeeTypeValid(feeType) {
		return nil, fmt.Errorf("invalid fee type")
	}

	return &Fee{
		Base:         base.NewBase(),
		Amount:       amount,
		Type:         feeType,
		CurrencyCode: currencyCode,
	}, nil
}

func IsFeeTypeValid(feeType Type) bool {
	_, ok := validFeeTypes[feeType]
	return ok
}
