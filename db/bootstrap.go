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

func InitDB() *Queries {
	ctx := context.Background()

	if err := util.LoadDotEnv(".env"); err != nil {
		log.Fatal(err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	dbConnectTimeout, err := util.EnvDuration("DB_CONNECT_TIMEOUT", defaultDBConnectTimeout)
	if err != nil {
		log.Fatal(err)
	}
	dbRetryInterval, err := util.EnvDuration("DB_CONNECT_RETRY_INTERVAL", defaultDBRetryInterval)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Database schema ensured")

	log.Printf("connecting to database (timeout=%s, retry_interval=%s)", dbConnectTimeout, dbRetryInterval)
	pool, err := NewPoolWithRetry(ctx, connString, dbConnectTimeout, dbRetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	if err := EnsureSchema(ctx, sqlDB); err != nil {
		log.Fatal(err)
	}

	queries := New(sqlDB)
	return queries
}
