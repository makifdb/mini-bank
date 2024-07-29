// internal/infrastructure/repository/account_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *account.Account) error {
	return r.db.WithContext(ctx).Create(acc).Error
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	var acc account.Account
	err := r.db.WithContext(ctx).First(&acc, id).Error
	return &acc, err
}

func (r *AccountRepository) FindByIDTx(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*account.Account, error) {
	var acc account.Account
	err := tx.WithContext(ctx).First(&acc, id).Error
	return &acc, err
}

func (r *AccountRepository) Update(ctx context.Context, acc *account.Account) error {
	return r.db.WithContext(ctx).Save(acc).Error
}

func (r *AccountRepository) UpdateTx(ctx context.Context, tx *gorm.DB, acc *account.Account) error {
	return tx.WithContext(ctx).Save(acc).Error
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&account.Account{}, id).Error
}

func (r *AccountRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]account.Account, error) {
	var accounts []account.Account
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) DB() *gorm.DB {
	return r.db
}
