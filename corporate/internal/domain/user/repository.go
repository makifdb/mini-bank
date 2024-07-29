package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int) ([]User, error)
	FindByIDWithAccounts(ctx context.Context, id uuid.UUID) (*User, error)
}
