package parser

import (
	"github.com/alecthomas/participle"
	"io/ioutil"
)

type Parser struct {
	parser *participle.Parser
}

func New() (parser *Parser, err error) {
	parser = &Parser{}
	parser.parser, err = participle.Build(&AST{},
		participle.Lexer(storkLexer),
		participle.Unquote("String"),
	)
	return
}

func (p *Parser) ParseCode(code string) (*AST, error) {
	ast := &AST{}
	err := p.parser.ParseString(code, ast)
	if err != nil {
		return nil, err
	}
	return ast, nil
}

func (p *Parser) ParseFile(fileName string) (*AST, error) {
	ast := &AST{}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = p.parser.ParseBytes(content, ast)
	if err != nil {
		return nil, err
	}
	return ast, nil
}