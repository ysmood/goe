package goe

import (
	"fmt"
	"reflect"
)

const (
	DOTENV = ".env"
)

type info struct{}

var Prefix = fmt.Sprintf("[%s]", reflect.TypeOf(info{}).PkgPath())
