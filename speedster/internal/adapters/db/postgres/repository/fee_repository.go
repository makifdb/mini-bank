package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/internal/core/domain"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
)

type FeeRepository struct {
	db *pgxpool.Pool
}

func NewFeeRepository(db *pgxpool.Pool) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) Create(ctx context.Context, t *domain.Fee) error {
	query := `INSERT INTO fees (id, external_id, created_at, updated_at, deleted_at, amount, type, currency) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, t.ID, t.ExternalID, t.CreatedAt, t.UpdatedAt, t.DeletedAt, t.Amount.String(), t.Type, t.Currency)
	return err
}

func (r *FeeRepository) FindByID(ctx context.Context, id string) (*domain.Fee, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees WHERE external_id = $1`
	row := r.db.QueryRow(ctx, query, id)
	f := &domain.Fee{}
	var amountStr string
	err := row.Scan(&f.ID, &f.ExternalID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &amountStr, &f.Type, &f.Currency)
	if err != nil {
		return nil, err
	}
	f.Amount, err = utils.NewBigDecimal(amountStr)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *FeeRepository) Update(ctx context.Context, t *domain.Fee) error {
	query := `UPDATE fees SET updated_at = $1, amount = $2, type = $3, currency = $4 
              WHERE external_id = $5`
	_, err := r.db.Exec(ctx, query, t.UpdatedAt, t.Amount.String(), t.Type, t.Currency, t.ExternalID)
	return err
}

func (r *FeeRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM fees WHERE external_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *FeeRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Fee, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fees := []domain.Fee{}
	for rows.Next() {
		f := domain.Fee{}
		var amountStr string
		if err := rows.Scan(&f.ID, &f.ExternalID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &amountStr, &f.Type, &f.Currency); err != nil {
			return nil, err
		}
		f.Amount, err = utils.NewBigDecimal(amountStr)
		if err != nil {
			return nil, err
		}
		fees = append(fees, f)
	}
	return fees, nil
}

func (r *FeeRepository) FindAllByType(ctx context.Context, feeType domain.FeeType, limit, offset int) ([]domain.Fee, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, amount, type, currency 
              FROM fees WHERE type = $1 LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, feeType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fees := []domain.Fee{}
	for rows.Next() {
		f := domain.Fee{}
		var amountStr string
		if err := rows.Scan(&f.ID, &f.ExternalID, &f.CreatedAt, &f.UpdatedAt, &f.DeletedAt, &amountStr, &f.Type, &f.Currency); err != nil {
			return nil, err
		}
		f.Amount, err = utils.NewBigDecimal(amountStr)
		if err != nil {
			return nil, err
		}
		fees = append(fees, f)
	}
	return fees, nil
}

func (r *FeeRepository) CalculateFee(ctx context.Context, feeType domain.FeeType, currency domain.CurrencyCode, amount *utils.BigDecimal) (*utils.BigDecimal, error) {
	query := `SELECT amount FROM fees WHERE type = $1 AND currency = $2`
	row := r.db.QueryRow(ctx, query, feeType, currency)
	var feeAmountStr string
	err := row.Scan(&feeAmountStr)
	if err != nil {
		return nil, err
	}
	feeAmount, err := utils.NewBigDecimal(feeAmountStr)
	if err != nil {
		return nil, err
	}

	return amount.Mul(feeAmount), nil
}
