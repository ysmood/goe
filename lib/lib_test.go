package lib_test

import (
	"testing"
	"time"

	"github.com/ysmood/dotenv/lib"
	"github.com/ysmood/got"
)

func TestGet(t *testing.T) {
	g := got.T(t)

	t.Setenv("BOOL", "true")
	t.Setenv("NUM", "2")
	t.Setenv("STR", "ok")
	t.Setenv("FLOAT", "1.2")
	t.Setenv("DURATION", "1m")

	g.Eq(lib.Get("BOOL", false), true)
	g.Eq(lib.Get("BOOL_DEFAULT", true), true)
	g.Eq(lib.Get("NUM", 0), 2)
	g.Eq(lib.Get("NUM_DEFAULT", 1), 1)
	g.Eq(lib.Get("STR", ""), "ok")
	g.Eq(lib.Get("STR_DEFAULT", "yes"), "yes")
	g.Eq(lib.Get("FLOAT", 0.0), 1.2)
	g.Eq(lib.Get("FLOAT_DEFAULT", 1.1), 1.1)
	g.Eq(lib.Get("DURATION", time.Second), time.Minute)
	g.Eq(lib.Get("DURATION_DEFAULT", time.Hour), time.Hour)

	g.Eq(lib.Require[int]("NUM"), 2)

	g.Eq(g.Panic(func() {
		lib.Require[int]("KEY")
	}), "required env variable not found: KEY")
}
