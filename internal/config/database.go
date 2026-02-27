package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectDB initializes and returns a connection pool to the PostgreSQL database.
func ConnectDB(dbURL string) (*pgxpool.Pool, error) {
	// 1. Parse the connection string from .env into a configuration object
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// 2. Setup pool configuration (Best practice for production)
	poolConfig.MaxConns = 10                      // Maximum number of active connections
	poolConfig.MinConns = 2                       // Minimum idle connections to keep alive
	poolConfig.MaxConnIdleTime = 30 * time.Minute // Close idle connections after 30 minutes

	// 3. Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 4. Ping the database to ensure the connection is actually established and valid
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
