package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Create(ctx context.Context, admin *models.Admin) error {
	query := `INSERT INTO admins (id, external_id, created_at, updated_at, deleted_at, email) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, admin.ID, admin.ExternalID, admin.CreatedAt, admin.UpdatedAt, admin.DeletedAt, admin.Email)
	return err
}

func (r *AdminRepository) FindByID(ctx context.Context, id string) (*models.Admin, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, email FROM admins WHERE external_id = $1`
	row := r.db.QueryRow(ctx, query, id)
	admin := &models.Admin{}
	err := row.Scan(&admin.ID, &admin.ExternalID, &admin.CreatedAt, &admin.UpdatedAt, &admin.DeletedAt, &admin.Email)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, email FROM admins WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	admin := &models.Admin{}
	err := row.Scan(&admin.ID, &admin.ExternalID, &admin.CreatedAt, &admin.UpdatedAt, &admin.DeletedAt, &admin.Email)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (r *AdminRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Admin, error) {
	query := `SELECT id, external_id, created_at, updated_at, deleted_at, email FROM admins LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	admins := []models.Admin{}
	for rows.Next() {
		admin := models.Admin{}
		if err := rows.Scan(&admin.ID, &admin.ExternalID, &admin.CreatedAt, &admin.UpdatedAt, &admin.DeletedAt, &admin.Email); err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

func (r *AdminRepository) Update(ctx context.Context, admin *models.Admin) error {
	query := `UPDATE admins SET updated_at = $1, email = $2 WHERE external_id = $3`
	_, err := r.db.Exec(ctx, query, admin.UpdatedAt, admin.Email, admin.ExternalID)
	return err
}

func (r *AdminRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM admins WHERE external_id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
