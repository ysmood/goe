package goe_test

import (
	"testing"
	"time"

	"github.com/ysmood/goe"
	"github.com/ysmood/got"
)

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
}

func TestLoad(t *testing.T) {
	g := got.T(t)

	g.E(goe.Load(false))

	g.Eq(goe.Get("STR", ""), "hello")
}
