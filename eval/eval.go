package eval

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrNodeNotLiteral       = "node: %v is not literal"
	ErrNotIntegerValue      = "value: %v is not integer"
	ErrNotBoolValue         = "value: %v is not boolean"
	ErrNotStringValue       = "value: %v is not string"
	ErrNotSupportedOperator = "operator is not supported: %v"
)

func Eval(node ast.Node) (object.Object, error) {
	switch v := node.(type) {
	case *ast.Program:
		return evalStatements(v.Stmts)
	case *ast.ExpressionStmt:
		return Eval(v.Expr)
	case *ast.Literal:
		return evalLiteral(v)
	case *ast.Binary:
		return evalBinary(v)
	}

	return nil, nil
}

func evalStatements(nodes []ast.Stmt) (object.Object, error) {
	var result object.Object
	var err error
	for _, stmt := range nodes {
		result, err = Eval(stmt)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalLiteral(literal *ast.Literal) (object.Object, error) {
	value := literal.Value
	switch value.TkType {
	case tokens.INTEGER:
		return evalLiteralInteger(value.Value)
	case tokens.TRUE, tokens.FALSE:
		return evalLiteralBool(value.Value)
	case tokens.STRING:
		return evalLiteralString(value.Value)
	case tokens.NIL:
		return evalLiteralNull(value.Value)
	}

	return nil, fmt.Errorf(ErrNodeNotLiteral, value)
}

func evalBinary(bin *ast.Binary) (object.Object, error) {
	leftResult, err := Eval(bin.Left)
	if err != nil {
		return nil, err
	}
	rightResult, err := Eval(bin.Right)
	if err != nil {
		return nil, err
	}

	leftValue := leftResult.(*object.Integrer)
	rightValue := rightResult.(*object.Integrer)

	switch bin.Operator.TkType {
	case tokens.PLUS:
		return &object.Integrer{Value: leftValue.Value + rightValue.Value}, nil
	case tokens.MINUS:
		return &object.Integrer{Value: leftValue.Value - rightValue.Value}, nil
	case tokens.ASTERISK:
		return &object.Integrer{Value: leftValue.Value * rightValue.Value}, nil
	case tokens.SLASH:
		return &object.Integrer{Value: leftValue.Value / rightValue.Value}, nil
	}

	return nil, fmt.Errorf(ErrNotSupportedOperator, bin.Operator.Literal)
}

func evalLiteralInteger(value interface{}) (*object.Integrer, error) {
	v, ok := value.(int64)
	if !ok {
		return nil, fmt.Errorf(ErrNotIntegerValue, value)
	}
	return &object.Integrer{Value: v}, nil
}

func evalLiteralBool(value interface{}) (*object.Bool, error) {
	v, ok := value.(bool)
	if !ok {
		return nil, fmt.Errorf(ErrNotBoolValue, value)
	}
	return &object.Bool{Value: v}, nil
}

func evalLiteralString(value interface{}) (*object.String, error) {
	v, ok := value.(string)
	if !ok {
		return nil, fmt.Errorf(ErrNotStringValue, value)
	}

	return &object.String{Value: v}, nil
}

func evalLiteralNull(value interface{}) (*object.Null, error) {
	if value == nil {
		return &object.Null{}, nil
	}

	return nil, fmt.Errorf("%t is not null", value)
}
