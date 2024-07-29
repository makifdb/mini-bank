package service

import (
	"context"

	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

type FeeService struct {
	feeRepo *repository.FeeRepository
}

func NewFeeService(feeRepo *repository.FeeRepository) *FeeService {
	return &FeeService{
		feeRepo: feeRepo,
	}
}

func (s *FeeService) CreateFee(ctx context.Context, amount string, feeType models.FeeType, currency models.CurrencyCode) (*models.Fee, error) {
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

func (s *FeeService) UpdateFee(ctx context.Context, id, amount string, feeType models.FeeType, currency models.CurrencyCode) (*models.Fee, error) {
	f, err := s.feeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	a, err := utils.NewBigDecimal(amount)
	if err != nil {
		return nil, err
	}

	f.Amount = a
	f.Type = feeType
	f.Currency = currency
	if err := s.feeRepo.Update(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *FeeService) DeleteFee(ctx context.Context, id string) error {
	return s.feeRepo.Delete(ctx, id)
}
