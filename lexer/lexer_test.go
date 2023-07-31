package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/forfd8960/simpleinterpreter/tokens"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"
	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.ASSIGN, "="},
		{tokens.PLUS, "+"},
		{tokens.LPRARENT, "("},
		{tokens.RPARENT, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RBRACE, "}"},
		{tokens.COMMA, ","},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, "eof"},
	}

	lexer := NewLexer(input)

	for _, tt := range tests {
		t.Run("test-"+tt.literal, func(t *testing.T) {
			token, err := lexer.NextToken()
			assert.Nil(t, err)
			assert.Equal(t, tt.expectType, token.TkType)
			assert.Equal(t, tt.literal, token.Literal)
		})
	}
}
