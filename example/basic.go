package main

import (
	"fmt"
	"os"

	_ "github.com/ysmood/dotenv"
)

func main() {
	fmt.Println(os.Getenv("VALUE"))
}
