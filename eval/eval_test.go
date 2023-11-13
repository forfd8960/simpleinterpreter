package eval

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/lexer"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/parser"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

func testEvalInput(input string) (object.Object, error) {
	env := object.NewEnvironment()
	tokens, err := lexer.TokensFromInput(input)
	if err != nil {
		return nil, err
	}

	parser := parser.NewParser(tokens)
	root, err := parser.ParseProgram()
	if err != nil {
		return nil, err
	}

	obj, err1 := Eval(root, env)
	if err1 != nil {
		return nil, err1
	}

	return obj, nil
}

func TestIntegrateEval(t *testing.T) {
	type args struct {
		input string
	}

	env := object.NewEnvironment()
	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval let",
			args: args{
				input: `let x = 5;`,
			},
			want:    &object.Integer{Value: 5},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := lexer.TokensFromInput(tt.args.input)
			for _, tk := range tokens {
				fmt.Printf("token: %s\n", tk)
			}

			fmt.Printf("lexer token err: %v\n", err)

			if assert.Nil(t, err) {
				parser := parser.NewParser(tokens)
				root, err := parser.ParseProgram()

				fmt.Printf("%+v\n", root)
				fmt.Printf("parser err: %v\n", err)

				if assert.Nil(t, err) {
					obj, err := Eval(root, env)
					assert.Nil(t, err)
					assert.Equal(t, tt.want, obj)
				}
			}
		})
	}
}

func TestEval(t *testing.T) {
	type args struct {
		input ast.Node
	}

	env := object.NewEnvironment()
	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval integer",
			args: args{
				input: ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
			},
			want: &object.Integer{Value: 64},
		},
		{
			name: "eval bool",
			args: args{
				input: ast.NewLiteral(tokens.NewToken(tokens.TRUE, "true", true)),
			},
			want: &object.Bool{Value: true},
		},
		{
			name: "eval string",
			args: args{
				input: ast.NewLiteral(tokens.NewToken(tokens.STRING, "ss", "ss")),
			},
			want: &object.String{Value: "ss"},
		},
		{
			name: "eval null",
			args: args{
				input: ast.NewLiteral(tokens.LookupTokenByIdent("null")),
			},
			want: &object.Null{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
		})
	}
}

func TestEvalFunction(t *testing.T) {
	type args struct {
		input string
	}

	fnEnv := object.NewEnvironment()
	fnEnv.Set("add", &object.Function{
		Parameters: []*ast.Identifier{
			ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
		},
		Body: ast.NewBlockStmt([]ast.Stmt{
			ast.NewReturnStmt(
				tokens.LookupTokenByIdent(tokens.KWReturn),
				ast.NewBinary(
					ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
					tokens.NewToken(tokens.PLUS, "+", "+"),
				),
			),
		}),
		Env: object.NewEnvironment(),
	})

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval function",
			args: args{
				input: `
				fn add(x) { return x + 10; }
				`,
			},
			want: &object.Function{
				Parameters: []*ast.Identifier{
					ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
				},
				Body: ast.NewBlockStmt([]ast.Stmt{
					ast.NewReturnStmt(
						tokens.LookupTokenByIdent(tokens.KWReturn),
						ast.NewBinary(
							ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
							ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
							tokens.NewToken(tokens.PLUS, "+", "+"),
						),
					),
				}),
				Env: fnEnv,
			},
			wantErr: false,
		},
		{
			name: "eval return call function",
			args: args{
				input: `
				fn add(x) { return x + 10; }
				return add(10);
				`,
			},
			want:    &object.Integer{Value: int64(20)},
			wantErr: false,
		},
		{
			name: "eval stack call function",
			args: args{
				input: `
				fn add(x) { return x + 10; }
				fn minus(x) { return x - 1; }
				fn divTwo(x) { return x / 2; }
				let x = 20;
				return add(minus(divTwo(x)));
				`,
			},
			want:    &object.Integer{Value: int64(19)},
			wantErr: false,
		},
		{
			name: "unhappy - eval function not enough arguments",
			args: args{
				input: `
				fn add(x, y) { return x + y; }
				return add(10);
				`,
			},
			want:    &object.Error{Message: "not engough params to function: fn(x,y), need 2 arguments"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := testEvalInput(tt.args.input)
			if tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, obj)
				return
			}

			if assert.NotNil(t, obj) {
				objFn, ok := obj.(*object.Function)
				wantFn, ok1 := tt.want.(*object.Function)
				if ok && ok1 {
					assert.Equal(t, wantFn.Parameters, objFn.Parameters)
					assert.Equal(t, wantFn.Body, objFn.Body)
					fn1, _ := wantFn.Env.Get("add")
					fn2, _ := objFn.Env.Get("add")
					assert.Equal(t, fn1.Inspect(), fn2.Inspect())
				} else {
					assert.Equal(t, tt.want, obj)
				}
			}
		})
	}
}

