// Package load loads environment variables from a .env file.
// It uses [goe.Load] to load the .env file, override set to false, expand set to true.
package load

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ysmood/goe"
)

type info struct{}

var prefix = fmt.Sprintf("[%s] ", reflect.TypeOf(info{}).PkgPath())

func init() {
	lg := log.New(log.Writer(), prefix, log.Flags())

	msg, err := goe.Load(false, true)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println(msg)
}
