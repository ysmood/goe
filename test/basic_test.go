package test_test

import (
	"fmt"
	"os"

	_ "github.com/ysmood/goe/load"
)

func Example_loadDotEnvFileRecursively() {
	fmt.Println(os.Getenv("STR"))

	// Output: hello
}
