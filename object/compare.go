package object

import "github.com/forfd8960/simpleinterpreter/tokens"

func (obj *Integer) Compare(op tokens.TokenType, other *Integer) bool {
	switch op {
	case tokens.GT:
		return obj.Value > other.Value
	case tokens.GTEQ:
		return obj.Value >= other.Value
	case tokens.LT:
		return obj.Value < other.Value
	case tokens.LTEQ:
		return obj.Value <= other.Value
	case tokens.NOTEQUAL:
		return obj.Value != other.Value
	case tokens.EQUAL:
		return obj.Value == other.Value
	}

	return false
}

func (obj *String) Compare(op tokens.TokenType, other *String) bool {
	switch op {
	case tokens.GT:
		return obj.Value > other.Value
	case tokens.GTEQ:
		return obj.Value >= other.Value
	case tokens.LT:
		return obj.Value < other.Value
	case tokens.LTEQ:
		return obj.Value <= other.Value
	case tokens.NOTEQUAL:
		return obj.Value != other.Value
	case tokens.EQUAL:
		return obj.Value == other.Value
	}

	return false
}
