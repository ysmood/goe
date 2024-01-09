package dotenv_test

import (
	"testing"

	"github.com/ysmood/goe/pkg/dotenv"
	"github.com/ysmood/got"
)

func TestBasic(t *testing.T) {
	g := got.T(t)

	doc, err := dotenv.Parse(`
a = "ok"
b = 'ok'
# comment
export c = 10
d = ok # comment

e = 'o\nk'
f = 'o
k'
g = "o\nk"
h = "o
k"
i = ""
j = ''
k = ${test}
l = $test

123 = 123
ok_ok = 1
OK_OK = 2
ok_OK = 3

# test
# test
`)
	g.E(err)

	g.Snapshot("basic", doc)

	g.Eq(doc.Get("k"), "${test}")
	g.Eq(doc.Get("xxx"), "")
}
