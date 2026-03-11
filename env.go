package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// loadDotEnv loads key=value pairs from a local .env file.
// Existing environment variables are not overridden.
func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fmt.Errorf("invalid .env line %d: missing '='", lineNo)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("invalid .env line %d: empty key", lineNo)
		}

		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		value = os.Expand(value, func(name string) string {
			if existing, exists := os.LookupEnv(name); exists {
				return existing
			}
			return ""
		})

		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("unable to set %s from .env: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("unable to read %s: %w", path, err)
	}

	return nil
}
