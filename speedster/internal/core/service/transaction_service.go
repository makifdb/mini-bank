package service

import (
	"context"
	"fmt"

	"github.com/makifdb/mini-bank/speedster/internal/adapters/db/postgres/repository"
	"github.com/makifdb/mini-bank/speedster/internal/core/domain"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.AccountRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, accountRepo *repository.AccountRepository) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (s *TransactionService) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount, fee *utils.BigDecimal) error {
	fromAccount, err := s.accountRepo.FindByID(ctx, fromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.FindByID(ctx, toAccountID)
	if err != nil {
		return err
	}

	if fromAccount.Balance.Cmp(amount.Add(fee)) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	fromAccount.Balance = fromAccount.Balance.Sub(amount.Add(fee))
	toAccount.Balance = toAccount.Balance.Add(amount)

	if err := s.accountRepo.Update(ctx, fromAccount); err != nil {
		return err
	}

	if err := s.accountRepo.Update(ctx, toAccount); err != nil {
		return err
	}

	txn, err := domain.NewTransaction(fromAccountID, toAccountID, amount, fee)
	if err != nil {
		return err
	}

	return s.transactionRepo.Create(ctx, txn)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	return s.transactionRepo.FindByID(ctx, id)
}

func (s *TransactionService) GetTransactions(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	return s.transactionRepo.FindAll(ctx, limit, offset)
}
