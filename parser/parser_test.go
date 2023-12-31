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

func TestParsePow(t *testing.T) {
	input := `
	let v = 1 + 2**5;
	`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	oneLiteral, _ := ast.NewLiteral1(1)
	twoLiteral, _ := ast.NewLiteral1(2)
	fiveLiteral, _ := ast.NewLiteral1(5)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		assert.Nil(t, err)
		assert.NotNil(t, program)
		if assert.Equal(t, 1, len(program.Stmts)) {
			assert.Equal(t,
				ast.NewLetStmt(
					tokens.NewIdentToken("v"),
					ast.NewBinary(oneLiteral, ast.NewBinary(
						twoLiteral,
						fiveLiteral,
						tokens.OPPow),
						tokens.OPPlus,
					),
				),
				program.Stmts[0],
			)
		}
	}
}

func TestParseFor(t *testing.T) {
	input := `for (i = 0; i <= 10; i=i+1) {
		print("%d", i);
	}`

	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	zeroLiteral, _ := ast.NewLiteral1(0)
	tenLiteral, _ := ast.NewLiteral1(10)
	oneLiteral, _ := ast.NewLiteral1(1)
	fmtLiteral, _ := ast.NewLiteral1("%d")

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
			assert.NotNil(t, program)

			assert.Equal(t,
				ast.NewBlockStmt([]ast.Stmt{
					ast.NewExpressionStmt(
						ast.NewAssign(tokens.NewIdentToken("i"), zeroLiteral),
					),
					ast.NewWhileStmt(
						ast.NewBinary(
							ast.NewIdentifier1("i"),
							tenLiteral,
							tokens.NewToken(tokens.LTEQ, "<=", "<="),
						),
						ast.NewBlockStmt([]ast.Stmt{
							ast.NewBlockStmt([]ast.Stmt{
								ast.NewPrintStmt([]ast.Expression{
									fmtLiteral,
									ast.NewIdentifier1("i"),
								}),
							}),
							ast.NewExpressionStmt(
								ast.NewAssign(tokens.NewIdentToken("i"), ast.NewBinary(
									ast.NewIdentifier1("i"),
									oneLiteral,
									tokens.OPPlus,
								)),
							),
						}),
					),
				}),
				program.Stmts[0],
			)
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
		if assert.Nil(t, err) {
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
}

func TestAssign(t *testing.T) {
	input := `
	let a = 1;
	a = a + 10;
	`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
			assert.NotNil(t, program)

			assert.Equal(t, []ast.Stmt{
				ast.NewLetStmt(
					tokens.NewToken(tokens.IDENT, "a", "a"),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "1", int64(1))),
				),
				ast.NewExpressionStmt(ast.NewAssign(
					tokens.NewToken(tokens.IDENT, "a", "a"),
					ast.NewBinary(
						ast.NewIdentifier1("a"),
						ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
						tokens.NewToken(tokens.PLUS, "+", "+"),
					),
				)),
			}, program.Stmts)
		}
	}
}

func TestIFStmt(t *testing.T) {
	input := `if (1 <= 2) {
		print("%d", 1);
	} else {
		print("%d", 22);
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
					ast.NewPrintStmt([]ast.Expression{
						ast.NewLiteral(
							tokens.NewToken(tokens.STRING, "%d", "%d"),
						),
						ast.NewLiteral(
							tokens.NewToken(tokens.INTEGER, "1", int64(1)),
						),
					}),
				}),
				ast.NewBlockStmt([]ast.Stmt{
					ast.NewPrintStmt([]ast.Expression{
						ast.NewLiteral(
							tokens.NewToken(tokens.STRING, "%d", "%d"),
						),
						ast.NewLiteral(
							tokens.NewToken(tokens.INTEGER, "22", int64(22)),
						),
					},
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

func TestParseBreakStmt(t *testing.T) {
	input := `while ( a < b) {
		a = a + 1;
		break;
	}`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	identA := tokens.NewToken(tokens.IDENT, "a", "a")
	identB := tokens.NewToken(tokens.IDENT, "b", "b")
	oneLiteral := tokens.NewToken(tokens.INTEGER, "1", int64(1))

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
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
					ast.NewBreakStmt(),
				}),
			), program.Stmts[0])
		}
	}
}

func TestParseSliceStmt(t *testing.T) {
	input := `
	let arr = [1,2,3];
	let arr1 = [true];
	let arr2 = [];
	`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	arr := tokens.NewToken(tokens.IDENT, "arr", "arr")
	arr1 := tokens.NewToken(tokens.IDENT, "arr1", "arr1")
	arr2 := tokens.NewToken(tokens.IDENT, "arr2", "arr2")
	oneLiteral, _ := ast.NewLiteral1(1)
	twoLiteral, _ := ast.NewLiteral1(2)
	threeLiteral, _ := ast.NewLiteral1(3)
	trueLiteral, _ := ast.NewLiteral1(true)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
			assert.NotNil(t, program)

			assert.Equal(t, ast.NewLetStmt(
				arr,
				ast.NewSlice(
					[]ast.Expression{
						oneLiteral,
						twoLiteral,
						threeLiteral,
					},
				),
			), program.Stmts[0])

			assert.Equal(t, ast.NewLetStmt(
				arr1,
				ast.NewSlice([]ast.Expression{trueLiteral}),
			), program.Stmts[1])

			assert.Equal(t,
				ast.NewLetStmt(
					arr2,
					ast.NewSlice([]ast.Expression{}),
				),
				program.Stmts[2],
			)
		}
	}
}

func TestParseSliceAccessStmt(t *testing.T) {
	input := `
	let arr = [1,2,3];
	arr[0]
	`
	tokenList, err := lexer.TokensFromInput(input)
	assert.Nil(t, err)

	arr := tokens.NewToken(tokens.IDENT, "arr", "arr")
	zeroLiteral, _ := ast.NewLiteral1(0)
	oneLiteral, _ := ast.NewLiteral1(1)
	twoLiteral, _ := ast.NewLiteral1(2)
	threeLiteral, _ := ast.NewLiteral1(3)

	p := NewParser(tokenList)
	if assert.NotNil(t, p) {
		program, err := p.ParseProgram()
		if assert.Nil(t, err) {
			assert.NotNil(t, program)

			assert.Equal(t, ast.NewLetStmt(
				arr,
				ast.NewSlice(
					[]ast.Expression{
						oneLiteral,
						twoLiteral,
						threeLiteral,
					},
				),
			), program.Stmts[0])

			assert.Equal(t,
				ast.NewExpressionStmt(
					ast.NewSliceAccess(
						ast.NewIdentifier1("arr"),
						zeroLiteral,
					),
				), program.Stmts[1])
		}
	}
}
