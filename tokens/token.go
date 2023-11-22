package tokens

import (
	"fmt"
	"strings"
)

const (
	KWLet    = "let"
	KWReturn = "return"
	KWPrint  = "print"
	KWClass  = "class"
	KWThis   = "this"
	KWFn     = "fn"
	KWIf     = "if"
	KWElse   = "else"
	KwFor    = "for"
	KwWhile  = "while"
	KWTrue   = "true"
	KWFlase  = "false"
	KWNull   = "null"
)

var (
	Assign = NewToken(ASSIGN, "=", "=")
)

var keywords = map[string]TokenType{
	KWLet:    LET,
	KWClass:  CLASS,
	KWThis:   THIS,
	KWFn:     FUNCTION,
	KWIf:     IF,
	KWElse:   ELSE,
	KwFor:    FOR,
	KWReturn: RETURN,
	KWTrue:   TRUE,
	KWFlase:  FALSE,
	KWPrint:  PRINT,
	KWNull:   NIL,
}

var keyword2Token = map[string]*Token{
	KWLet: {
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
	KWFn: {
		TkType:  FUNCTION,
		Literal: KWFn,
		Value:   "function",
	},
	KWIf: {
		TkType:  IF,
		Literal: "if",
		Value:   "if",
	},
	KwFor: {
		TkType:  FOR,
		Literal: "for",
		Value:   "for",
	},
	KwWhile: {
		TkType:  WHILE,
		Literal: "while",
		Value:   "while",
	},
	KWElse: {
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
	KWTrue: {
		TkType:  TRUE,
		Literal: "true",
		Value:   true,
	},
	KWFlase: {
		TkType:  FALSE,
		Literal: "false",
		Value:   false,
	},
	KWNull: {
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
