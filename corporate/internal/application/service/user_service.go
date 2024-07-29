package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/domain/user"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, firstName, lastName, email string) (*user.User, error) {
	u, err := user.NewUser(firstName, lastName, email)
	if err != nil {
		return nil, err
	}
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) GetUsers(ctx context.Context, limit, offset int) ([]user.User, error) {
	return s.userRepo.FindAll(ctx, limit, offset)
}

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName, email string) (*user.User, error) {
	u, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	u.FirstName = firstName
	u.LastName = lastName
	u.Email = email
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
