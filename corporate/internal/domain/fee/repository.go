package fee

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/shopspring/decimal"
)

type FeeRepository interface {
	Create(ctx context.Context, t *Fee) error
	FindByID(ctx context.Context, id uuid.UUID) (*Fee, error)
	Update(ctx context.Context, t *Fee) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int) ([]Fee, error)
	FindAllByType(ctx context.Context, feeType Type, limit, offset int) ([]Fee, error)
	CalculateFee(ctx context.Context, feeType Type, currency account.CurrencyCode, amount decimal.Decimal) (decimal.Decimal, error)
}
