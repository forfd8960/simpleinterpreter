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
