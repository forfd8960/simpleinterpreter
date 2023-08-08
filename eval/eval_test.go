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
			want: &object.Integrer{Value: 64},
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
			want: &object.Integrer{Value: 64},
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
			want: &object.Integrer{Value: 128},
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
			want: &object.Integrer{Value: -10},
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
