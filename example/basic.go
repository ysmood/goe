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

	// Get a list from a env variable
	list := goe.GetList("LIST", []int{})

	// Get a map from a env variable
	numMap := goe.GetMap("MAP", map[int]int{})

	// Get a expanded env variable from expression like EXPANDED=${ENV}
	expanded := goe.Get("EXPANDED", "")

	// It will output the env variables in "../.env"
	// Output:
	//	1 hello true [1 2] map[1:2 3:4] dev
	fmt.Println(num, str, isDev, list, numMap, expanded)
}
