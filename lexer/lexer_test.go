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
	}

	lexer := NewLexer(input)

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			token := lexer.NextToken()
			assert.Equal(t, tt.expectType, token.TkType)
			assert.Equal(t, tt.literal, token.Literal)
		})
	}
}
