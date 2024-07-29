package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         int    `mapstructure:"PORT"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	DBURL        string `mapstructure:"DB_URL"`
	RedisAddr    string `mapstructure:"REDIS_ADDR"`
	SmtpServer   string `mapstructure:"SMTP_SERVER"`
	SmtpPort     int    `mapstructure:"SMTP_PORT"`
	SmtpUser     string `mapstructure:"SMTP_USER"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
	Environment  string `mapstructure:"ENVIRONMENT"`
}

func NewConfig() *Config {
	return &Config{
		Port:         getEnvInt("PORT", 3000),
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		DBURL:        getEnv("DB_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		SmtpServer:   getEnv("SMTP_SERVER", "smtp.resend.com"),
		SmtpPort:     getEnvInt("SMTP_PORT", 587),
		SmtpUser:     getEnv("SMTP_USER", "resend"),
		SmtpPassword: getEnv("SMTP_PASSWORD", "password"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv reads an environment variable or returns a default value if not found
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt reads an environment variable as an integer or returns a default value if not found
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
