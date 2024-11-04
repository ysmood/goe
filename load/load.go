// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"fmt"
	"os"

	"github.com/ysmood/goe"
)

func init() {
	err := load()
	if err != nil {
		fmt.Fprintln(os.Stderr, goe.Prefix, err)
		os.Exit(1)
	}
}

func load() error {
	return goe.AutoLoad(goe.DOTENV)
}
