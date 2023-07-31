package lexer

import (
	"github.com/forfd8960/simpleinterpreter/tokens"
)

func CondExp(condition bool, v1, v2 tokens.TokenType) tokens.TokenType {
	if condition {
		return v1
	}
	return v2
}
