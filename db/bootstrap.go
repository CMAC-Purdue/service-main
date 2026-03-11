package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schemaSQL string

// EnsureSchema applies the SQL schema used by this service.
func EnsureSchema(ctx context.Context, sqlDB *sql.DB) error {
	if _, err := sqlDB.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("unable to apply schema: %w", err)
	}
	return nil
}
