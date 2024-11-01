package utils

import "fmt"

func Println(a ...interface{}) {
	_, _ = fmt.Println(a...) //nolint: forbidigo
}
