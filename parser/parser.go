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
	case p.match(tokens.CLASS):
		return p.parseClassStmt()
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
	if p.match(tokens.ASSIGN) {
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

func (p *Parser) parseClassStmt() (*ast.ClassStmt, error) {
	className, err := p.consume(tokens.IDENT, "expect class name")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(tokens.LBRACE, "expect { before class body"); err != nil {
		return nil, err
	}

	var methods = []*ast.Function{}
	for !p.isAtEnd() && !p.check(tokens.RBRACE) {
		fn, err := p.function()
		if err != nil {
			return nil, err
		}

		methods = append(methods, fn)
	}

	if _, err := p.consume(tokens.RBRACE, "expect } after class body"); err != nil {
		return nil, err
	}

	return ast.NewClassStmt(className, methods), nil
}

// Todo: anonymous functions
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
		return block, nil
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	p.consume(tokens.LPRARENT, "Expect `(` after for.")

	var initializer ast.Stmt
	var err error

	if p.match(tokens.SEMICOLON) {
		initializer = nil
	} else if p.match(tokens.LET) {
		initializer, err = p.parseLetStmt()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var cond ast.Expression
	if !p.check(tokens.SEMICOLON) {
		cond, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(tokens.SEMICOLON, "Expect `;` after loop condition."); err != nil {
		return nil, err
	}

	var increment ast.Expression
	if !p.check(tokens.SEMICOLON) {
		increment, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(tokens.RPARENT, "Expect `)` after for clauses."); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = ast.NewBlockStmt([]ast.Stmt{body, ast.NewExpressionStmt(increment)})
	}

	if cond == nil {
		cond = ast.NewLiteral(tokens.NewToken(tokens.TRUE, "true", true))
	}
	body = ast.NewWhileStmt(cond, body)

	if initializer != nil {
		body = ast.NewBlockStmt([]ast.Stmt{initializer, body})
	}

	return body, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	p.consume(tokens.LPRARENT, `expect "(" after 'if'`)
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	p.consume(tokens.RPARENT, `expect ")" after if condition `)

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(tokens.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return ast.NewIFStmt(cond, thenBranch, elseBranch), nil
}

// printStatement
func (p *Parser) printStatement() (ast.Stmt, error) {
	if _, err := p.consume(tokens.LPRARENT, "Expect `(` after print."); err != nil {
		return nil, err
	}

	values := []ast.Expression{}
	format, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	values = append(values, format)
	for {
		if !p.match(tokens.COMMA) {
			break
		}

		v, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}

	if _, err := p.consume(tokens.RPARENT, `Expect ) after print`); err != nil {
		return nil, err
	}

	if _, err := p.consume(tokens.SEMICOLON, `Expect ";" after print()`); err != nil {
		return nil, err
	}
	return ast.NewPrintStmt(values), nil
}

// whileStatement
func (p *Parser) whileStatement() (ast.Stmt, error) {
	p.consume(tokens.LPRARENT, `expect "(" after 'while'.`)
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	p.consume(tokens.RPARENT, `expect ")" after condition.`)

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return ast.NewWhileStmt(cond, body), nil
}

func (p *Parser) block() (*ast.Block, error) {
	statements := make([]ast.Stmt, 0, 1)

	for !p.check(tokens.RBRACE) && !p.isAtEnd() {
		d, err := p.declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, d)
	}

	if _, err := p.consume(tokens.RBRACE, `Expect "}" after block!`); err != nil {
		return nil, err
	}

	if p.check(tokens.SEMICOLON) {
		p.consume(tokens.SEMICOLON, `Expect ; after block!`)
	}

	return ast.NewBlockStmt(statements), nil
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

		switch v := exp.(type) {
		case *ast.Assign:
			return ast.NewAssign(v.Name, value), nil
		case *ast.Get:
			return ast.NewSet(v.Expr, v.Name, value), nil
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

	for p.match(tokens.GT, tokens.GTEQ, tokens.LT, tokens.LTEQ, tokens.EQUAL, tokens.NOTEQUAL) {
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

	return p.call()
}

func (p *Parser) call() (ast.Expression, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(tokens.LPRARENT) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(tokens.DOT) {
			name, err := p.consume(tokens.IDENT, "Expect property name after .")
			if err != nil {
				return nil, err
			}
			expr = ast.NewGet(expr, name)
		} else {
			break
		}
	}

	return expr, nil
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
	case p.match(tokens.THIS):
		return ast.NewThisExpr(p.previous()), nil
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

func (p *Parser) finishCall(callee ast.Expression) (ast.Expression, error) {
	arguments := make([]ast.Expression, 0, 10)
	if !p.check(tokens.RPARENT) {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, expr)

		for p.check(tokens.COMMA) {
			p.consume(tokens.COMMA, "expect comma after (")

			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)
		}
	}

	if _, err := p.consume(tokens.RPARENT, "Expect ) after call"); err != nil {
		return nil, err
	}

	return ast.NewCall(callee, arguments), nil
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
