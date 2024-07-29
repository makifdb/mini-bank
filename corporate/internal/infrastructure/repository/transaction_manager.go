package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	Begin(ctx context.Context) (*gorm.DB, error)
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
}

type GormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{db: db}
}

func (m *GormTransactionManager) Begin(ctx context.Context) (*gorm.DB, error) {
	return m.db.WithContext(ctx).Begin(), nil
}

func (m *GormTransactionManager) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (m *GormTransactionManager) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
