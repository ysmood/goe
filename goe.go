// Package goe provide helpers to load environment variables.
package goe

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ysmood/goe/pkg/dotenv"
	"golang.org/x/exp/constraints"
)

// Load .env file and return informative message about what this function has done.
// It will recursively search for the file in parent folders until it finds one.
// It uses [LoadDotEnv] to parse and load the content.
func Load(override, expand bool, file string) error {
	path, err := LookupFile(file)
	if err != nil {
		return fmt.Errorf("failed to find .env file: %w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}

	err = LoadContent(override, expand, string(content))
	if err != nil {
		return fmt.Errorf("failed to load .env content: %w", err)
	}

	return nil
}

// LoadContent load the .env content.
// If override is true, it will override the existing env vars.
// If expand is true, it will expand the env vars via [os.ExpandEnv].
func LoadContent(override, expand bool, content string) error {
	doc, err := dotenv.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse .env file: %w", err)
	}

	for _, p := range doc.Lines {
		k, v := p.Key, p.Value

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
func LookupFile(file string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	prev := ""

	for dir != prev {
		p := filepath.Join(dir, file)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}

		prev = dir
		dir = filepath.Dir(dir)
	}

	return "", fmt.Errorf("%w: %s", os.ErrNotExist, file)
}

type EnvType interface {
	EnvKeyType | ~[]byte
}

type EnvKeyType interface {
	~bool | ~string | time.Duration | constraints.Float | constraints.Integer
}

// Is checks if the env var with the name is equal to the val.
func Is(name string, val string) bool {
	return Get(name, "") == val
}

// Has checks if the env var with the name exists.
func Has(name string) bool {
	_, has := os.LookupEnv(name)

	return has
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
	l, err := GetListWithSep(name, ",", defaultVal)
	if err != nil {
		panic("failed to parse list: " + err.Error())
	}

	return l
}

// GetListWithSep returns env var with the name. It will return the defaultVal if it's not found.
// It will parse the value as a list of type T with separator.
func GetListWithSep[T EnvType](name, separator string, defaultVal []T) ([]T, error) {
	if _, has := os.LookupEnv(name); !has {
		return defaultVal, nil
	}

	var out []T //nolint: prealloc

	for _, s := range strings.Split(Get(name, ""), separator) {
		v, err := Parse[T](s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse list: %w", err)
		}

		out = append(out, v)
	}

	return out, nil
}

// GetMap is a shortcut for [GetMapWithSep] with pairSep set to "," and kvSep set to ":".
func GetMap[K EnvKeyType, V EnvType](name string, defaultVal map[K]V) map[K]V {
	m, err := GetMapWithSep(name, ",", ":", defaultVal)
	if err != nil {
		panic("failed to parse map: " + err.Error())
	}

	return m
}

var ErrInvalidMapFormat = errors.New("invalid map format")

// GetMapWithSep returns env var with the name.
// It will override the key-value pairs in defaultVal with the parsed pairs.
// It will parse the value as a map of type K, V with two types of separators,
// the pairSep is for key-value pairs, and the kvSep is for key and value.
func GetMapWithSep[K EnvKeyType, V EnvType](name, pairSep, kvSep string, defaultVal map[K]V) (map[K]V, error) {
	str := Get(name, "")

	for _, s := range strings.Split(str, pairSep) {
		kv := strings.Split(s, kvSep)
		if len(kv) != 2 { //nolint: mnd
			return nil, fmt.Errorf("%w: %s", ErrInvalidMapFormat, str)
		}

		k, err := Parse[K](kv[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse map key: %w", err)
		}

		v, err := Parse[V](kv[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse map value: %w", err)
		}

		defaultVal[k] = v
	}

	return defaultVal, nil
}

func GetWithParser[T any](name string, parser func(string) (T, error), defaultVal T) T {
	if _, has := os.LookupEnv(name); has {
		return RequireWithParser(name, parser)
	}

	return defaultVal
}

// Require load and parse the env var with the name.
// It will panic if the env var is not found or failed to parse.
// It uses [Parse] to parse the env var.
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

var ErrUnsupportedSliceType = errors.New("unsupported slice type")

// Parse the str to the type T.
// It will auto detect the type of the env var and parse it.
// If T is []byte and str is a existing file path, the file content will be the env var,
// or the str will be parsed as base64 and used as the env var.
func Parse[T EnvType](str string) (T, error) { //nolint: funlen,cyclop
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

	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return empty, fmt.Errorf("%w: %s", ErrUnsupportedSliceType, v.Type().String())
		}

		b, err := os.ReadFile(str)
		if err == nil {
			return convert(b, v).Interface().(T), nil
		}

		b, err = base64.StdEncoding.DecodeString(str)
		if err != nil {
			return empty, fmt.Errorf("failed to parse base64: %w", err)
		}

		v = convert(b, v)
	}

	return v.Interface().(T), nil
}

// Unset env var with [os.Unsetenv].
// Useful for secret unset to prevent leaking by other packages.
func Unset(name string) struct{} {
	err := os.Unsetenv(name)
	if err != nil {
		panic("failed to unset env variable: " + name)
	}

	return struct{}{}
}

// Time parse the str to time.Time.
func Time(str string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return t, nil
}

func convert(from any, to reflect.Value) reflect.Value {
	return reflect.ValueOf(from).Convert(to.Type())
}
