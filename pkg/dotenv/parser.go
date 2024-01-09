package dotenv

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Doc struct {
	Lines []*Line `@@*`
}

type Line struct {
	Key   string `(("export" @Ident) | @Ident) "="`
	Value string `@(Ident | String)`
}

var (
	docLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Ident`, `([a-zA-Z_\d]+)|(\$\{[a-zA-Z_\d]+\})|(\$[a-zA-Z_\d]+)`},
		{`String`, `("[^"]*")|('[^']*')`},
		{`Punct`, `=`},
		{`export`, `export`},
		{"comment", `#[^\n]*`},
		{"whitespace", `\s+`},
	})

	parser = participle.MustBuild[Doc](
		participle.Lexer(docLexer),
		participle.Unquote("String"),
	)
)

func Parse(doc string) (*Doc, error) {
	return parser.ParseString("", doc) //nolint: wrapcheck
}

func (d *Doc) Get(key string) string {
	for _, line := range d.Lines {
		if line.Key == key {
			return line.Value
		}
	}

	return ""
}
