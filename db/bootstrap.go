package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"service-main/util"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed schema.sql
var schemaSQL string

const (
	defaultDBConnectTimeout = 45 * time.Second
	defaultDBRetryInterval  = 2 * time.Second
)

// EnsureSchema applies the SQL schema used by this service.
func EnsureSchema(ctx context.Context, sqlDB *sql.DB) error {
	if _, err := sqlDB.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("unable to apply schema: %w", err)
	}
	return nil
}

func InitDB() (*Queries, func(), error) {
	ctx := context.Background()

	if err := util.LoadDotEnv(".env"); err != nil {
		return nil, nil, fmt.Errorf("unable to load env: %w", err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	dbConnectTimeout, err := util.EnvDuration("DB_CONNECT_TIMEOUT", defaultDBConnectTimeout)
	if err != nil {
		return nil, nil, err
	}
	dbRetryInterval, err := util.EnvDuration("DB_CONNECT_RETRY_INTERVAL", defaultDBRetryInterval)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("connecting to database (timeout=%s, retry_interval=%s)", dbConnectTimeout, dbRetryInterval)
	pool, err := NewPoolWithRetry(ctx, connString, dbConnectTimeout, dbRetryInterval)
	if err != nil {
		return nil, nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	if err := EnsureSchema(ctx, sqlDB); err != nil {
		_ = sqlDB.Close()
		pool.Close()
		return nil, nil, err
	}
	log.Print("database schema ensured")

	queries := New(sqlDB)
	cleanup := func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("error closing sql db: %v", err)
		}
		pool.Close()
	}

	return queries, cleanup, nil
}
