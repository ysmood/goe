package main

import (
	"fmt"

	dotenv "github.com/ysmood/dotenv"
)

func main() {
	num := dotenv.Get("NUM", 0)
	str := dotenv.Get("STR", "")

	// It will output the env variables in "../.env"
	// The output will be:
	//	2 hello world
	fmt.Println(num+1, str+" world")
}
