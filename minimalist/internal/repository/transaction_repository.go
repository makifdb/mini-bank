package repository

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, txn *models.Transaction) error {
	query := `INSERT INTO transactions (id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	amountStr := txn.Amount.Text('f', -1)
	feeStr := txn.Fee.Text('f', -1)
	_, err := r.db.ExecContext(ctx, query, txn.ID, txn.CreatedAt, txn.UpdatedAt, txn.DeletedAt, txn.FromAccountID, txn.ToAccountID, amountStr, feeStr)
	return err
}

func (r *TransactionRepository) FindByID(ctx context.Context, id int64) (*models.Transaction, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee 
              FROM transactions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	txn := &models.Transaction{}
	var amountStr, feeStr string
	var deletedAt sql.NullTime
	err := row.Scan(&txn.ID, &txn.CreatedAt, &txn.UpdatedAt, &deletedAt, &txn.FromAccountID, &txn.ToAccountID, &amountStr, &feeStr)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		txn.DeletedAt = &deletedAt.Time
	}
	txn.Amount, _, _ = big.ParseFloat(amountStr, 10, 0, big.ToNearestEven)
	txn.Fee, _, _ = big.ParseFloat(feeStr, 10, 0, big.ToNearestEven)
	return txn, nil
}

func (r *TransactionRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee 
              FROM transactions LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		txn := models.Transaction{}
		var amountStr, feeStr string
		var deletedAt sql.NullTime
		if err := rows.Scan(&txn.ID, &txn.CreatedAt, &txn.UpdatedAt, &deletedAt, &txn.FromAccountID, &txn.ToAccountID, &amountStr, &feeStr); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			txn.DeletedAt = &deletedAt.Time
		}
		txn.Amount, _, _ = big.ParseFloat(amountStr, 10, 0, big.ToNearestEven)
		txn.Fee, _, _ = big.ParseFloat(feeStr, 10, 0, big.ToNearestEven)
		transactions = append(transactions, txn)
	}
	return transactions, nil
}
