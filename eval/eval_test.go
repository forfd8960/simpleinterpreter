package eval

import (
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
