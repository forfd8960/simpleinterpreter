package ast

import (
	"fmt"
	"reflect"

	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrUnsuportedLiteralType = "unsuported literal type: %v"
)

func NewLiteral1(literal any) (*Literal, error) {
	switch v := literal.(type) {
	case int64:
		return NewLiteral(tokens.NewIntegerToken(v)), nil
	case int:
		return NewLiteral(tokens.NewIntegerToken(int64(v))), nil
	case int32:
		return NewLiteral(tokens.NewIntegerToken(int64(v))), nil
	case bool:
		return NewLiteral(tokens.NewBoolToken(v)), nil
	case string:
		return NewLiteral(tokens.NewStringToken(v)), nil
	default:
		return nil, fmt.Errorf(ErrUnsuportedLiteralType, reflect.TypeOf(literal))
	}
}
