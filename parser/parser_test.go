package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

func TestParseProgram(t *testing.T) {
	input := `
	let x = 5;
	let y = 6;
	let foobar = 989858;
	`
	l := lexer.NewLexer(input)
	tokenList := make([]*tokens.Token, 0, 1)
	for {
		tk, err := l.NextToken()
		if err != nil {
			break
		}

		tokenList = append(tokenList, tk)
		if tk.TkType == tokens.EOF {
			break
		}
	}
	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)
		assert.Equal(t, 3, len(program.Stmts))
	}
}

func TestParseLiteral(t *testing.T) {
	input := `
	true;
	false;
	100;
	`
	l := lexer.NewLexer(input)
	tokenList := make([]*tokens.Token, 0, 1)
	for {
		tk, err := l.NextToken()
		if err != nil {
			break
		}

		tokenList = append(tokenList, tk)
		if tk.TkType == tokens.EOF {
			break
		}
	}

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)

		if assert.Equal(t, 3, len(program.Stmts)) {
			assert.Equal(t,
				ast.NewExpressionStmt(ast.NewLiteral(tokens.NewToken(
					tokens.TRUE,
					"true",
					true,
				))),
				program.Stmts[0],
			)
			assert.Equal(t,
				ast.NewExpressionStmt(ast.NewLiteral(tokens.NewToken(
					tokens.FALSE,
					"false",
					false,
				))),
				program.Stmts[1],
			)
			assert.Equal(t,
				ast.NewExpressionStmt(ast.NewLiteral(tokens.NewToken(
					tokens.INTEGER,
					"100",
					int64(100),
				))),
				program.Stmts[2],
			)
		}
	}
}