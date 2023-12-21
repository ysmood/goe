// Package goe provide helpers to load environment variables.
package goe

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	envparse "github.com/ysmood/goe/pkg/envparse"
	"golang.org/x/exp/constraints"
)

// Load .env file and return informative message about what this function has done.
// It will recursively search for the `.env` file in parent folders until it finds one.
// It uses [LoadDotEnv] to parse and load the .env content.
func Load(override, expand bool) (string, error) {
	path := LookupFile(".env")
	if path == "" {
		return "No .env file to load", nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to open .env file: %w", err)
	}

	err = LoadDotEnv(override, expand, content)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Loaded environment variables from: %s", path), nil
}

// LoadDotEnv load the .env content.
// If override is true, it will override the existing env vars.
// If expand is true, it will expand the env vars via [os.ExpandEnv].
func LoadDotEnv(override, expand bool, content []byte) error {
	ps, err := envparse.Parse(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse .env file: %w", err)
	}

	for _, p := range ps {
		k, v := p.Key, p.Val

		if !override {
			if _, has := os.LookupEnv(k); has {
				continue
			}
		}

		if expand {
			v = os.ExpandEnv(v)
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
	~bool | ~string | time.Duration | constraints.Float | constraints.Integer
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

// GetList is a shortcut for [GetListWithSep] with separator set to ",".
func GetList[T EnvType](name string, defaultVal []T) []T {
	return GetListWithSep(name, ",", defaultVal)
}

// GetListWithSep returns env var with the name. It will return the defaultVal if it's not found.
// It will parse the value as a list of type T with separator.
func GetListWithSep[T EnvType](name, separator string, defaultVal []T) []T {
	if _, has := os.LookupEnv(name); !has {
		return defaultVal
	}

	var out []T

	for _, s := range strings.Split(Get(name, ""), separator) {
		v, err := Parse[T](s)
		if err != nil {
			panic(err)
		}

		out = append(out, v)
	}

	return out
}

// GetMap is a shortcut for [GetMapWithSep] with pairSep set to "," and kvSep set to ":".
func GetMap[K, V EnvType](name string, defaultVal map[K]V) map[K]V {
	return GetMapWithSep(name, ",", ":", defaultVal)
}

// GetMapWithSep returns env var with the name.
// It will override the key-value pairs in defaultVal with the parsed pairs.
// It will parse the value as a map of type K, V with two types of separators,
// the pairSep is for key-value pairs, and the kvSep is for key and value.
func GetMapWithSep[K, V EnvType](name, pairSep, kvSep string, defaultVal map[K]V) map[K]V {
	str := Get(name, "")

	for _, s := range strings.Split(str, pairSep) {
		kv := strings.Split(s, kvSep)
		if len(kv) != 2 { //nolint: gomnd
			panic("invalid map format: " + str)
		}

		k, err := Parse[K](kv[0])
		if err != nil {
			panic(err)
		}

		v, err := Parse[V](kv[1])
		if err != nil {
			panic(err)
		}

		defaultVal[k] = v
	}

	return defaultVal
}

func GetWithParser[T any](name string, parser func(string) (T, error), defaultVal T) T {
	if _, has := os.LookupEnv(name); has {
		return RequireWithParser(name, parser)
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

// RequireWithParser load and parse the env var with the name.
func RequireWithParser[T any](name string, parser func(string) (T, error)) T {
	v, err := parser(Require[string](name))
	if err != nil {
		panic("failed to parse env variable: " + name + ": " + err.Error())
	}

	return v
}

// Parse the str to the type T.
func Parse[T EnvType](str string) (T, error) { //nolint: cyclop
	v := reflect.ValueOf(new(T)).Elem()
	empty := v.Interface().(T)

	switch v.Interface().(type) { //nolint: gocritic
	case time.Duration:
		d, err := time.ParseDuration(str)
		if err != nil {
			return empty, fmt.Errorf("failed to parse duration: %w", err)
		}

		v.Set(reflect.ValueOf(d))

		return v.Interface().(T), nil
	}

	switch v.Kind() { //nolint: exhaustive
	case reflect.Bool:
		b, err := strconv.ParseBool(str)
		if err != nil {
			return empty, fmt.Errorf("failed to parse bool: %w", err)
		}

		v.Set(reflect.ValueOf(b))

	case reflect.String:
		v = convert(str, v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to parse int: %w", err)
		}

		v = convert(i, v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to parse uint: %w", err)
		}

		v = convert(i, v)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return empty, fmt.Errorf("failed to parse float: %w", err)
		}

		v = convert(f, v)
	}

	return v.Interface().(T), nil
}

// Time parse the str to time.Time.
func Time(str string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return t, nil
}

// ReadFile read file and return the content as string.
// Useful when expanding file path in env var.
func ReadFile(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func convert(from any, to reflect.Value) reflect.Value {
	return reflect.ValueOf(from).Convert(to.Type())
}
