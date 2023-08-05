package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

func TestParseLetStmt(t *testing.T) {
	input := `
	let x = 5;
	let y = 6;
	let foobar = 989858;
	`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)
		if assert.Equal(t, 3, len(program.Stmts)) {
			assert.Equal(t, ast.NewLetStmt(
				tokens.NewToken(
					tokens.IDENT,
					"x",
					"x",
				),
				ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "5", int64(5))),
			), program.Stmts[0])

			assert.Equal(t, ast.NewLetStmt(
				tokens.NewToken(
					tokens.IDENT,
					"y",
					"y",
				),
				ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "6", int64(6))),
			), program.Stmts[1])

			assert.Equal(t, ast.NewLetStmt(
				tokens.NewToken(
					tokens.IDENT,
					"foobar",
					"foobar",
				),
				ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "989858", int64(989858))),
			), program.Stmts[2])
		}
	}
}

func TestIFStmt(t *testing.T) {
	input := `if (1 <= 2) {
		print(1);
	} else {
		print(22);
	}`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	for _, tk := range tokenList {
		fmt.Printf("token: %s\n", tk)
	}

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)

		assert.Equal(t,
			ast.NewIFStmt(
				ast.NewBinary(
					ast.NewLiteral(
						tokens.NewToken(tokens.INTEGER, "1", int64(1)),
					),
					ast.NewLiteral(
						tokens.NewToken(tokens.INTEGER, "2", int64(2)),
					),
					tokens.NewToken(
						tokens.LTEQ,
						"<=",
						"<=",
					),
				),
				ast.NewBlockStmt([]ast.Stmt{
					ast.NewPrintStmt(
						ast.NewGrouping(
							ast.NewLiteral(
								tokens.NewToken(tokens.INTEGER, "1", int64(1)),
							),
						),
					),
				}),
				ast.NewBlockStmt([]ast.Stmt{
					ast.NewPrintStmt(
						ast.NewGrouping(
							ast.NewLiteral(
								tokens.NewToken(tokens.INTEGER, "22", int64(22)),
							),
						),
					),
				}),
			),
			program.Stmts[0],
		)
	}
}

func TestParseLiteral(t *testing.T) {
	input := `
	true;
	false;
	100;
	`

	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

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
