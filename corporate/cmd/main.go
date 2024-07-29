package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/makifdb/mini-bank/corporate/internal/app"
	"github.com/makifdb/mini-bank/corporate/internal/application/middleware"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/logging"
	"github.com/makifdb/mini-bank/corporate/internal/infrastructure/repository"
	"github.com/makifdb/mini-bank/corporate/internal/interface/handler"
	"github.com/makifdb/mini-bank/corporate/pkg/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title Mini Bank API
// @version 1.0
// @description This is a sample server for a mini bank application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

func registerRoutes(
	router *gin.Engine,
	adminHandler *handler.AdminHandler,
	accountHandler *handler.AccountHandler,
	userHandler *handler.UserHandler,
	feeHandler *handler.FeeHandler,
	transferHandler *handler.TransferHandler,
) {
	authGroup := router.Group("/v1")
	adminHandler.RegisterRoutes(authGroup)

	api := router.Group("/v1")
	api.Use(middleware.AuthMiddleware())
	{
		accountHandler.RegisterRoutes(api)
		userHandler.RegisterRoutes(api)
		feeHandler.RegisterRoutes(api)
		handler.RegisterTransferRoutes(api, transferHandler)
	}
}

func main() {
	app := fx.New(
		app.Module,
		fx.Invoke(registerRoutes),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(lc fx.Lifecycle, cfg *config.Config, r *gin.Engine, db *repository.Database, logger *zap.Logger) {
	logger = logging.ConfigureLogger(logger, logging.Info)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting HTTP server", zap.Int("port", cfg.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("ListenAndServe failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			db.Close()
			return srv.Shutdown(shutdownCtx)
		},
	})
}
