package account

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, acc *Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*Account, error)
	Update(ctx context.Context, acc *Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAllByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Account, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	DB() *gorm.DB
}
