package lexer

import (
	"unicode"

	"github.com/forfd8960/simpleinterpreter/tokens"
)

type Lexer struct {
	input   string
	runes   []rune
	start   int // current index in input
	current int // current read index in input after pos
	char    rune
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		runes: []rune(input),
	}
}

func (l *Lexer) NextToken() *tokens.Token {
	return nil
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.runes)
}

func (l *Lexer) buildToken(tkType tokens.TokenType, value interface{}) *tokens.Token {
	literal := string(l.runes[l.start:l.current])
	return &tokens.Token{
		TkType:  tkType,
		Literal: literal,
		Value:   value,
	}
}

func (l *Lexer) advance() rune {
	pos := l.current
	l.current++
	return l.runes[pos]
}

func (l *Lexer) match(r rune) bool {
	if l.isAtEnd() {
		return false
	}

	if l.runes[l.current] != r {
		return false
	}

	l.advance()
	return true
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isAlphaNumberic(r rune) bool {
	return isAlpha(r) || isDigit(r)
}
