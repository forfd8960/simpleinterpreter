package eval

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

func TestEval(t *testing.T) {
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
			obj, err := Eval(tt.args.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, obj)
		})
	}
}

func TestEval1(t *testing.T) {
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
			obj, err := Eval(tt.args.input)
			fmt.Printf("%+v\n", obj)
			fmt.Printf("type: %+v\n", reflect.TypeOf(obj))
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Eval(tt.args.input)
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
