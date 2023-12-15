// Package load loads environment variables from a .env file,
// it will recursively search for the file in parent folders until it finds one.
// It will print all the side effects to the console to make it easier to debug.
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

	msg, err := goe.Load(false)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println(msg)
}
