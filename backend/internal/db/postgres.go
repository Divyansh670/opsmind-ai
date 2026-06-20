package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Config struct {
	DSN             string
	MaxOpenConns    int32
	ConnMaxLifetime time.Duration
}

func DefaultConfig(dsn string) Config {
	return Config{
		DSN:             dsn,
		MaxOpenConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

func New(cfg Config) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✓ Database connection pool established")
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	fmt.Println("Closing database connection pool...")
	db.Pool.Close()
}

func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}
