// Package goe provide helpers to load environment variables.
package goe

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	envparse "github.com/hashicorp/go-envparse"
	"golang.org/x/exp/constraints"
)

// Load .env file and return informative message about what this function has done.
// It will recursively search for the `.env` file in parent folders until it finds one.
func Load(override bool) (string, error) {
	path := LookupFile(".env")
	if path == "" {
		return "No .env file to load", nil
	}

	content, err := os.ReadFile(path) //nolint: gosec
	if err != nil {
		return "", fmt.Errorf("failed to open .env file: %w", err)
	}

	err = LoadDotEnv(override, content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Loaded environment variables from: %s", path), nil
}

// LoadDotEnv load the .env content.
func LoadDotEnv(override bool, content []byte) error {
	dict, err := envparse.Parse(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse .env file: %w", err)
	}

	for k, v := range dict {
		if !override {
			if _, has := os.LookupEnv(k); has {
				continue
			}
		}

		err = os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("failed to set env variable: %w", err)
		}
	}

	return nil
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
	bool | string | time.Duration | constraints.Float | constraints.Integer
}

// Is check if the env var with the name is equal to the val.
// If the env var is not found, it will return false.
func Is[T EnvType](name string, val T) bool {
	if _, has := os.LookupEnv(name); !has {
		return false
	}

	return Require[T](name) == val
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

	v, err := Parse[T](envStr)
	if err != nil {
		panic(err)
	}

	return v
}

// Parse the str to the type T.
func Parse[T EnvType](str string) (T, error) {
	var v any = *new(T)

	switch v.(type) {
	case bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return v.(T), fmt.Errorf("failed to parse bool: %w", err)
		}

		v = b

	case string:
		v = str

	case int, int8, int16, int32, int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return v.(T), fmt.Errorf("failed to parse int: %w", err)
		}

		v = convert(i, v)

	case uint, uint8, uint16, uint32, uint64:
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return v.(T), fmt.Errorf("failed to parse uint: %w", err)
		}

		v = convert(i, v)

	case float32, float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return v.(T), fmt.Errorf("failed to parse float: %w", err)
		}

		v = convert(f, v)

	case time.Duration:
		d, err := time.ParseDuration(str)
		if err != nil {
			return v.(T), fmt.Errorf("failed to parse duration: %w", err)
		}

		v = d
	}

	return v.(T), nil
}

// ReadFile read file and return the content as string.
func ReadFile(path string) string {
	b, err := os.ReadFile(path) //nolint: gosec
	if err != nil {
		panic(err)
	}

	return string(b)
}

func convert(from any, to any) any {
	return reflect.ValueOf(from).Convert(reflect.TypeOf(to)).Interface()
}
