// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/ysmood/goe"
)

type info struct{}

var prefix = fmt.Sprintf("[%s]", reflect.TypeOf(info{}).PkgPath())

func init() {
	err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, prefix, err)
		os.Exit(1)
	}
}

func load() error {
	err := goe.Load(false, true, goe.DOTENV)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(prefix + " .env file not found, skipped loading.")

			return nil
		}

		return err //nolint:wrapcheck
	}

	path, _ := goe.LookupFile(goe.DOTENV)

	fmt.Println(prefix+" Loaded environment variables from:", path)

	return nil
}
