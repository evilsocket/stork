package parser

import (
	"github.com/alecthomas/participle/lexer"
)

var storkLexer = lexer.Must(lexer.Regexp(
	`(?m)` +
		`(\s+)` +
		`|(^[#].*$)` +
		`|(?P<Punct>[,=])` +
		`|(?P<Variable>\$[a-zA-Z][a-zA-Z_\d]*)` +
		`|(?P<Ident>[a-zA-Z][a-zA-Z_\:\d]*)` +
		`|(?P<String>(?:(?:"(?:\\.|[^\"])*")|(?:'(?:\\.|[^'])*')))`,
))
