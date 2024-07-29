package service

import (
	"context"
	"fmt"
	"os"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/pkg/models"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
	"github.com/redis/go-redis/v9"
)

type AdminService struct {
	adminRepo   *repository.AdminRepository
	redisClient *redis.Client
	mailService *utils.MailService
}

func NewAdminService(adminRepo *repository.AdminRepository, redisClient *redis.Client, mailService *utils.MailService) *AdminService {
	return &AdminService{
		adminRepo:   adminRepo,
		redisClient: redisClient,
		mailService: mailService,
	}
}

func (s *AdminService) SignUp(ctx context.Context, email string) (*models.Admin, error) {
	admin, err := models.NewAdmin(email)
	if err != nil {
		return nil, err
	}

	if err := s.adminRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	verifyCode, err := utils.GenerateVerificationCode()
	if err != nil {
		return nil, err
	}

	if err := s.redisClient.Set(ctx, "login:"+admin.Email, verifyCode, 10*time.Minute).Err(); err != nil {
		return nil, err
	}

	if err := s.mailService.SendVerificationEmail(admin.Email, verifyCode); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AdminService) Login(ctx context.Context, email string) error {
	admin, err := s.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	verifyCode, err := utils.GenerateVerificationCode()
	if err != nil {
		return err
	}

	if err := s.redisClient.Set(ctx, "login:"+admin.Email, verifyCode, 10*time.Minute).Err(); err != nil {
		return err
	}

	if err := s.mailService.SendLoginEmail(admin.Email, verifyCode); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) VerifyAdmin(ctx context.Context, email, code string) (string, error) {
	storedCode, err := s.redisClient.Get(ctx, "login:"+email).Result()
	if err != nil {
		return "", err
	}

	if code != storedCode {
		return "", fmt.Errorf("invalid verification code")
	}

	admin, err := s.adminRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := s.redisClient.Del(ctx, "login:"+email).Err(); err != nil {
		return "", err
	}

	tokenString, err := generateJWTToken(admin.ExternalID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AdminService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	newTokenString, err := generateJWTToken(claims["sub"].(string))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}

func generateJWTToken(adminID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "mini-bank",
		"sub": adminID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET")))
}
