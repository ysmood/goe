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

	cmd := exec.Command("go", "run", "./example")
	out, err := cmd.CombinedOutput()
	g.Desc("%s", string(out)).E(err)

	g.Has(string(out), `.env
1 hello true [1 2] map[1:2 3:4] dev hello true 2023`)
}

func TestGet(t *testing.T) {
	g := got.T(t)

	t.Setenv("ENV", "dev")
	t.Setenv("BOOL", "true")
	t.Setenv("NUM", "2")
	t.Setenv("STR", "ok")
	t.Setenv("FLOAT", "1.2")
	t.Setenv("DURATION", "1m")
	t.Setenv("TIME", "2023-12-21T15:41:51+08:00")

	type MyInt int

	g.Eq(goe.Get("NUM", MyInt(0)), 2)

	g.Eq(goe.Get("BOOL", false), true)
	g.Eq(goe.Get("BOOL_DEFAULT", true), true)
	g.Eq(goe.Get("NUM", 0), 2)
	g.Eq(goe.Get("NUM", uint(0)), 2)
	g.Eq(goe.Get("NUM_DEFAULT", uint(1)), 1)
	g.Eq(goe.Get("STR", ""), "ok")
	g.Eq(goe.Get("STR_DEFAULT", "yes"), "yes")
	g.Eq(goe.Get("FLOAT", 0.0), 1.2)
	g.Eq(goe.Get("FLOAT_DEFAULT", 1.1), 1.1)
	g.Eq(goe.Get("DURATION", time.Second), time.Minute)
	g.Eq(goe.Get("DURATION_DEFAULT", time.Hour), time.Hour)

	g.Eq(goe.GetWithParser("TIME", goe.Time, time.Time{}).Format("2006"), "2023")

	tm, err := time.Parse("2006", "2023")
	g.E(err)
	g.Eq(goe.GetWithParser("TIME_DEFAULT", goe.Time, tm).Format("2006"), "2023")

	g.Eq(goe.Require[int]("NUM"), 2)

	g.Eq(g.Panic(func() {
		goe.Require[int]("KEY")
	}), "required env variable not found: KEY")

	g.Eq(g.Panic(func() {
		t.Setenv("WRONG_INT", "xxx")
		goe.Get("WRONG_INT", 0)
	}).(error).Error(), `failed to parse int: strconv.ParseInt: parsing "xxx": invalid syntax`)
}

func TestIs(t *testing.T) {
	g := got.T(t)

	t.Setenv("ENV", "dev")

	g.True(goe.Is("ENV", "dev"))
	g.False(goe.Is("ENV", "stg"))
	g.False(goe.Is("NOT_EXISTS", "dev"))
	g.True(goe.Is("NOT_EXISTS", ""))
}

func TestHas(t *testing.T) {
	g := got.T(t)

	t.Setenv("ENV", "dev")

	g.True(goe.Has("ENV"))
	g.False(goe.Has("NOT_EXISTS"))
}

func TestLoad(t *testing.T) {
	g := got.T(t)

	g.Chdir("example")

	g.E(goe.Load(false, true, ".env"))

	g.Eq(goe.Get("SECRET", ""), "hello")
	g.Has(goe.Get("EXPANDED", ""), "dev")
	g.Eq(goe.Get("BIN", []byte{}), []byte("hello"))
}

func TestGetList(t *testing.T) {
	g := got.T(t)

	t.Setenv("LIST", "1,2")
	g.Eq(goe.GetList("LIST", []int{}), []int{1, 2})

	g.Eq(goe.GetList("DEFAULT", []int{1}), []int{1})
}

func TestRequireOneOf(t *testing.T) {
	g := got.T(t)
	list := []string{"apple", "banana"}

	t.Setenv("ONE", "apple")
	g.Eq(goe.RequireOneOf("ONE", list...), "apple")

	t.Setenv("TWO", "banana")
	g.Eq(goe.RequireOneOf("TWO", list...), "banana")

	t.Setenv("THREE", "orange")
	g.Eq(g.Panic(func() {
		goe.RequireOneOf("THREE", list...)
	}), "env variable 'THREE' value 'orange' not in the list: [apple banana]")
}

func TestGetMap(t *testing.T) {
	g := got.T(t)

	t.Setenv("MAP", "a:1,b:2")
	g.Eq(
		goe.GetMap("MAP", map[string]int{"a": 2, "c": 3}),
		map[string]int{"a": 1, "b": 2, "c": 3},
	)
}

func TestUnset(t *testing.T) {
	g := got.T(t)

	t.Setenv("ENV", "dev")
	goe.Unset("ENV")

	g.Eq(goe.Get("ENV", "stg"), "stg")
}
