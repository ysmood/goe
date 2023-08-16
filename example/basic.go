package main

import (
	"fmt"
	"os"

	_ "github.com/ysmood/dotenv"
)

func main() {
	// It will output the VALUE in "../.env" which is "ok"
	fmt.Println(os.Getenv("VALUE"))
}
