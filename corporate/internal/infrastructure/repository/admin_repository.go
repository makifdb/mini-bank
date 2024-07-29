package repository

import (
	"context"

	"github.com/makifdb/mini-bank/corporate/internal/domain/admin"
	"gorm.io/gorm"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *admin.Admin) error
	FindByEmail(ctx context.Context, email string) (*admin.Admin, error)
	Delete(ctx context.Context, id string) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(ctx context.Context, admin *admin.Admin) error {
	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *adminRepository) FindByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	var admin admin.Admin
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error
	return &admin, err
}

func (r *adminRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&admin.Admin{}, id).Error
}
