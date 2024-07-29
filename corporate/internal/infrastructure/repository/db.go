// db.go
package repository

import (
	"fmt"

	"github.com/makifdb/mini-bank/corporate/internal/domain/account"
	"github.com/makifdb/mini-bank/corporate/internal/domain/transaction"
	"github.com/makifdb/mini-bank/corporate/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database struct to hold the GORM DB instance
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new GORM database connection
func NewDatabase(dbURL string) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Automatically migrate the schema, while adding indexes
	db.AutoMigrate(
		&user.User{},
		&account.Account{},
		&transaction.Transaction{},
	)

	return &Database{DB: db}, nil
}

// Close to close the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
