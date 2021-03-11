package parser

import (
	"fmt"
	"github.com/alecthomas/participle/lexer"

	"github.com/evilsocket/stork/commands"
)

type AST struct {
	Steps []*Statement `@@*`
}

type Statement struct {
	Command *Command `( @@`
	Set     *Set     `| @@)`
}

type Set struct {
	Pos lexer.Position

	Identifier string `@Variable`
	Value      *Value `"=" @@`
}

type Command struct {
	Pos lexer.Position

	Identifier string   `@Ident`
	Parameters []*Value `( @@ ( "," @@ )* )?`
}

type Value struct {
	Pos lexer.Position

	String   *string `  @String`
	Variable *string `| @Variable`
}

func (p Value) Resolve(vars commands.Variables) (string, error) {
	if p.String != nil {
		return *p.String, nil
	} else if val, found := vars[(*p.Variable)[1:]]; found {
		return val, nil
	} else {
		return "", fmt.Errorf("%v %s undefined", p.Pos, *p.Variable)
	}
}
