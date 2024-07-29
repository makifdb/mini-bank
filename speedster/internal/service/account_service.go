package service

import (
	"context"

	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
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

func (s *AccountService) CreateAccount(ctx context.Context, userID string, currency models.CurrencyCode, amount string) (*models.Account, error) {
	acc, err := models.NewAccount(currency, amount, userID)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return s.accountRepo.FindByID(ctx, id)
}

func (s *AccountService) UpdateAccount(ctx context.Context, id string, currency models.CurrencyCode, balance string) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	acc.Currency = currency
	acc.Balance, err = utils.NewBigDecimal(balance)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID string, limit, offset int) ([]models.Account, error) {
	return s.accountRepo.FindAllByUserID(ctx, userID, limit, offset)
}

func (s *AccountService) DeleteAccount(ctx context.Context, id string) error {
	return s.accountRepo.Delete(ctx, id)
}

func (s *AccountService) Deposit(ctx context.Context, id string, amount string) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amt, err := utils.NewBigDecimal(amount)
	if err != nil {
		return nil, err
	}

	acc.Balance = acc.Balance.Add(amt)
	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	err = s.mailService.SendDepositNotification(user.Email, amount)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *AccountService) Withdraw(ctx context.Context, id string, amount string) (*models.Account, error) {
	acc, err := s.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	amt, err := utils.NewBigDecimal(amount)
	if err != nil {
		return nil, err
	}

	if err := acc.Withdraw(amt.String()); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Update(ctx, acc); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, acc.UserID)
	if err != nil {
		return nil, err
	}

	err = s.mailService.SendWithdrawalNotification(user.Email, amount)
	if err != nil {
		return nil, err
	}

	return acc, nil
}
