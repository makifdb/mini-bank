package utils

import (
	"errors"

	"github.com/cockroachdb/apd/v3"
)

type BigDecimal struct {
	value *apd.Decimal
	ctx   *apd.Context
}

func NewBigDecimal(value string) (*BigDecimal, error) {
	d, _, err := apd.NewFromString(value)
	if err != nil {
		return nil, err
	}
	return &BigDecimal{value: d, ctx: apd.BaseContext.WithPrecision(34)}, nil
}

func (bd *BigDecimal) Add(b *BigDecimal) *BigDecimal {
	result := new(apd.Decimal)
	_, err := bd.ctx.Add(result, bd.value, b.value)
	if err != nil {
		return nil
	}
	return &BigDecimal{value: result, ctx: bd.ctx}
}

func (bd *BigDecimal) Sub(b *BigDecimal) *BigDecimal {
	result := new(apd.Decimal)
	_, err := bd.ctx.Sub(result, bd.value, b.value)
	if err != nil {
		return nil
	}
	return &BigDecimal{value: result, ctx: bd.ctx}
}

func (bd *BigDecimal) Mul(b *BigDecimal) *BigDecimal {
	result := new(apd.Decimal)
	_, err := bd.ctx.Mul(result, bd.value, b.value)
	if err != nil {
		return nil
	}
	return &BigDecimal{value: result, ctx: bd.ctx}
}

func (bd *BigDecimal) Div(b *BigDecimal) (*BigDecimal, error) {
	if b.value.IsZero() {
		return nil, errors.New("division by zero")
	}
	result := new(apd.Decimal)
	_, err := bd.ctx.Quo(result, bd.value, b.value)
	if err != nil {
		return nil, err
	}
	return &BigDecimal{value: result, ctx: bd.ctx}, nil
}

func (bd *BigDecimal) Cmp(b *BigDecimal) int {
	return bd.value.Cmp(b.value)
}

func (bd *BigDecimal) String() string {
	return bd.value.Text('f')
}

// Getter method for value
func (bd *BigDecimal) Value() *apd.Decimal {
	return bd.value
}
