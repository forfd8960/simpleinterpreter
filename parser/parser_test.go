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

func TestWhileStmt(t *testing.T) {
	input := `while ( a < b) {
		a = a + 1;
	}`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	identA := tokens.NewToken(tokens.IDENT, "a", "a")
	identB := tokens.NewToken(tokens.IDENT, "b", "b")
	oneLiteral := tokens.NewToken(tokens.INTEGER, "1", int64(1))

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)

		assert.Equal(t, ast.NewWhileStmt(
			ast.NewBinary(
				ast.NewIdentifier(identA),
				ast.NewIdentifier(identB),
				tokens.NewToken(tokens.LT, "<", "<"),
			),
			ast.NewBlockStmt([]ast.Stmt{
				ast.NewExpressionStmt(
					ast.NewAssign(
						identA,
						ast.NewBinary(
							ast.NewIdentifier(identA),
							ast.NewLiteral(oneLiteral),
							tokens.NewToken(tokens.PLUS, "+", "+"),
						),
					),
				),
			}),
		), program.Stmts[0])
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

func TestParseFunction(t *testing.T) {
	input := `
	fn add(x) {
		return x + 1;
	}
	`

	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
			if assert.NotNil(t, program) {
				if assert.Equal(t, 1, len(program.Stmts)) {
					assert.Equal(t, ast.NewFunctionStmt(
						tokens.NewToken(tokens.IDENT, "add", "add"),
						[]*tokens.Token{
							tokens.NewToken(tokens.IDENT, "x", "x"),
						},
						ast.NewBlockStmt([]ast.Stmt{
							ast.NewReturnStmt(
								tokens.LookupTokenByIdent(tokens.KWReturn),
								ast.NewBinary(
									ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
									ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "1", int64(1))),
									tokens.NewToken(tokens.PLUS, "+", "+"),
								),
							),
						}),
					),
						program.Stmts[0],
					)
				}
			}
		}
	}
}

func TestParseCallExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		program *ast.Program
		err     error
	}{
		{
			name:  "parse one function call",
			input: "add(1, 2);",
			program: &ast.Program{
				Stmts: []ast.Stmt{
					ast.NewExpressionStmt(ast.NewCall(
						ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "add", "add")),
						[]ast.Expression{
							ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "1", int64(1))),
							ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "2", int64(2))),
						},
					)),
				},
			},
			err: nil,
		},
		{
			name:  "parse stacked function call",
			input: "add(minus(1), div(2));",
			program: &ast.Program{
				Stmts: []ast.Stmt{
					ast.NewExpressionStmt(ast.NewCall(
						ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "add", "add")),
						[]ast.Expression{
							ast.NewCall(
								ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "minus", "minus")),
								[]ast.Expression{ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "1", int64(1)))},
							),
							ast.NewCall(
								ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "div", "div")),
								[]ast.Expression{ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "2", int64(2)))},
							),
						},
					)),
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenList, err := lexer.TokensFromInput(tt.input)
			assert.Nil(t, err)

			p := NewParser(tokenList)
			if assert.NotNil(t, p) {
				program, err := p.ParseProgram()
				if assert.Equal(t, tt.err, err) {
					if assert.NotNil(t, program) {
						if assert.Equal(t, 1, len(program.Stmts)) {
							fmt.Printf("%+v\n", program.Stmts[0])
							assert.Equal(t, tt.program, program)
						}
					}
				}
			}
		})
	}
}
