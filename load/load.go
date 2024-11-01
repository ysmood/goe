// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"errors"
	"fmt"
	"os"

	"github.com/ysmood/goe"
	"github.com/ysmood/goe/pkg/utils"
)

func init() {
	err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, goe.Prefix, err)
		os.Exit(1)
	}
}

func load() error {
	err := goe.Load(false, true, goe.DOTENV)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			utils.Println(goe.Prefix, ".env file not found, skipped loading.")

			return nil
		}

		return err //nolint:wrapcheck
	}

	path, _ := goe.LookupFile(goe.DOTENV)

	utils.Println(goe.Prefix, "Loaded environment variables from:", path)

	return nil
}
