package main

import (
	"fmt"

	dotenv "github.com/ysmood/dotenv"
)

func main() {
	// Get a optional env variable
	num := dotenv.Get("NUM", 0)

	// Get a required env variable
	str := dotenv.Require[string]("STR")

	// It will output the env variables in "../.env"
	// The output will be:
	//	2 hello world
	fmt.Println(num+1, str+" world")
}
