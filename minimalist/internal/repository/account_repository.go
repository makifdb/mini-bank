package repository

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *models.Account) error {
	query := `INSERT INTO accounts (id, balance, currency, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, acc.ID, acc.Balance.String(), acc.Currency, acc.UserID, acc.CreatedAt, acc.UpdatedAt)
	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id int64) (*models.Account, error) {
	query := `SELECT id, balance, currency, user_id, created_at, updated_at FROM accounts WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	acc := &models.Account{}
	var balanceStr string

	err := row.Scan(&acc.ID, &balanceStr, &acc.Currency, &acc.UserID, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	balance, _, err := big.ParseFloat(balanceStr, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, err
	}
	acc.Balance = *balance

	return acc, nil
}

func (r *AccountRepository) Update(ctx context.Context, acc *models.Account) error {
	query := `UPDATE accounts SET balance = ?, currency = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, acc.Balance.String(), acc.Currency, acc.UpdatedAt, acc.ID)
	return err
}

func (r *AccountRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM accounts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *AccountRepository) FindAllByUserID(ctx context.Context, userID int64, limit, offset int) ([]models.Account, error) {
	query := `SELECT id, balance, currency, user_id, created_at, updated_at FROM accounts WHERE user_id = ? LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account

	for rows.Next() {
		acc := models.Account{}
		var balanceStr string

		err := rows.Scan(&acc.ID, &balanceStr, &acc.Currency, &acc.UserID, &acc.CreatedAt, &acc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		balance, _, err := big.ParseFloat(balanceStr, 10, 0, big.ToNearestEven)
		if err != nil {
			return nil, err
		}
		acc.Balance = *balance

		accounts = append(accounts, acc)
	}

	return accounts, nil
}
