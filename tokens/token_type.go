package tokens

//go:generate go run github.com/dmarkham/enumer -type=TokenType
type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INTEGER
	STRING
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /
	LT       // <
	GT       // >
	EQUAL    // ==
	NOTEQUAL // !=

	OR  // ||
	AND // &&

	COMMA     // ,
	SEMICOLON // ;

	LPRARENT // (
	RPARENT  // )
	LBRACE   // {
	RBRACE   // }

	FUNCTION // fn
	LET      // let
	IF       // if
	ELSE     // else
	RETURN   // return
	TRUE     // true
	FALSE    // false
)

const (
	LiteralEOF = "eof"
)
