package service

import (
	"context"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
	"github.com/makifdb/mini-bank/minimalist/internal/repository"
)

type FeeService struct {
	feeRepo *repository.FeeRepository
}

func NewFeeService(feeRepo *repository.FeeRepository) *FeeService {
	return &FeeService{
		feeRepo: feeRepo,
	}
}

func (s *FeeService) CreateFee(ctx context.Context, amount float64, feeType models.FeeType, currency models.CurrencyCode) (*models.Fee, error) {
	f, err := models.NewFee(amount, feeType, currency)
	if err != nil {
		return nil, err
	}
	if err := s.feeRepo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeeService) GetFees(ctx context.Context, limit, offset int) ([]models.Fee, error) {
	return s.feeRepo.FindAll(ctx, limit, offset)
}

func (s *FeeService) UpdateFee(ctx context.Context, id int64, amount float64, feeType models.FeeType, currency models.CurrencyCode) (*models.Fee, error) {
	f, err := s.feeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	f.Amount = amount
	f.Type = feeType
	f.Currency = currency
	if err := s.feeRepo.Update(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeeService) DeleteFee(ctx context.Context, id int64) error {
	return s.feeRepo.Delete(ctx, id)
}
