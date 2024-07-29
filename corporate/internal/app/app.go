package app

import (
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/makifdb/mini-bank/corporate/docs"
	"github.com/makifdb/mini-bank/corporate/internal/application/middleware"
	"github.com/makifdb/mini-bank/corporate/internal/application/service"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/email"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/logging"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
	"github.com/makifdb/mini-bank/corporate/internal/interface/handler"
	"github.com/makifdb/mini-bank/corporate/pkg/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RateLimiterMiddleware(60))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

var Module = fx.Options(
	fx.Provide(
		config.NewConfig,
		logging.NewLogger,
		newRouter,
		repository.NewDatabase,
		repository.NewUserRepository,
		repository.NewAdminRepository,
		repository.NewAccountRepository,
		repository.NewTransactionRepository,
		repository.NewFeeRepository,
		repository.NewGormTransactionManager,
		email.NewMailService,
		service.NewAccountService,
		service.NewAdminService,
		service.NewUserService,
		service.NewFeeService,
		service.NewTransferService,
		handler.NewAccountHandler,
		handler.NewAdminHandler,
		handler.NewUserHandler,
		handler.NewFeeHandler,
		handler.NewTransferHandler,
	),
	fx.Invoke(registerHooks),
)

func registerHooks(lifecycle fx.Lifecycle, logger *zap.Logger) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Info("Application starting")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Application stopping")
				return nil
			},
		},
	)
}
