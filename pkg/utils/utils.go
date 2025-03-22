package utils

import (
	"fmt"
	"os"
)

func Println(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
}
