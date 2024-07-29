package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *models.Account) error {
	query := `INSERT INTO accounts (id, balance, currency, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, acc.ID, acc.Balance.String(), acc.Currency, acc.UserID, acc.CreatedAt, acc.UpdatedAt)
	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (*models.Account, error) {
	query := `SELECT id, balance, currency, user_id, created_at, updated_at FROM accounts WHERE external_id = $1`
	row := r.db.QueryRow(ctx, query, id)

	acc := &models.Account{}
	var balanceStr string

	err := row.Scan(&acc.ID, &balanceStr, &acc.Currency, &acc.UserID, &acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	acc.Balance, err = utils.NewBigDecimal(balanceStr)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (r *AccountRepository) Update(ctx context.Context, acc *models.Account) error {
	query := `UPDATE accounts SET balance = $1, currency = $2, updated_at = $3 WHERE external_id = $4`
	_, err := r.db.Exec(ctx, query, acc.Balance.String(), acc.Currency, acc.UpdatedAt, acc.ID)
	return err
}

func (r *AccountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE external_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *AccountRepository) FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Account, error) {
	query := `SELECT id, balance, currency, user_id, created_at, updated_at FROM accounts WHERE user_id = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
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

		acc.Balance, err = utils.NewBigDecimal(balanceStr)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}