func TestEvalLet(t *testing.T) {
	type args struct {
		input ast.Node
	}

	env := object.NewEnvironment()

	wantEnv1 := object.NewEnvironment()
	wantEnv1.Set("x", &object.Integer{Value: 64})

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantEnv *object.Environment
		wantErr bool
	}{
		{
			name: "eval let",
			args: args{
				input: ast.NewLetStmt(
					tokens.NewToken(tokens.IDENT, "x", "x"),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
				),
			},
			wantEnv: wantEnv1,
			want:    &object.Integer{Value: 64},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			fmt.Printf("%+v\n", obj)
			fmt.Printf("type: %+v\n", reflect.TypeOf(obj))

			envObj, ok := env.Get("x")
			fmt.Printf("envObj: %+v, %+v\n", envObj, ok)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
			assert.Equal(t, tt.wantEnv, env)
		})
	}
}

func TestEval1(t *testing.T) {
	type args struct {
		input ast.Node
	}

	env := object.NewEnvironment()

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval integer",
			args: args{
				input: &ast.Program{
					Stmts: []ast.Stmt{
						&ast.ExpressionStmt{
							Expr: ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
						},
					},
				},
			},
			want: &object.Integer{Value: 64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			fmt.Printf("%+v\n", obj)
			fmt.Printf("type: %+v\n", reflect.TypeOf(obj))
			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
		})
	}
}

func TestEvalIdentifier(t *testing.T) {
	type args struct {
		input ast.Node
	}

	env := object.NewEnvironment()
	env.Set("x", &object.Integer{Value: 100})

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval identifier",
			args: args{
				input: &ast.Program{
					Stmts: []ast.Stmt{
						ast.NewIdentifier(tokens.NewToken(tokens.IDENT, "x", "x")),
					},
				},
			},
			want:    &object.Integer{Value: 100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
		})
	}
}

func TestEvalBinary(t *testing.T) {
	type args struct {
		input ast.Node
	}
	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval integer",
			args: args{
				input: ast.NewBinary(
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
					tokens.NewToken(tokens.PLUS, "+", "+"),
				),
			},
			want: &object.Integer{Value: 128},
		},
		{
			name: "eval compare",
			args: args{
				input: ast.NewBinary(
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "64", int64(64))),
					tokens.NewToken(tokens.EQUAL, "==", "=="),
				),
			},
			want: &object.Bool{Value: true},
		},
		{
			name: "eval compare string",
			args: args{
				input: ast.NewBinary(
					ast.NewLiteral(tokens.NewToken(tokens.STRING, "Hello", "Hello")),
					ast.NewLiteral(tokens.NewToken(tokens.STRING, "Hello", "Hello")),
					tokens.NewToken(tokens.EQUAL, "==", "=="),
				),
			},
			want: &object.Bool{Value: true},
		},
		{
			name: "eval compare string",
			args: args{
				input: ast.NewBinary(
					ast.NewLiteral(tokens.NewToken(tokens.STRING, "abc", "abc")),
					ast.NewLiteral(tokens.NewToken(tokens.STRING, "ABC", "ABC")),
					tokens.NewToken(tokens.GT, ">", ">"),
				),
			},
			want: &object.Bool{Value: true},
		},
	}
	env := object.NewEnvironment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			fmt.Printf("%+v\n", obj)
			fmt.Printf("type: %+v\n", reflect.TypeOf(obj))
			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
		})
	}
}

