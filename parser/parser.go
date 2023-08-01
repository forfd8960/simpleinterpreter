package parser

import (
	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

type Parser struct {
	lxer      *lexer.Lexer
	currentTk *tokens.Token
	nextTk    *tokens.Token
}

func NewParser(l *lexer.Lexer) (*Parser, error) {
	p := &Parser{lxer: l}

	if err := p.nextToken(); err != nil {
		return nil, err
	}
	if err := p.nextToken(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) nextToken() error {
	p.currentTk = p.nextTk
	var err error
	p.nextTk, err = p.lxer.NextToken()
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	return nil, nil
}
