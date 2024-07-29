package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (id, external_id, created_at, updated_at, deleted_at, first_name, last_name, email) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, u.ID, u.ExternalID, u.CreatedAt, u.UpdatedAt, u.DeletedAt, u.FirstName, u.LastName, u.Email)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, first_name, last_name, email FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.ExternalID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.FirstName, &u.LastName, &u.Email)
	return u, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, first_name, last_name, email FROM users WHERE external_id = $1`
	row := r.db.QueryRow(ctx, query, id)
	u := &models.User{}
	err := row.Scan(&u.ID, &u.ExternalID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.FirstName, &u.LastName, &u.Email)
	return u, err
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET updated_at = $2, deleted_at = $3, first_name = $4, last_name = $5, email = $6 WHERE external_id = $1`
	_, err := r.db.Exec(ctx, query, u.ID, u.UpdatedAt, u.DeletedAt, u.FirstName, u.LastName, u.Email)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE external_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, first_name, last_name, email FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		u := models.User{}
		if err := rows.Scan(&u.ID, &u.ExternalID, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt, &u.FirstName, &u.LastName, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) FindByIDWithAccounts(ctx context.Context, id string) (*models.User, error) {
	user, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := `SELECT id, external_id, created_at, updated_at, deleted_at, balance, currency, user_id FROM accounts WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		acc := models.Account{}
		if err := rows.Scan(&acc.ID, &acc.ExternalID, &acc.CreatedAt, &acc.UpdatedAt, &acc.DeletedAt, &acc.Balance, &acc.Currency, &acc.UserID); err != nil {
			return nil, err
		}
		user.Accounts = append(user.Accounts, acc)
	}

	return user, nil
}
