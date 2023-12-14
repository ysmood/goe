// Package dotenv loads environment variables from a .env file.
package dotenv

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ysmood/dotenv/lib"
)

type info struct{}

var prefix = fmt.Sprintf("[%s] ", reflect.TypeOf(info{}).PkgPath())

func init() {
	lg := log.New(log.Writer(), prefix, log.Flags())

	msg, err := lib.Load()
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println(msg)
}

// Get is a shortcut for [lib.Get].
func Get[T lib.EnvType](name string, defaultVal T) T {
	return lib.Get(name, defaultVal)
}

// Require is a shortcut for [lib.Require].
func Require[T lib.EnvType](name string) T {
	return lib.Require[T](name)
}
