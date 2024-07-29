package service

import (
	"context"
	"time"

	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, firstName, lastName, email string) (*models.User, error) {
	user, err := models.NewUser(firstName, lastName, email)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, firstName, lastName, email string) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	return s.userRepo.FindAll(ctx, limit, offset)
}

func (s *UserService) GetUserByIDWithAccounts(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByIDWithAccounts(ctx, id)
}
