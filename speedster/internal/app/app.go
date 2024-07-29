package app

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/makifdb/mini-bank/speedster/internal/handler"
	"github.com/redis/go-redis/v9"

	"github.com/makifdb/mini-bank/speedster/internal/middleware"
	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/internal/service"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"
	"github.com/rs/zerolog"
)

func NewApp(logger zerolog.Logger, dbpool *pgxpool.Pool, redis *redis.Client, mailService *utils.MailService) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	// Middlewares
	app.Use(middleware.RateLimiter())
	app.Use(middleware.AuthMiddleware())
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// Repositories
	accountRepo := repository.NewAccountRepository(dbpool)
	adminRepo := repository.NewAdminRepository(dbpool)
	feeRepo := repository.NewFeeRepository(dbpool)
	transactionRepo := repository.NewTransactionRepository(dbpool)
	userRepo := repository.NewUserRepository(dbpool)

	// Services
	accountService := service.NewAccountService(accountRepo, userRepo, feeRepo, mailService)
	adminService := service.NewAdminService(adminRepo, redis, mailService)
	feeService := service.NewFeeService(feeRepo)
	transferService := service.NewTransactionService(transactionRepo, accountRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	accountHandler := handler.NewAccountHandler(accountService)
	authHandler := handler.NewAuthHandler(adminService)
	feeHandler := handler.NewFeeHandler(feeService)
	transferHandler := handler.NewTransferHandler(transferService)
	userHandler := handler.NewUserHandler(userService)

	// Register routes and middlewares
	accountHandler.RegisterRoutes(app)
	authHandler.RegisterRoutes(app)
	feeHandler.RegisterRoutes(app)
	transferHandler.RegisterRoutes(app)
	userHandler.RegisterRoutes(app)

	return app
}
