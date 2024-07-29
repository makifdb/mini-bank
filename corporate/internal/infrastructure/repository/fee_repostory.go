package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/fee"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type FeeRepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) Create(ctx context.Context, t *fee.Fee) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *FeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*fee.Fee, error) {
	var f fee.Fee
	err := r.db.WithContext(ctx).First(&f, id).Error
	return &f, err
}

func (r *FeeRepository) Update(ctx context.Context, t *fee.Fee) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *FeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&fee.Fee{}, id).Error
}

func (r *FeeRepository) FindAll(ctx context.Context, limit, offset int) ([]fee.Fee, error) {
	var fees []fee.Fee
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&fees).Error
	return fees, err
}

func (r *FeeRepository) FindAllByType(ctx context.Context, feeType fee.Type, limit, offset int) ([]fee.Fee, error) {
	var fees []fee.Fee
	err := r.db.WithContext(ctx).Where("type = ?", feeType).Limit(limit).Offset(offset).Find(&fees).Error
	return fees, err
}

func (r *FeeRepository) CalculateFee(ctx context.Context, feeType fee.Type, currency account.CurrencyCode, amount decimal.Decimal) (decimal.Decimal, error) {
	var f fee.Fee
	err := r.db.WithContext(ctx).Where("type = ? AND currency_code = ?", feeType, currency).First(&f).Error
	if err != nil {
		return decimal.Zero, err
	}

	return amount.Mul(f.Amount), nil
}
