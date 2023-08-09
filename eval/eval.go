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
	case *ast.Block:
		return evalBockStmts(v)
	case *ast.IFStmt:
		return evalIfStmt(v)
	case *ast.ReturnStmt:
		return evalReturn(v)
	case *ast.ExpressionStmt:
		return Eval(v.Expr)
	case *ast.Literal:
		return evalLiteral(v)
	case *ast.Binary:
		return evalBinary(v)
	case *ast.Unary:
		return evalUnary(v)
	}

	return nil, nil
}

func evalStatements(nodes []ast.Stmt) (object.Object, error) {
	var result object.Object
	var err error
	for _, stmt := range nodes {
		result, err = Eval(stmt)

		if ret, ok := result.(*object.Return); ok {
			return ret.Value, nil
		}

		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalBockStmts(b *ast.Block) (object.Object, error) {
	var obj object.Object
	var err error
	for _, stmt := range b.Statements {
		obj, err = Eval(stmt)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

func evalIfStmt(v *ast.IFStmt) (object.Object, error) {
	cond, err := Eval(v.Condition)
	if err != nil {
		return nil, err
	}
	truth, ok := cond.(*object.Bool)
	if !ok {
		return nil, fmt.Errorf("bad if condition: %v", cond)
	}

	if truth.Value {
		return Eval(v.ThenBranch)
	} else if v.ElseBranch != nil {
		return Eval(v.ElseBranch)
	} else {
		return nil, nil
	}
}

func evalReturn(v *ast.ReturnStmt) (object.Object, error) {
	result, err := Eval(v.Value)
	if err != nil {
		return nil, err
	}

	return &object.Return{Value: result}, nil
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

	op := bin.Operator.TkType

	switch op {
	case tokens.GT, tokens.GTEQ, tokens.LT, tokens.LTEQ, tokens.NOTEQUAL, tokens.EQUAL:
		resul, err := compareObj(leftResult, rightResult, op)
		return &object.Bool{Value: resul}, err
	default:
		leftValue := leftResult.(*object.Integer)
		rightValue := rightResult.(*object.Integer)
		switch op {
		case tokens.PLUS:
			return &object.Integer{Value: leftValue.Value + rightValue.Value}, nil
		case tokens.MINUS:
			return &object.Integer{Value: leftValue.Value - rightValue.Value}, nil
		case tokens.ASTERISK:
			return &object.Integer{Value: leftValue.Value * rightValue.Value}, nil
		case tokens.SLASH:
			return &object.Integer{Value: leftValue.Value / rightValue.Value}, nil
		}

	}

	return nil, fmt.Errorf(ErrNotSupportedOperator, bin.Operator.Literal)
}

func compareObj(obj1, obj2 object.Object, op tokens.TokenType) (bool, error) {
	if obj1.Type() != obj2.Type() {
		return false, fmt.Errorf("can not compare 2 different type: %v, %v", obj1, obj2)
	}

	switch obj1.Type() {
	case object.OBJ_INTEGER:
		left, _ := obj1.(*object.Integer)
		right, _ := obj2.(*object.Integer)
		return compareInteger(left.Value, right.Value, op), nil
	case object.OBJ_STRING:
		left, _ := obj1.(*object.String)
		right, _ := obj2.(*object.String)
		return compareString(left.Value, right.Value, op), nil
	}

	return false, fmt.Errorf("unsupported compare type:  %v, %v", obj1, obj2)
}

func compareInteger(v1, v2 int64, op tokens.TokenType) bool {
	switch op {
	case tokens.GT:
		return v1 > v2
	case tokens.GTEQ:
		return v1 >= v2
	case tokens.LT:
		return v1 < v2
	case tokens.LTEQ:
		return v1 <= v2
	case tokens.NOTEQUAL:
		return v1 != v2
	case tokens.EQUAL:
		return v1 == v2
	}

	return false
}

func compareString(v1, v2 string, op tokens.TokenType) bool {
	switch op {
	case tokens.GT:
		return v1 > v2
	case tokens.GTEQ:
		return v1 >= v2
	case tokens.LT:
		return v1 < v2
	case tokens.LTEQ:
		return v1 <= v2
	case tokens.NOTEQUAL:
		return v1 != v2
	case tokens.EQUAL:
		return v1 == v2
	}

	return false
}

func evalUnary(node *ast.Unary) (object.Object, error) {
	op := node.Operator
	obj, err := Eval(node.Right)
	if err != nil {
		return nil, err
	}

	switch op.TkType {
	case tokens.BANG:
		v, ok := obj.(*object.Bool)
		if !ok {
			return nil, fmt.Errorf("right value must be boolean: %v", obj)
		}
		return &object.Bool{Value: !v.Value}, nil
	case tokens.MINUS:
		v, ok := obj.(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("right value must be integer: %v", obj)
		}
		return &object.Integer{Value: -v.Value}, nil
	}

	return nil, fmt.Errorf("unsupported unary Operator: %v", op)
}

func evalLiteralInteger(value interface{}) (*object.Integer, error) {
	v, ok := value.(int64)
	if !ok {
		return nil, fmt.Errorf(ErrNotIntegerValue, value)
	}
	return &object.Integer{Value: v}, nil
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
