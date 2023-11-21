package tokens

import (
	"fmt"
	"strings"
)

const (
	KWReturn = "return"
	KWPrint  = "print"
	KWClass  = "class"
	KWThis   = "this"
)

var (
	Assign = NewToken(ASSIGN, "=", "=")
)

var keywords = map[string]TokenType{
	"let":    LET,
	KWClass:  CLASS,
	KWThis:   THIS,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	KWReturn: RETURN,
	"true":   TRUE,
	"false":  FALSE,
	KWPrint:  PRINT,
	"null":   NIL,
}

var keyword2Token = map[string]*Token{
	"let": {
		TkType:  LET,
		Literal: "let",
		Value:   "let",
	},
	KWClass: {
		TkType:  CLASS,
		Literal: "class",
		Value:   "class",
	},
	KWThis: {
		TkType:  THIS,
		Literal: "this",
		Value:   "this",
	},
	"fn": {
		TkType:  FUNCTION,
		Literal: "fn",
		Value:   "function",
	},
	"if": {
		TkType:  IF,
		Literal: "if",
		Value:   "if",
	},
	"for": {
		TkType:  FOR,
		Literal: "for",
		Value:   "for",
	},
	"while": {
		TkType:  WHILE,
		Literal: "while",
		Value:   "while",
	},
	"else": {
		TkType:  ELSE,
		Literal: "else",
		Value:   "else",
	},
	KWReturn: {
		TkType:  RETURN,
		Literal: "return",
		Value:   "return",
	},
	KWPrint: {
		TkType:  PRINT,
		Literal: "print",
		Value:   "print",
	},
	"true": {
		TkType:  TRUE,
		Literal: "true",
		Value:   true,
	},
	"false": {
		TkType:  FALSE,
		Literal: "false",
		Value:   false,
	},
	"null": {
		TkType:  NIL,
		Literal: "null",
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
