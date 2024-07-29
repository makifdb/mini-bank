package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/fee"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
	"github.com/shopspring/decimal"
)

type FeeService struct {
	feeRepo *repository.FeeRepository
}

func NewFeeService(feeRepo *repository.FeeRepository) *FeeService {
	return &FeeService{
		feeRepo: feeRepo,
	}
}

func (s *FeeService) CreateFee(ctx context.Context, amount decimal.Decimal, feeType fee.Type, currency account.CurrencyCode) (*fee.Fee, error) {
	f, err := fee.NewFee(amount, feeType, currency)
	if err != nil {
		return nil, err
	}
	if err := s.feeRepo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeeService) GetFees(ctx context.Context, limit, offset int) ([]fee.Fee, error) {
	return s.feeRepo.FindAll(ctx, limit, offset)
}

func (s *FeeService) UpdateFee(ctx context.Context, id uuid.UUID, amount decimal.Decimal, feeType fee.Type, currency account.CurrencyCode) (*fee.Fee, error) {
	f, err := s.feeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	f.Amount = amount
	f.Type = feeType
	f.CurrencyCode = currency
	if err := s.feeRepo.Update(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeeService) DeleteFee(ctx context.Context, id uuid.UUID) error {
	return s.feeRepo.Delete(ctx, id)
}
