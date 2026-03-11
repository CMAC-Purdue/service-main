package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

func NewPoolWithRetry(ctx context.Context, connString string, timeout time.Duration, retryInterval time.Duration) (*pgxpool.Pool, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}
	if retryInterval <= 0 {
		return nil, fmt.Errorf("retry interval must be greater than 0")
	}

	timedCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var lastErr error
	for {
		pool, err := NewPool(timedCtx, connString)
		if err == nil {
			return pool, nil
		}
		lastErr = err

		timer := time.NewTimer(retryInterval)
		select {
		case <-timedCtx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, fmt.Errorf("unable to connect to database within %s: %w", timeout, lastErr)
		case <-timer.C:
		}
	}
}
