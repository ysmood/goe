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

	list := goe.GetList("LIST", []int{})

	numMap := goe.GetMap("MAP", map[int]int{})

	// It will output the env variables in "../.env"
	// Output:
	//	1 hello true [1 2] map[1:2 3:4]
	fmt.Println(num, str, isDev, list, numMap)
}
