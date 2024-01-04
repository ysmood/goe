// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ysmood/goe"
)

type info struct{}

var prefix = fmt.Sprintf("[%s]", reflect.TypeOf(info{}).PkgPath())

func init() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, prefix, r)
			os.Exit(1)
		}
	}()

	err := goe.Load(false, true, goe.DOTENV)
	if err != nil {
		panic(err)
	}

	path, err := goe.LookupFile(goe.DOTENV)
	if err != nil {
		panic(err)
	}

	fmt.Println(prefix+" Loaded environment variables from:", path)
}
