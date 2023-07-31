package lexer

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrUnSupportedToken = "unsupported token: %v"
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

func (l *Lexer) NextToken() (*tokens.Token, error) {
	if !l.isAtEnd() {
		l.start = l.current
		return l.scanToken()
	}

	return tokens.NewToken(tokens.EOF, tokens.LiteralEOF, nil), nil
}

func (l *Lexer) scanToken() (*tokens.Token, error) {
	r := l.advance()

	var (
		err error
		tok *tokens.Token
	)

	switch r {
	case '=':
		tok = l.buildToken(tokens.ASSIGN, "=")
	case ';':
		tok = l.buildToken(tokens.SEMICOLON, ";")
	case ',':
		tok = l.buildToken(tokens.COMMA, ",")
	case '+':
		tok = l.buildToken(tokens.PLUS, "+")
	case '(':
		tok = l.buildToken(tokens.LPRARENT, "(")
	case ')':
		tok = l.buildToken(tokens.RPARENT, ")")
	case '{':
		tok = l.buildToken(tokens.LBRACE, "{")
	case '}':
		tok = l.buildToken(tokens.RBRACE, "}")
	case ' ', '\r', '\n', '\t':
	default:
		if isDigit(r) {
			tok, err = l.parseInteger()
		} else if isAlpha(r) {
			tok, err = l.parseIdent()
		} else {
			err = fmt.Errorf(ErrUnSupportedToken, r)
		}
	}
	if err != nil {
		return nil, err
	}

	return tok, nil
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
	l.current++
	return l.runes[l.current-1]
}

func (l *Lexer) parseInteger() (*tokens.Token, error) {
	for isDigit(l.peek()) {
		l.advance()
	}

	text := string(l.runes[l.start:l.current])
	num, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return nil, err
	}

	return tokens.NewToken(tokens.INTEGER, text, num), nil
}

func (l *Lexer) parseIdent() (*tokens.Token, error) {
	return nil, nil
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

func (l *Lexer) peek() rune {
	if l.current >= len(l.runes) {
		return '\000'
	}

	return l.runes[l.current]
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
