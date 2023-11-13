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

func TestLexerLetStmt(t *testing.T) {
	input := "let x = 5;"
	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "x"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "5"},
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

func TestComplexTokens(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	};
	let result = add(five, ten);
	`
	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "five"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "5"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "ten"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "add"},
		{tokens.ASSIGN, "="},
		{tokens.FUNCTION, "fn"},
		{tokens.LPRARENT, "("},
		{tokens.IDENT, "x"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "y"},
		{tokens.RPARENT, ")"},
		{tokens.LBRACE, "{"},
		{tokens.IDENT, "x"},
		{tokens.PLUS, "+"},
		{tokens.IDENT, "y"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "result"},
		{tokens.ASSIGN, "="},
		{tokens.IDENT, "add"},
		{tokens.LPRARENT, "("},
		{tokens.IDENT, "five"},
		{tokens.COMMA, ","},
		{tokens.IDENT, "ten"},
		{tokens.RPARENT, ")"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, tokens.LiteralEOF},
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

func TestMoreTokens(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	if (five < ten) {
		return true;
	} else {
		return false;
	};
	`
	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "five"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "5"},
		{tokens.SEMICOLON, ";"},
		{tokens.LET, "let"},
		{tokens.IDENT, "ten"},
		{tokens.ASSIGN, "="},
		{tokens.INTEGER, "10"},
		{tokens.SEMICOLON, ";"},
		{tokens.IF, "if"},
		{tokens.LPRARENT, "("},
		{tokens.IDENT, "five"},
		{tokens.LT, "<"},
		{tokens.IDENT, "ten"},
		{tokens.RPARENT, ")"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.TRUE, "true"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.ELSE, "else"},
		{tokens.LBRACE, "{"},
		{tokens.RETURN, "return"},
		{tokens.FALSE, "false"},
		{tokens.SEMICOLON, ";"},
		{tokens.RBRACE, "}"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, tokens.LiteralEOF},
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

func TestMathOperator(t *testing.T) {
	input := `+ / * -;
	=;
	==;
	!
	!=;
	`
	lexer := NewLexer(input)

	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.PLUS, "+"},
		{tokens.SLASH, "/"},
		{tokens.ASTERISK, "*"},
		{tokens.MINUS, "-"},
		{tokens.SEMICOLON, ";"},
		{tokens.ASSIGN, "="},
		{tokens.SEMICOLON, ";"},
		{tokens.EQUAL, "=="},
		{tokens.SEMICOLON, ";"},
		{tokens.BANG, "!"},
		{tokens.NOTEQUAL, "!="},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, tokens.LiteralEOF},
	}

	for _, tt := range tests {
		t.Run("test-"+tt.literal, func(t *testing.T) {
			token, err := lexer.NextToken()

			assert.Nil(t, err)
			assert.Equal(t, tt.expectType, token.TkType)
			assert.Equal(t, tt.literal, token.Literal)
		})
	}
}

func TestLexerString(t *testing.T) {
	input := `
	let a = "abc";
	`
	lexer := NewLexer(input)
	// for !lexer.isAtEnd() {
	// 	token, err := lexer.NextToken()
	// 	if err != nil {
	// 		fmt.Println("lexer err: ", err)
	// 		break
	// 	}
	// 	fmt.Printf("token: %+v\n", token)
	// }

	tests := []struct {
		expectType tokens.TokenType
		literal    string
	}{
		{tokens.LET, "let"},
		{tokens.IDENT, "a"},
		{tokens.ASSIGN, "="},
		{tokens.STRING, "abc"},
		{tokens.SEMICOLON, ";"},
		{tokens.EOF, tokens.LiteralEOF},
	}

	for _, tt := range tests {
		t.Run("test-"+tt.literal, func(t *testing.T) {
			token, err := lexer.NextToken()

			assert.Nil(t, err)
			assert.Equal(t, tt.expectType, token.TkType)
			assert.Equal(t, tt.literal, token.Literal)
		})
	}
}
