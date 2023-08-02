package parser

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrNotSupportToken = "not supported token: %s"
)

type Parser struct {
	tokens    []*tokens.Token
	current   int
	errors    []error
	lxer      *lexer.Lexer
	currentTk *tokens.Token
	nextTk    *tokens.Token
}

func NewParser(tokens []*tokens.Token) *Parser {
	p := &Parser{tokens: tokens}
	return p
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{
		Stmts: make([]ast.Stmt, 1),
	}

	for !p.isAtEnd() {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
	}

	return program, nil
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	switch {
	case p.match(tokens.LET):
		return p.parseLetStmt()
	default:
		return nil, fmt.Errorf(ErrNotSupportToken)
	}

}

func (p *Parser) parseLetStmt() (ast.Stmt, error) {
	identToken, err := p.consume(tokens.IDENT, "expect identifier name")
	if err != nil {
		return nil, err
	}

	var initExpr ast.Expression
	if p.match(tokens.EQUAL) {
		initExpr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(tokens.SEMICOLON, "expect `;` after identifier name"); err != nil {
		return nil, err
	}

	return ast.NewLetStmt(identToken, initExpr), nil
}

func (p *Parser) parseExpr() (ast.Expression, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expression, error) {
	exp, err := p.or()
	if err != nil {
		return nil, err
	}
	//exp := p.equality()

	if p.match(tokens.ASSIGN) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := exp.(*ast.Identifier); ok {
			return ast.NewAssign(v.Token, value), nil
		}

		return nil, fmt.Errorf("invalid assignment: %v", equals)
	}

	return exp, nil
}

func (p *Parser) or() (ast.Expression, error) {
	exp, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.OR) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		exp = ast.NewLogical(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) and() (ast.Expression, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.AND) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		exp = ast.NewLogical(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) equality() (ast.Expression, error) {
	exp := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		op := p.previous()
		right := p.comparison()
		exp = expr.NewBinary(exp, right, op)
	}
	return exp
}

func (p *Parser) comparison() (ast.Expression, error) {
	exp := p.term()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		op := p.previous()
		right := p.term()
		exp = expr.NewBinary(exp, right, op)
	}
	return exp
}

func (p *Parser) term() (ast.Expression, error) {
	exp := p.factor()

	for p.match(token.Minus, token.Plus) {
		op := p.previous()
		right := p.factor()
		exp = expr.NewBinary(exp, right, op)
	}

	return exp
}

func (p *Parser) factor() (ast.Expression, error) {
	exp := p.unary()

	for p.match(token.Slash, token.Star) {
		op := p.previous()
		right := p.unary()
		exp = expr.NewBinary(exp, right, op)
	}

	return exp
}

func (p *Parser) unary() (ast.Expression, error) {
	if p.match(tokens.BANG, tokens.MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expr.NewUnary(op, right)
	}

	// return p.call()
	return nil, fmt.Errorf("not support call")
}

func (p *Parser) consume(tkType tokens.TokenType, msg string) (*tokens.Token, error) {
	if p.check(tkType) {
		return p.advance(), nil
	}

	return nil, fmt.Errorf(msg)
}

func (p *Parser) match(tkTypes ...tokens.TokenType) bool {
	for _, t := range tkTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tkType tokens.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TkType == tkType
}

func (p *Parser) advance() *tokens.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *tokens.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TkType == tokens.EOF
}

func (p *Parser) peek() *tokens.Token {
	return p.tokens[p.current]
}
