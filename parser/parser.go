package parser

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrNotSupportToken = "not supported token: %s"
)

type Parser struct {
	tokens  []*tokens.Token
	current int
}

func NewParser(tokens []*tokens.Token) *Parser {
	p := &Parser{tokens: tokens}
	return p
}

/*
expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary

	| primary ;

primary        → NUMBER | STRING | "true" | "false" | "nil"

	| "(" expression ")" ;
*/
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{
		Stmts: make([]ast.Stmt, 0, 1),
	}

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}

		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
	}

	return program, nil
}

func (p *Parser) declaration() (ast.Stmt, error) {
	switch {
	case p.match(tokens.LET):
		return p.parseLetStmt()
	case p.match(tokens.FUNCTION):
		return p.function()
	default:
		return p.statement()
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

func (p *Parser) function() (*ast.Function, error) {
	name, err := p.consume(tokens.IDENT, "expect function name.")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokens.LPRARENT, "Expect `(` after function name."); err != nil {
		return nil, err
	}

	var params []*tokens.Token
	if !p.check(tokens.RPARENT) {
		ident, err := p.consume(tokens.IDENT, "Expect parameter name.")
		if err != nil {
			return nil, err
		}
		params = append(params, ident)

		for p.match(tokens.COMMA) {
			if len(params) > 8 {
				return nil, fmt.Errorf("cannot have more than 8 parameters")
			}

			ident, err = p.consume(tokens.IDENT, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, ident)
		}
	}

	if _, err := p.consume(tokens.RPARENT, "Expect `)` after parameters"); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokens.LBRACE, "Expect `{` before function body."); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return ast.NewFunctionStmt(name, params, body), nil
}

func (p *Parser) parseReturnStmt() (ast.Stmt, error) {
	kw := p.previous()

	var value ast.Expression
	var err error
	if !p.match(tokens.SEMICOLON) {
		value, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	if _, err = p.consume(tokens.SEMICOLON, "expect `;` after return value"); err != nil {
		return nil, err
	}

	return ast.NewReturnStmt(kw, value), nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	switch {
	case p.match(tokens.FOR):
		return p.forStatement()
	case p.match(tokens.IF):
		return p.ifStatement()
	case p.match(tokens.RETURN):
		return p.parseReturnStmt()
	case p.match(tokens.PRINT):
		return p.printStatement()
	case p.match(tokens.WHILE):
		return p.whileStatement()
	case p.match(tokens.LBRACE):
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return ast.NewBlockStmt(block), nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	return nil, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	return nil, nil
}

// printStatement
func (p *Parser) printStatement() (ast.Stmt, error) {
	return nil, nil
}

// whileStatement
func (p *Parser) whileStatement() (ast.Stmt, error) {
	return nil, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0, 1)

	for !p.check(tokens.RBRACE) && !p.isAtEnd() {
		d, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, d)
	}

	p.consume(tokens.RBRACE, `Expect "}" after block!`)
	return statements, nil
}

func (p *Parser) expressionStatement() (*ast.ExpressionStmt, error) {
	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	p.consume(tokens.SEMICOLON, `Expect ":" after value.`)
	return ast.NewExpressionStmt(value), nil
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
	exp, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.NOTEQUAL, tokens.EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) comparison() (ast.Expression, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.GT, tokens.GT, tokens.LT, tokens.LTEQ) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		exp = ast.NewBinary(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) term() (ast.Expression, error) {
	exp, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.MINUS, tokens.PLUS) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		exp = ast.NewBinary(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) factor() (ast.Expression, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.SLASH, tokens.ASTERISK) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, right, op)
	}

	return exp, nil
}

func (p *Parser) unary() (ast.Expression, error) {
	if p.match(tokens.BANG, tokens.MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return ast.NewUnary(op, right), nil
	}

	// return p.call()
	return p.primary()
}

func (p *Parser) primary() (ast.Expression, error) {
	switch {
	case p.match(tokens.FALSE):
		return ast.NewLiteral(p.previous()), nil
	case p.match(tokens.TRUE):
		return ast.NewLiteral(p.previous()), nil
	case p.match(tokens.NIL):
		return ast.NewLiteral(p.previous()), nil
	case p.match(tokens.INTEGER, tokens.STRING):
		return ast.NewLiteral(p.previous()), nil
	case p.match(tokens.IDENT):
		return ast.NewIdentifier(p.previous()), nil
	case p.match(tokens.LPRARENT):
		exp, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(tokens.RPARENT, "Expect ')' after expression"); err != nil {
			return nil, err
		}
		return ast.NewGrouping(exp), nil

	}

	return nil, fmt.Errorf("unknow expr")
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