package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres holds the concurrency-safe connection pool for pgx.
type Postgres struct {
	Pool *pgxpool.Pool
}

// NewPostgres initializes and returns a wrapped PostgreSQL connection pool.
func NewPostgres(ctx context.Context, connString string) (*Postgres, error) {
	// Parse the connection string into a configuration struct.
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database url: %w", err)
	}

	// Tune pool settings for enterprise scale
	config.MaxConns = 15
	config.MinConns = 2
	config.MaxConnIdleTime = 15 * time.Minute
	config.MaxConnLifetime = 1 * time.Hour

	// Enforce a strict timeout for the initial connection attempt
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Initialize the pool with the mapped configuration
	pool, err := pgxpool.NewWithConfig(connectCtx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify the connection is actively accepting requests
	if err := pool.Ping(connectCtx); err != nil {
		return nil, fmt.Errorf("unable to ping the database: %w", err)
	}

	return &Postgres{Pool: pool}, nil
}

// Close safely drains the connection pool. Blocks until all connections are returned to the pool and closed.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
