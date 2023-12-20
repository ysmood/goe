package main

import (
	"fmt"

	"github.com/ysmood/goe"
	_ "github.com/ysmood/goe/load" // load the .env file
)

func main() {
	// Get a optional env variable
	num := goe.Get("NUM", 0)

	// Get a required env variable
	str := goe.Require[string]("STR")

	// Check if the env variable is equal to specified value
	isDev := goe.Is("ENV", "dev")

	// It will output the env variables in "../.env"
	// The output will be:
	//	2 hello world true
	fmt.Println(num+1, str+" world", isDev)
}
