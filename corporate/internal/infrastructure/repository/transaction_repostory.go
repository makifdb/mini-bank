package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/transaction"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, txn *transaction.Transaction) error {
	return r.db.WithContext(ctx).Create(txn).Error
}

func (r *TransactionRepository) CreateTx(ctx context.Context, tx *gorm.DB, txn *transaction.Transaction) error {
	return tx.WithContext(ctx).Create(txn).Error
}

func (r *TransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	var txn transaction.Transaction
	err := r.db.WithContext(ctx).First(&txn, id).Error
	return &txn, err
}
