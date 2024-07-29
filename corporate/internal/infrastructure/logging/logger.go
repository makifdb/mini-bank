package logging

import (
	"github.com/makifdb/mini-bank/corporate/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger based on the environment.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var config zap.Config
	if cfg.Environment == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

// Add more log levels
var (
	Debug = zap.DebugLevel
	Info  = zap.InfoLevel
	Warn  = zap.WarnLevel
	Error = zap.ErrorLevel
	Fatal = zap.FatalLevel
)

// ConfigureLogger sets the log level based on environment
func ConfigureLogger(logger *zap.Logger, level zapcore.Level) *zap.Logger {
	return logger.WithOptions(zap.IncreaseLevel(level))
}
