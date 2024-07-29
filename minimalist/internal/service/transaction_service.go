package service

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
	"github.com/makifdb/mini-bank/minimalist/internal/repository"
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

func (s *TransactionService) Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount, fee *big.Float) error {
	fromAccount, err := s.accountRepo.FindByID(ctx, fromAccountID)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.FindByID(ctx, toAccountID)
	if err != nil {
		return err
	}

	total := new(big.Float).Add(amount, fee)

	if fromAccount.Balance.Cmp(total) < 0 {
		return fmt.Errorf("insufficient balance")
	}

	fromAccount.Balance.Sub(&fromAccount.Balance, total)
	toAccount.Balance.Add(&toAccount.Balance, amount)

	if err := s.accountRepo.Update(ctx, fromAccount); err != nil {
		return err
	}

	if err := s.accountRepo.Update(ctx, toAccount); err != nil {
		return err
	}

	txn := &models.Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
		Fee:           fee,
	}

	return s.transactionRepo.Create(ctx, txn)
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error) {
	txn, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}
	return txn, nil
}

func (s *TransactionService) GetTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
