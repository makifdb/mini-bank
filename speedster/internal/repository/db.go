package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(dbURL string) (*Database, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL: %v", err)
	}

	// Optimize connection pool settings
	config.MaxConns = 10                     // Set the maximum number of connections
	config.MinConns = 2                      // Set the minimum number of connections
	config.MaxConnIdleTime = time.Minute * 5 // Set the maximum idle time for connections
	config.MaxConnLifetime = time.Hour * 1   // Set the maximum lifetime for a connection

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
