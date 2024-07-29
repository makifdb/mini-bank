package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/internal/core/domain"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, txn *domain.Transaction) error {
	query := `INSERT INTO transactions (id, external_id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, txn.ID, txn.ExternalID, txn.CreatedAt, txn.UpdatedAt, txn.DeletedAt, txn.FromAccountID, txn.ToAccountID, txn.Amount, txn.Fee)
	return err
}

func (r *TransactionRepository) FindByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee FROM transactions WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	txn := &domain.Transaction{}
	err := row.Scan(&txn.ID, &txn.ExternalID, &txn.CreatedAt, &txn.UpdatedAt, &txn.DeletedAt, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.Fee)
	return txn, err
}

func (r *TransactionRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, from_account_id, to_account_id, amount, fee FROM transactions LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []domain.Transaction{}
	for rows.Next() {
		txn := domain.Transaction{}
		if err := rows.Scan(&txn.ID, &txn.ExternalID, &txn.CreatedAt, &txn.UpdatedAt, &txn.DeletedAt, &txn.FromAccountID, &txn.ToAccountID, &txn.Amount, &txn.Fee); err != nil {
			return nil, err
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}
