package tokens

type Token struct {
	TkType TokenType
	Literal string
	Value interface{}
}

func NewToken(tkType TokenType, literal string, value interface{}) *Token {
	return &Token{
		TkType: tkType,
		Literal: literal,
		Value: value,
	}
}

type Scanner struct {
	Text string
	Tokens []*Token
	start int
	current int
}

func (tk *Scanner) ScanTokens() {}

func (tk *Scanner) parseInt() (*Token, error)
