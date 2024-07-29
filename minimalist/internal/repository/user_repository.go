package repository

import (
	"context"
	"database/sql"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (id, created_at, updated_at, deleted_at, first_name, last_name, email) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.CreatedAt, u.UpdatedAt, u.DeletedAt, u.FirstName, u.LastName, u.Email)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, first_name, last_name, email 
              FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)
	u := &models.User{}
	var deletedAt sql.NullTime
	err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &deletedAt, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}
	return u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, first_name, last_name, email 
              FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	u := &models.User{}
	var deletedAt sql.NullTime
	err := row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &deletedAt, &u.FirstName, &u.LastName, &u.Email)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		u.DeletedAt = &deletedAt.Time
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET updated_at = ?, deleted_at = ?, first_name = ?, last_name = ?, email = ? 
              WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, u.UpdatedAt, u.DeletedAt, u.FirstName, u.LastName, u.Email, u.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, first_name, last_name, email 
              FROM users LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		u := models.User{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &deletedAt, &u.FirstName, &u.LastName, &u.Email); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			u.DeletedAt = &deletedAt.Time
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) FindByIDWithAccounts(ctx context.Context, id int64) (*models.User, error) {
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, created_at, updated_at, deleted_at, balance, currency, user_id 
              FROM accounts WHERE user_id = ?`
	rows, err := r.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		acc := models.Account{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&acc.ID, &acc.CreatedAt, &acc.UpdatedAt, &deletedAt, &acc.Balance, &acc.Currency, &acc.UserID); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			acc.DeletedAt = &deletedAt.Time
		}
		user.Accounts = append(user.Accounts, acc)
	}

	return user, nil
}