func TestEvalUnary(t *testing.T) {
	type args struct {
		input ast.Node
	}
	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval unary",
			args: args{
				input: ast.NewUnary(
					tokens.NewToken(tokens.BANG, "!", "!"),
					ast.NewLiteral(tokens.NewToken(tokens.TRUE, "true", true)),
				),
			},
			want: &object.Bool{Value: false},
		},
		{
			name: "eval unary",
			args: args{
				input: ast.NewUnary(
					tokens.NewToken(tokens.MINUS, "-", "-"),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
				),
			},
			want: &object.Integer{Value: -10},
		},
		{
			name: "eval unary err",
			args: args{
				input: ast.NewUnary(
					tokens.NewToken(tokens.BANG, "!", "!"),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
				),
			},
			wantErr: true,
			want:    nil,
		},
	}

	env := object.NewEnvironment()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			fmt.Printf("obj: %v, err: %v\n", obj, err)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			if assert.Nil(t, err) {
				assert.Equal(t, tt.want, obj)
			}
		})
	}
}

func TestEvalReturn(t *testing.T) {
	type args struct {
		input ast.Node
	}

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "eval return bool",
			args: args{
				input: ast.NewReturnStmt(
					tokens.NewToken(tokens.RETURN, "return", "return"),
					ast.NewLiteral(tokens.NewToken(tokens.TRUE, "true", true)),
				),
			},
			want: &object.Return{Value: &object.Bool{Value: true}},
		},
		{
			name: "eval return integer",
			args: args{
				input: ast.NewReturnStmt(
					tokens.NewToken(tokens.RETURN, "return", "return"),
					ast.NewLiteral(tokens.NewToken(tokens.INTEGER, "10", int64(10))),
				),
			},
			want: &object.Return{Value: &object.Integer{Value: 10}},
		},
		{
			name: "eval return unary",
			args: args{
				input: ast.NewReturnStmt(
					tokens.NewToken(tokens.RETURN, "return", "return"),
					ast.NewUnary(
						tokens.NewToken(tokens.BANG, "!", "!"),
						ast.NewLiteral(tokens.NewToken(tokens.TRUE, "true", true)),
					),
				),
			},
			want: &object.Return{Value: &object.Bool{Value: false}},
		},
		{
			name: "eval return string",
			args: args{
				input: ast.NewReturnStmt(
					tokens.NewToken(tokens.RETURN, "return", "return"),
					ast.NewLiteral(tokens.NewToken(tokens.STRING, "hello", "hello")),
				),
			},
			wantErr: false,
			want:    &object.Return{Value: &object.String{Value: "hello"}},
		},
	}
	env := object.NewEnvironment()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input, env)
			fmt.Printf("obj: %+v, err: %v\n", obj, err)
			fmt.Printf("obj type: %+v\n", reflect.ValueOf(obj).String())

			v := obj.(*object.Return).Value
			fmt.Printf("value: %+v\n", v)
			fmt.Printf("value type: %+v\n", reflect.ValueOf(v).String())
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			if assert.Nil(t, err) {
				assert.Equal(t, tt.want, obj)
			}
		})
	}
}

func TestStringCompare(t *testing.T) {
	s1 := "abc"
	s2 := "ABC"
	fmt.Printf("%t\n", s1 > s2)
	assert.True(t, s1 > s2)
}

func TestGroupExpression(t *testing.T) {
	var x = 100
	{
		var x = x + 10
		fmt.Println("x: ", x)
	}
	assert.Equal(t, 100, x)
}

func TestEvalBlock(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				input: `
				let a = 100;
				{
					let a = a + 10;
					return a
				}
				`,
			},
			want:    &object.Integer{Value: int64(110)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := testEvalInput(tt.args.input)
			if tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, obj)
				return
			}

			if assert.NotNil(t, obj) {
				assert.Equal(t, tt.want, obj)
			}
		})
	}
}

func TestEvalPrint(t *testing.T) {
	type args struct {
		input string
	}

	tests := []struct {
		name    string
		args    args
		want    object.Object
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				input: `
				print("%d\n", 100)
				`,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "happy path",
			args: args{
				input: `
				fn add(x, y) {
					return x + y;
				}
				print("%d\n", add(10, 10));
				`,
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := testEvalInput(tt.args.input)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}

			assert.Equal(t, tt.want, obj)
		})
	}
}
