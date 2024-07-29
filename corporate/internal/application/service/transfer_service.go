package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/fee"
	"github.com/makifdb/mini-bank/corporate/internal/domain/transaction"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/email"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type TransferService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	userRepo        repository.UserRepository
	feeRepo         repository.FeeRepository
	mailService     email.MailService
	logger          *zap.Logger
	txManager       repository.TransactionManager
}

func NewTransferService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	userRepo repository.UserRepository,
	feeRepo repository.FeeRepository,
	mailService email.MailService,
	logger *zap.Logger,
	txManager repository.TransactionManager,
) *TransferService {
	return &TransferService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		userRepo:        userRepo,
		feeRepo:         feeRepo,
		mailService:     mailService,
		logger:          logger,
		txManager:       txManager,
	}
}

func (s *TransferService) Transfer(ctx context.Context, fromAccountID, toAccountID uuid.UUID, amount decimal.Decimal) error {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction", zap.Error(err))
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.txManager.Rollback(tx)
		}
	}()

	fromAccount, err := s.accountRepo.FindByIDTx(ctx, tx, fromAccountID)
	if err != nil {
		s.logger.Error("Failed to find from account", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	toAccount, err := s.accountRepo.FindByIDTx(ctx, tx, toAccountID)
	if err != nil {
		s.logger.Error("Failed to find to account", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if fromAccount.ID == toAccount.ID {
		s.txManager.Rollback(tx)
		return &transactionError{"Cannot transfer to the same account"}
	}

	if fromAccount.Currency != toAccount.Currency {
		s.txManager.Rollback(tx)
		return &transactionError{"Cannot transfer between different currencies"}
	}

	transactionFee, err := s.feeRepo.CalculateFee(ctx, fee.Transfer, fromAccount.Currency, amount)
	if err != nil {
		s.logger.Error("Failed to calculate fee", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if fromAccount.Balance.LessThan(amount.Add(transactionFee)) {
		s.txManager.Rollback(tx)
		return &transactionError{"Insufficient balance"}
	}

	if err := fromAccount.Transfer(toAccount, amount); err != nil {
		s.logger.Error("Failed to transfer funds", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := fromAccount.Withdraw(transactionFee); err != nil {
		s.logger.Error("Failed to withdraw fee", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := s.accountRepo.UpdateTx(ctx, tx, fromAccount); err != nil {
		s.logger.Error("Failed to update from account", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := s.accountRepo.UpdateTx(ctx, tx, toAccount); err != nil {
		s.logger.Error("Failed to update to account", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	txn, err := transaction.NewTransaction(fromAccountID, toAccountID, amount, transactionFee)
	if err != nil {
		s.logger.Error("Failed to create transaction", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := s.transactionRepo.CreateTx(ctx, tx, txn); err != nil {
		s.logger.Error("Failed to save transaction", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	fromUser, err := s.userRepo.FindByID(ctx, fromAccount.UserID)
	if err != nil {
		s.logger.Error("Failed to find from user", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	toUser, err := s.userRepo.FindByID(ctx, toAccount.UserID)
	if err != nil {
		s.logger.Error("Failed to find to user", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := s.mailService.SendTransferNotification(fromUser.Email, toUser.Email, amount); err != nil {
		s.logger.Error("Failed to send email notifications", zap.Error(err))
		s.txManager.Rollback(tx)
		return err
	}

	if err := s.txManager.Commit(tx); err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

type transactionError struct {
	message string
}

func (e *transactionError) Error() string {
	return e.message
}
