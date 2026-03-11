package util

import (
	"database/sql"
	"strings"
)

func ToNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return sql.NullString{}
	}

	return sql.NullString{String: trimmed, Valid: true}
}

func NullStringToPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	result := value.String
	return &result
}
