package main

import (
	"context"
	"fmt"
	"os"

	"github.com/makifdb/mini-bank/speedster/internal/app"
	"github.com/makifdb/mini-bank/speedster/internal/config"
	"github.com/makifdb/mini-bank/speedster/internal/redis"
	"github.com/makifdb/mini-bank/speedster/internal/repository"
	"github.com/makifdb/mini-bank/speedster/pkg/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

// @title Mini Bank API
// @version 1.0
// @description This is a sample server for a mini bank.
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
func main() {
	// Load configuration
	var cfg config.Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Set up logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Set up database connection pool
	db, err := repository.NewDatabase(cfg.DBURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Set up redis client
	redisClient, err := redis.NewClient(cfg.RedisAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	// mail service
	mailService := utils.NewMailService(cfg.SmtpServer, cfg.SmtpUser, cfg.SmtpPassword, cfg.SmtpPort)

	// Register routes and middlewares
	app := app.NewApp(logger, db.Pool, redisClient, mailService)

	// Start the server
	logger.Info().Msgf("Starting server on port %d", cfg.Port)
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
