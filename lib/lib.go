// Package lib ...
package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

type EnvType interface {
	bool | string | int | float64 | time.Duration
}

// Get env var with default value.
func Get[T EnvType](name string, defaultVal T) T {
	envStr, has := os.LookupEnv(name)
	if !has {
		return defaultVal
	}

	var v any = defaultVal

	switch v.(type) {
	case bool:
		b, err := strconv.ParseBool(envStr)
		if err != nil {
			panic(err)
		}

		v = b

	case string:
		v = envStr

	case int:
		i, err := strconv.ParseInt(envStr, 10, 64)
		if err != nil {
			panic(err)
		}

		v = int(i)

	case float64:
		f, err := strconv.ParseFloat(envStr, 64)
		if err != nil {
			panic(err)
		}

		v = f

	case time.Duration:
		d, err := time.ParseDuration(envStr)
		if err != nil {
			panic(err)
		}

		v = d
	}

	return v.(T)
}
