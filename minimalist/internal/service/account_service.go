package service

import (
	"context"
	"fmt"
	"math/big"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
	"github.com/makifdb/mini-bank/minimalist/internal/repository"
	"github.com/makifdb/mini-bank/minimalist/pkg/utils"
)

type AccountService struct {
	accountRepo *repository.AccountRepository
	userRepo    *repository.UserRepository
	feeRepo     *repository.FeeRepository
	mailService *utils.MailService
}

func NewAccountService(accountRepo *repository.AccountRepository, userRepo *repository.UserRepository, feeRepo *repository.FeeRepository, mailService *utils.MailService) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		userRepo:    userRepo,
		feeRepo:     feeRepo,
		mailService: mailService,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64, amount float64, currency models.CurrencyCode) (*models.Account, error) {
	amountFloat := big.NewFloat(amount)
	acc, err := models.NewAccount(currency, *amountFloat, userID)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id int64) (*models.Account, error) {
	return s.accountRepo.FindByID(ctx, id)
}

func (s *AccountService) UpdateAccount(ctx context.Context, id int64, currency models.CurrencyCode, balance float64) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	balanceFloat := big.NewFloat(balance)
	acc.Currency = currency
	acc.Balance = *balanceFloat

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID int64, limit, offset int) ([]models.Account, error) {
	return s.accountRepo.FindAllByUserID(ctx, userID, limit, offset)
}

func (s *AccountService) DeleteAccount(ctx context.Context, id int64) error {
	return s.accountRepo.Delete(ctx, id)
}

func (s *AccountService) Deposit(ctx context.Context, id int64, amount float64) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amt := big.NewFloat(amount)
	acc.Balance.Add(&acc.Balance, amt)
	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	err = s.mailService.SendDepositNotification(user.Email, fmt.Sprintf("%f", amount))
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) Withdraw(ctx context.Context, id int64, amount float64) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amt := big.NewFloat(amount)
	amtStr := amt.Text('f', -1)
	if err := acc.Withdraw(amtStr); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	err = s.mailService.SendWithdrawalNotification(user.Email, fmt.Sprintf("%f", amount))
	if err != nil {
		return nil, err
	}

	return acc, nil
}
