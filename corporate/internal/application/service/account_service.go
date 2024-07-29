// internal/application/service/account_service.go
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/fee"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/email"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
	"github.com/shopspring/decimal"
)

type AccountService struct {
	accountRepo repository.AccountRepository
	userRepo    repository.UserRepository
	feeRepo     repository.FeeRepository
	mailService email.MailService
}

func NewAccountService(accountRepo repository.AccountRepository, userRepo repository.UserRepository, feeRepo repository.FeeRepository, mailService email.MailService) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		feeRepo:     feeRepo,
		mailService: mailService,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID uuid.UUID, currency account.CurrencyCode, amount decimal.Decimal) (*account.Account, error) {
	acc, err := account.NewAccount(currency, amount, userID)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	return s.accountRepo.FindByID(ctx, id)
}

func (s *AccountService) UpdateAccount(ctx context.Context, id uuid.UUID, currency account.CurrencyCode, balance decimal.Decimal) (*account.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	acc.Currency = currency
	acc.Balance = balance

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]account.Account, error) {
	accounts, err := s.accountRepo.FindAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *AccountService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return s.accountRepo.Delete(ctx, id)
}

func (s *AccountService) Deposit(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*account.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	feeAmount, err := s.feeRepo.CalculateFee(ctx, fee.Deposit, acc.Currency, amount)
	if err != nil {
		return nil, err
	}

	if err := acc.Deposit(amount.Sub(feeAmount)); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	if err := s.mailService.SendDepositNotification(user.Email, amount); err != nil {
		return nil, err
	}

	return acc, nil

}

func (s *AccountService) Withdraw(ctx context.Context, id uuid.UUID, amount decimal.Decimal) (*account.Account, error) {

	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	feeAmount, err := s.feeRepo.CalculateFee(ctx, fee.Withdraw, acc.Currency, amount)
	if err != nil {
		return nil, err
	}

	if err := acc.Withdraw(amount.Add(feeAmount)); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	if err := s.mailService.SendWithdrawalNotification(user.Email, amount); err != nil {
		return nil, err
	}

	return acc, nil
}
