package app

import (
	"database/sql"
	"net/http"

	"github.com/makifdb/mini-bank/minimalist/internal/handler"
	"github.com/redis/go-redis/v9"

	"github.com/makifdb/mini-bank/minimalist/internal/repository"
	"github.com/makifdb/mini-bank/minimalist/internal/service"
	"github.com/makifdb/mini-bank/minimalist/pkg/utils"
)

func NewApp(db *sql.DB, rds *redis.Client, mailService *utils.MailService) *http.ServeMux {
	router := http.NewServeMux()

	// Repositories
	accountRepo := repository.NewAccountRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	feeRepo := repository.NewFeeRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Services
	accountService := service.NewAccountService(accountRepo, userRepo, feeRepo, mailService)
	adminService := service.NewAdminService(adminRepo, rds, mailService)
	feeService := service.NewFeeService(feeRepo)
	transferService := service.NewTransactionService(transactionRepo, accountRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	accountHandler := handler.NewAccountHandler(accountService)
	adminHandler := handler.NewAdminHandler(adminService)
	feeHandler := handler.NewFeeHandler(feeService)
	transferHandler := handler.NewTransferHandler(transferService)
	userHandler := handler.NewUserHandler(userService)

	// Register routes and middlewares
	accountHandler.RegisterRoutes(router)
	adminHandler.RegisterRoutes(router)
	feeHandler.RegisterRoutes(router)
	transferHandler.RegisterRoutes(router)
	userHandler.RegisterRoutes(router)

	return router
}
