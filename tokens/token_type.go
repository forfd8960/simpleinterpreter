package tokens

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INTEGER
	STRING
	ASSIGN
	PLUS

	COMMA
	SEMICOLON

	LPRARENT // (
	RPARENT  // )
	LBRACE   // {
	RBRACE   // }

	FUNCTION // FN
	LET      // let
)

const (
	LiteralEOF = "eof"
)
