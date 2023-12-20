// Package goe provide helpers to load environment variables.
package goe

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	envparse "github.com/hashicorp/go-envparse"
)

// Load .env file and return informative message about what this function has done.
// It will recursively search for the `.env` file in parent folders until it finds one.
func Load(override bool) (string, error) {
	path := LookupFile(".env")
	if path == "" {
		return "No .env file to load", nil
	}

	file, err := os.Open(path) //nolint: gosec
	if err != nil {
		return "", fmt.Errorf("failed to open .env file: %w", err)
	}

	dict, err := envparse.Parse(file)
	if err != nil {
		return "", fmt.Errorf("failed to parse .env file: %w", err)
	}

	for k, v := range dict {
		if !override {
			if _, has := os.LookupEnv(k); has {
				continue
			}
		}

		err = os.Setenv(k, v)
		if err != nil {
			return "", fmt.Errorf("failed to set env variable: %w", err)
		}
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

// Get env var with the name. It will return the defaultVal if it's not found.
// If the env var is found, it will use [Require] to parse the value.
func Get[T EnvType](name string, defaultVal T) T {
	if _, has := os.LookupEnv(name); has {
		return Require[T](name)
	}

	return defaultVal
}

// Require load and parse the env var with the name.
// It will auto detect the type of the env var and parse it.
// It will panic if the env var is not found.
func Require[T EnvType](name string) T {
	envStr, has := os.LookupEnv(name)
	if !has {
		panic("required env variable not found: " + name)
	}

	var v any = *new(T)

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