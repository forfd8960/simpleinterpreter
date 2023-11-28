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
	DPlus    // ++
	DMinus   // --
	MINUS    // -
	BANG     // !
	ASTERISK // *
	POW      // **
	SLASH    // /
	LT       // <
	LTEQ     // <=
	GT       // >
	GTEQ     // >=

	EQUAL    // ==
	NOTEQUAL // !=

	OR  // ||
	AND // &&

	COMMA     // ,
	SEMICOLON // ;
	DOT       // .

	LPRARENT   // (
	RPARENT    // )
	LBRACE     // {
	RBRACE     // }
	LSQBRACKET // left square bracket [
	RSQBRACKET // right square breacket ]

	CLASS    // class
	THIS     // this
	FUNCTION // fn
	LET      // let
	IF       // if
	ELSE     // else
	RETURN   // return
	TRUE     // true
	FALSE    // false
	NIL      // nil

	FOR   // for
	WHILE // while
	PRINT // print()
	BREAK // break
)

const (
	LiteralEOF = "eof"
)
