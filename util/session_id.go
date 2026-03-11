package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const defaultSessionIDBytes = 32

// GenerateSessionID creates a URL-safe random session ID.
func GenerateSessionID() (string, error) {
	return GenerateSessionIDSize(defaultSessionIDBytes)
}

// GenerateSessionIDSize creates a URL-safe random session ID from size random bytes.
func GenerateSessionIDSize(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("size must be greater than 0")
	}

	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("unable to generate random session id: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
