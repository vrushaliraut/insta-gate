package database

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresMigrationsAndPool(t *testing.T) {
	ctx := context.Background()

	// Spin up the Postgres Test Container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second),
		),
	)

	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	// Ensure the container is torn down after the test finishes
	defer pgContainer.Terminate(ctx)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	// Test Connection Pool Wrapper
	db, err := NewPostgres(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to postgres via pgxpool: %v", err)
	}
	defer db.Close()

	// Verify pool is alive and successfully pinged
	if err := db.Pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	// 3. Test Migrations (Up and Down)
	// Resolve the absolute path to the migrations folder
	migrationsPath, err := filepath.Abs("../../migrations")
	if err != nil {
		t.Fatalf("failed to get migrations path: %v", err)
	}

	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		t.Fatalf("failed to initialize golang-migrate: %v", err)
	}
	defer m.Close()

	// Run Migrations Up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations UP: %v", err)
	}

	// Run Migrations Down
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations DOWN: %v", err)
	}
}
