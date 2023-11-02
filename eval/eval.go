package eval

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrNodeNotLiteral          = "node: %v is not literal"
	ErrNotIntegerValue         = "value: %v is not integer"
	ErrNotBoolValue            = "value: %v is not boolean"
	ErrNotStringValue          = "value: %v is not string"
	ErrDivideByZero            = "integer divide by zero"
	ErrNotSupportedOperator    = "operator is not supported: %v"
	ErrIdentifierNotFound      = "identifier: %s is not found"
	ErrIdentifierIsNotCallable = "%s is not callable(it shoud be function or xxx)"
)

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch v := node.(type) {
	case *ast.Program:
		return evalStatements(v.Stmts, env)
	case *ast.Function:
		return evalFunctionStmt(v, env)
	case *ast.Block:
		return evalBockStmts(v, env)
	case *ast.LetStmt:
		return evalLetStmt(v, env)
	case *ast.Identifier:
		return evalIdent(v, env)
	case *ast.IFStmt:
		return evalIfStmt(v, env)
	case *ast.ReturnStmt:
		return evalReturn(v, env)
	case *ast.ExpressionStmt:
		return Eval(v.Expr, env)
	case *ast.Grouping:
		return evalGroup(v, env)
	case *ast.Literal:
		return evalLiteral(v)
	case *ast.Binary:
		return evalBinary(v, env)
	case *ast.Unary:
		return evalUnary(v, env)
	case *ast.Call:
		return evalCall(v, env)
	}

	return nil, nil
}

func evalStatements(nodes []ast.Stmt, env *object.Environment) (object.Object, error) {
	var result object.Object
	var err error
	for _, stmt := range nodes {
		result, err = Eval(stmt, env)

		if ret, ok := result.(*object.Return); ok {
			return ret.Value, nil
		}

		if err != nil {
			return newError(err.Error()), nil
		}
	}
	return result, nil
}

func evalFunctionStmt(astFn *ast.Function, env *object.Environment) (object.Object, error) {
	params := make([]*ast.Identifier, 0, len(astFn.Parameters))
	for _, token := range astFn.Parameters {
		params = append(params, ast.NewIdentifier(token))
	}

	fn := &object.Function{
		Parameters: params,
		Body:       astFn.Body,
		Env:        env,
	}

	// register to env, and call expression can find the function object later
	env.Set(astFn.Name.String(), fn)
	return fn, nil
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
	}
}

func evalBockStmts(b *ast.Block, env *object.Environment) (object.Object, error) {
	var obj object.Object
	var err error
	for _, stmt := range b.Statements {
		obj, err = Eval(stmt, env)
		if err != nil {
			return nil, err
		}

		if obj != nil && obj.Type() == object.OBJ_RETURN {
			return obj, nil
		}
	}

	return obj, nil
}

func evalLetStmt(let *ast.LetStmt, env *object.Environment) (object.Object, error) {
	obj, err := Eval(let.InitExpr, env)
	if err != nil {
		return nil, err
	}

	env.Set(let.Ident.Name, obj)
	return obj, nil
}

func evalIdent(ident *ast.Identifier, env *object.Environment) (object.Object, error) {
	obj, ok := env.Get(ident.Name)
	if !ok {
		return nil, fmt.Errorf(ErrIdentifierNotFound, ident.Name)
	}

	return obj, nil
}

func evalIfStmt(v *ast.IFStmt, env *object.Environment) (object.Object, error) {
	cond, err := Eval(v.Condition, env)
	if err != nil {
		return nil, err
	}
	truth, ok := cond.(*object.Bool)
	if !ok {
		return nil, fmt.Errorf("bad if condition: %v", cond)
	}

	if truth.Value {
		return Eval(v.ThenBranch, env)
	} else if v.ElseBranch != nil {
		return Eval(v.ElseBranch, env)
	} else {
		return nil, nil
	}
}

func evalReturn(v *ast.ReturnStmt, env *object.Environment) (object.Object, error) {
	result, err := Eval(v.Value, env)
	if err != nil {
		return nil, err
	}

	return &object.Return{Value: result}, nil
}

func evalGroup(g *ast.Grouping, env *object.Environment) (object.Object, error) {
	return Eval(g.Expr, env)
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

func evalBinary(bin *ast.Binary, env *object.Environment) (object.Object, error) {
	leftResult, err := Eval(bin.Left, env)
	if err != nil {
		return nil, err
	}
	rightResult, err := Eval(bin.Right, env)
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
			if rightValue.Value == 0 {
				return nil, fmt.Errorf(ErrDivideByZero)
			}
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
		return left.Compare(op, right), nil
	case object.OBJ_STRING:
		left, _ := obj1.(*object.String)
		right, _ := obj2.(*object.String)
		return left.Compare(op, right), nil
	}

	return false, fmt.Errorf("unsupported compare type:  %v, %v", obj1.Type(), obj2.Type())
}

func evalUnary(node *ast.Unary, env *object.Environment) (object.Object, error) {
	op := node.Operator
	obj, err := Eval(node.Right, env)
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

func evalCall(callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	// callExpr.Callee is a identifier, and after Eval, it should return a function object
	callee, err := Eval(callExpr.Callee, globalEnv)
	if err != nil {
		return nil, err
	}

	fn, ok := callee.(*object.Function)
	if !ok {
		return nil, fmt.Errorf(ErrIdentifierIsNotCallable, callee.Inspect())
	}

	arguments := make([]object.Object, 0, len(fn.Parameters))
	for _, param := range fn.Parameters {
		v, err := Eval(param, globalEnv)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, v)
	}

	var env = object.NewEnvWithOutter(fn.Env)
	if len(arguments) > 0 {
		for idx, arg := range arguments {
			env.Set(fn.Parameters[idx].Name, arg)
		}
	}

	return Eval(fn.Body, env)
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
