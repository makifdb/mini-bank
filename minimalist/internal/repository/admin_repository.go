package repository

import (
	"context"
	"database/sql"

	"github.com/makifdb/mini-bank/minimalist/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	query := `INSERT INTO admins (id, created_at, updated_at, deleted_at, email) 
              VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, admin.ID, admin.CreatedAt, admin.UpdatedAt, admin.DeletedAt, admin.Email)
	return err
}

func (r *AdminRepository) FindByID(ctx context.Context, id int64) (*models.Admin, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, email FROM admins WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	admin := &models.Admin{}
	var deletedAt sql.NullTime
	err := row.Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt, &deletedAt, &admin.Email)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		admin.DeletedAt = &deletedAt.Time
	}
	return admin, nil
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, email FROM admins WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)
	admin := &models.Admin{}
	var deletedAt sql.NullTime
	err := row.Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt, &deletedAt, &admin.Email)
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		admin.DeletedAt = &deletedAt.Time
	}
	return admin, nil
}

func (r *AdminRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Admin, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, email FROM admins LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := []models.Admin{}
	for rows.Next() {
		admin := models.Admin{}
		var deletedAt sql.NullTime
		if err := rows.Scan(&admin.ID, &admin.CreatedAt, &admin.UpdatedAt, &deletedAt, &admin.Email); err != nil {
			return nil, err
		}
		if deletedAt.Valid {
			admin.DeletedAt = &deletedAt.Time
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

func (r *AdminRepository) Update(ctx context.Context, admin *models.Admin) error {
	query := `UPDATE admins SET updated_at = ?, email = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, admin.UpdatedAt, admin.Email, admin.ID)
	return err
}

func (r *AdminRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM admins WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
