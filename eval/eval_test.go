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
