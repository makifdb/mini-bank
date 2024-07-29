package admin

import (
	"context"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *Admin) error
	FindByEmail(ctx context.Context, email string) (*Admin, error)
	Delete(ctx context.Context, id string) error
}
