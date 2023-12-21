package main

import (
	"fmt"

	"github.com/ysmood/goe"
	_ "github.com/ysmood/goe/load" // load the .env file
)

var (
	// Get a optional env variable.
	num = goe.Get("NUM", 0)

	// Get a required env variable.
	secret = goe.Require[string]("SECRET")
	_      = goe.Unset("SECRET")

	// Check if the env variable is equal to specified value.
	isDev = goe.Is("ENV", "dev")

	// Get a list from a env variable.
	list = goe.GetList("LIST", []int{})

	// Get a map from a env variable.
	numMap = goe.GetMap("MAP", map[int]int{})

	// Get a expanded env variable from expression like EXPANDED=${ENV}.
	expanded = goe.Get("EXPANDED", "")

	// Get a custom type from env variable.
	time = goe.RequireWithParser("TIME", goe.Time)
)

func main() {
	// It will output the env variables in "../.env"
	// Output:
	//	1 hello true [1 2] map[1:2 3:4] dev 2023
	fmt.Println(num, secret, isDev, list, numMap, expanded, time.Format("2006"))
}
