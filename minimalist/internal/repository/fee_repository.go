package repository

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
)

type FeeRepository struct {
	db *sql.DB
}

func NewFeeRepository(db *sql.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) Create(ctx context.Context, t *models.Fee) error {
	query := `INSERT INTO fees (id, created_at, updated_at, deleted_at, amount, type, currency) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, t.ID, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Amount, t.Type, t.Currency)
	return err
}

func (r *FeeRepository) FindByID(ctx context.Context, id int64) (*models.Fee, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	f := &models.Fee{}
	var deletedAt sql.NullTime
	err := row.Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt, &deletedAt, &f.Amount, &f.Type, &f.Currency)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		f.DeletedAt = &deletedAt.Time
	}
	return f, nil
}

func (r *FeeRepository) Update(ctx context.Context, t *models.Fee) error {
	query := `UPDATE fees SET updated_at = ?, amount = ?, type = ?, currency = ? 
              WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, t.UpdatedAt, t.Amount, t.Type, t.Currency, t.ID)
	return err
}

func (r *FeeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM fees WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *FeeRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Fee, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fees := []models.Fee{}
	for rows.Next() {
		f := models.Fee{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt, &deletedAt, &f.Amount, &f.Type, &f.Currency); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			f.DeletedAt = &deletedAt.Time
		}
		fees = append(fees, f)
	}
	return fees, nil
}

func (r *FeeRepository) FindAllByType(ctx context.Context, feeType models.FeeType, limit, offset int) ([]models.Fee, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees WHERE type = ? LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, feeType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fees := []models.Fee{}
	for rows.Next() {
		f := models.Fee{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt, &deletedAt, &f.Amount, &f.Type, &f.Currency); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			f.DeletedAt = &deletedAt.Time
		}
		fees = append(fees, f)
	}
	return fees, nil
}

func (r *FeeRepository) CalculateFee(ctx context.Context, feeType models.FeeType, currency models.CurrencyCode, amount float64) (*big.Float, error) {
	query := `SELECT amount FROM fees WHERE type = ? AND currency = ?`
	row := r.db.QueryRowContext(ctx, query, feeType, currency)
	var feeAmountStr string
	err := row.Scan(&feeAmountStr)
	if err != nil {
		return nil, err
	}
	feeAmount, _, err := big.ParseFloat(feeAmountStr, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, err
	}

	amountFloat := big.NewFloat(amount)
	result := new(big.Float).Mul(amountFloat, feeAmount)
	return result, nil
}
