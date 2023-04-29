// Package lib ...
package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Load .env file return informative message.
func Load() (string, error) {
	path := LookupFile(".env")
	if path == "" {
		return "No .env file to load", nil
	}

	err := godotenv.Load(path)
	if err != nil {
		return "", fmt.Errorf("godotenv: %w", err)
	}

	return fmt.Sprintf("Loaded environment variables from: %s", path), nil
}

// LookupFile file recursively.
func LookupFile(file string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	prev := ""

	for dir != prev {
		p := filepath.Join(dir, file)
		if _, err := os.Stat(p); err == nil {
			return p
		}

		prev = dir
		dir = filepath.Dir(dir)
	}

	return ""
}
