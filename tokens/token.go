package tokens

import (
	"fmt"
	"strings"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
}

var keyword2Token = map[string]*Token{
	"let": {
		TkType:  LET,
		Literal: "let",
		Value:   "let",
	},
	"fn": &Token{
		TkType:  FUNCTION,
		Literal: "fn",
		Value:   "function",
	},
	"if": &Token{
		TkType:  IF,
		Literal: "fn",
		Value:   "function",
	},
	"else": &Token{
		TkType:  ELSE,
		Literal: "fn",
		Value:   "function",
	},
	"return": &Token{
		TkType:  RETURN,
		Literal: "return",
		Value:   "return",
	},
	"true": &Token{
		TkType:  TRUE,
		Literal: "true",
		Value:   true,
	},
	"false": &Token{
		TkType:  FALSE,
		Literal: "false",
		Value:   false,
	},
	"nil": &Token{
		TkType:  NIL,
		Literal: "nil",
		Value:   nil,
	},
}

type Token struct {
	TkType  TokenType
	Literal string
	Value   interface{}
}

func NewToken(tkType TokenType, literal string, value interface{}) *Token {
	return &Token{
		TkType:  tkType,
		Literal: literal,
		Value:   value,
	}
}

func (tk *Token) String() string {
	return strings.Join([]string{
		"TkType: " + tk.TkType.String(),
		"Literal: " + tk.Literal,
		fmt.Sprintf("%v", tk.Value),
	}, ",")
}

func LookupIdent(ident string) TokenType {
	if tp, ok := keywords[ident]; ok {
		return tp
	}

	return IDENT
}

func LookupTokenByIdent(ident string) *Token {
	if tp, ok := keyword2Token[ident]; ok {
		return tp
	}

	return nil
}
