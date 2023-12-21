package goe_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/ysmood/goe"
	"github.com/ysmood/got"
)

func TestExample(t *testing.T) {
	g := got.T(t)

	out, err := exec.Command("go", "run", "./example").CombinedOutput()
	g.E(err)

	g.Has(string(out), `goe/.env
1 hello true [1 2] map[1:2 3:4] dev`)
}

func TestGet(t *testing.T) {
	g := got.T(t)

	t.Setenv("ENV", "dev")
	t.Setenv("BOOL", "true")
	t.Setenv("NUM", "2")
	t.Setenv("STR", "ok")
	t.Setenv("FLOAT", "1.2")
	t.Setenv("DURATION", "1m")

	g.True(goe.Is("ENV", "dev"))
	g.False(goe.Is("ENV", "stg"))
	g.False(goe.Is("NOT_EXISTS", "dev"))

	g.Eq(goe.Get("BOOL", false), true)
	g.Eq(goe.Get("BOOL_DEFAULT", true), true)
	g.Eq(goe.Get("NUM", 0), 2)
	g.Eq(goe.Get("NUM_DEFAULT", uint(1)), 1)
	g.Eq(goe.Get("STR", ""), "ok")
	g.Eq(goe.Get("STR_DEFAULT", "yes"), "yes")
	g.Eq(goe.Get("FLOAT", 0.0), 1.2)
	g.Eq(goe.Get("FLOAT_DEFAULT", 1.1), 1.1)
	g.Eq(goe.Get("DURATION", time.Second), time.Minute)
	g.Eq(goe.Get("DURATION_DEFAULT", time.Hour), time.Hour)

	g.Eq(goe.Require[int]("NUM"), 2)

	g.Eq(g.Panic(func() {
		goe.Require[int]("KEY")
	}), "required env variable not found: KEY")

	g.Eq(g.Panic(func() {
		t.Setenv("WRONG_INT", "xxx")
		goe.Get("WRONG_INT", 0)
	}).(error).Error(), `failed to parse int: strconv.ParseInt: parsing "xxx": invalid syntax`)

	g.Has(goe.ReadFile("go.mod"), "github.com/ysmood/goe")
}

func TestLoad(t *testing.T) {
	g := got.T(t)

	g.E(goe.Load(false, true))

	g.Eq(goe.Get("STR", ""), "hello")
	g.Has(goe.Get("EXPANDED", ""), "dev")
}

func TestGetList(t *testing.T) {
	g := got.T(t)

	t.Setenv("LIST", "1,2")
	g.Eq(goe.GetList("LIST", []int{}), []int{1, 2})

	g.Eq(goe.GetList("DEFAULT", []int{1}), []int{1})
}

func TestGetMap(t *testing.T) {
	g := got.T(t)

	t.Setenv("MAP", "a:1,b:2")
	g.Eq(
		goe.GetMap("MAP", map[string]int{"a": 2, "c": 3}),
		map[string]int{"a": 1, "b": 2, "c": 3},
	)
}
