package eval

import (
	"fmt"
	"math"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrNodeNotLiteral                = "node: %v is not literal"
	ErrNotIntegerValue               = "value: %v is not integer"
	ErrNotBoolValue                  = "value: %v is not boolean"
	ErrNotStringValue                = "value: %v is not string"
	ErrDivideByZero                  = "integer divide by zero"
	ErrNotSupportedOperator          = "operator is not supported: %v"
	ErrIdentifierNotFound            = "identifier: %s is not found"
	ErrIdentifierIsNotCallable       = "%s is not callable(it shoud be function or xxx)"
	ErrOnlyClassInstanceHaveProperty = "expr: %s can not get property, only class instance have property"
	ErrThisNotFoundClassInstance     = "this can not found the class instance"
	ErrCondMustBeBoolValue           = "condition must be a bool value: %v"
)

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch v := node.(type) {
	case *ast.Program:
		return evalStatements(v.Stmts, env)
	case *ast.ClassStmt:
		return evalClassStmt(v, env)
	case *ast.Function:
		return evalFunctionStmt(v, env)
	case *ast.Block:
		return evalBockStmts(v, env)
	case *ast.WhileStmt:
		return evalWhileStmt(v, env)
	case *ast.LetStmt:
		return evalLetStmt(v, env)
	case *ast.Assign:
		return evalAssign(v, env)
	case *ast.Identifier:
		return evalIdent(v, env)
	case *ast.IFStmt:
		return evalIfStmt(v, env)
	case *ast.ReturnStmt:
		return evalReturn(v, env)
	case *ast.PrintStmt:
		return evalPrintStmt(v, env)
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
	case *ast.DExp: // a--, b++
		return evalDExp(v, env)
	case *ast.Get: // classObj.property
		return evalGetStmt(v, env)
	case *ast.Set: // classObj.property = value
		return evalSetStmt(v, env)
	case *ast.ThisExpr:
		return evalThisExpr(v, env)
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

func evalClassStmt(cls *ast.ClassStmt, env *object.Environment) (*object.Class, error) {
	methods := make(map[string]*object.Function, len(cls.Methods))
	for _, fn := range cls.Methods {
		method, err := evalFunctionStmt(fn, env)
		if err != nil {
			return nil, err
		}

		methods[fn.Name.Literal] = method
	}

	clsObj := object.NewClass(cls.NameIdent.Name, methods, object.NewEnvWithOutter(env))
	// store class object into env
	env.Set(clsObj.Name, clsObj)

	return clsObj, nil
}

// Todo: anonymous functions
func evalFunctionStmt(astFn *ast.Function, env *object.Environment) (*object.Function, error) {
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
	env.Set(astFn.Name.Literal, fn)
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

func evalWhileStmt(wl *ast.WhileStmt, env *object.Environment) (object.Object, error) {
	var result object.Object
	for {
		cond, err := Eval(wl.Condition, env)
		if err != nil {
			return nil, err
		}

		truth, ok := cond.(*object.Bool)
		if !ok {
			return nil, fmt.Errorf(ErrCondMustBeBoolValue, cond)
		}

		if !truth.Value {
			break
		}

		result, err = Eval(wl.Body, env)
		if err != nil {
			return nil, err
		}

		ret, isRet := result.(*object.Return)
		if isRet {
			return ret, nil
		}
	}

	return result, nil
}

func evalLetStmt(let *ast.LetStmt, env *object.Environment) (object.Object, error) {
	obj, err := Eval(let.InitExpr, env)
	if err != nil {
		return nil, err
	}

	env.Set(let.Ident.Name, obj)
	return obj, nil
}

func evalAssign(assign *ast.Assign, env *object.Environment) (object.Object, error) {
	obj, err := Eval(assign.Value, env)
	if err != nil {
		return nil, err
	}

	env.Set(assign.Name.Literal, obj)
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

func evalPrintStmt(v *ast.PrintStmt, env *object.Environment) (object.Object, error) {
	return evalBuiltInPrint(
		ast.NewCall(v, v.Values),
		env,
	)
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
	case tokens.PLUS:
		return plusObj(leftResult, rightResult)
	default:
		return doMath(leftResult, rightResult, op)
	}
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

func plusObj(obj1, obj2 object.Object) (object.Object, error) {
	if obj1.Type() != obj2.Type() {
		return nil, fmt.Errorf("can not plus 2 different type: %v, %v", obj1, obj2)
	}

	switch obj1.Type() {
	case object.OBJ_INTEGER:
		left, _ := obj1.(*object.Integer)
		right, _ := obj2.(*object.Integer)
		return &object.Integer{Value: left.Value + right.Value}, nil
	case object.OBJ_STRING:
		left, _ := obj1.(*object.String)
		right, _ := obj2.(*object.String)
		return &object.String{Value: left.Value + right.Value}, nil
	}

	return nil, fmt.Errorf("unsuported + for %s", obj1.Type())
}

func doMath(obj1, obj2 object.Object, op tokens.TokenType) (object.Object, error) {
	leftValue, ok := obj1.(*object.Integer)
	if !ok {
		return nil, fmt.Errorf("%v must be number", obj1.Inspect())
	}
	rightValue, ok := obj2.(*object.Integer)
	if !ok {
		return nil, fmt.Errorf("%v must be number", obj2.Inspect())
	}

	switch op {
	case tokens.MINUS:
		return &object.Integer{Value: leftValue.Value - rightValue.Value}, nil
	case tokens.ASTERISK:
		return &object.Integer{Value: leftValue.Value * rightValue.Value}, nil
	case tokens.POW:
		return &object.Integer{Value: int64(math.Pow(float64(leftValue.Value), float64(rightValue.Value)))}, nil
	case tokens.SLASH:
		if rightValue.Value == 0 {
			return nil, fmt.Errorf(ErrDivideByZero)
		}
		return &object.Integer{Value: leftValue.Value / rightValue.Value}, nil
	}

	return nil, fmt.Errorf("unsupported operator: %s", op.String())
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

	switch v := callee.(type) {
	case *object.Class:
		return evalCallClass(v, callExpr, globalEnv)
	case *object.Function:
		return evalCallFunction(v, callExpr, globalEnv)
	}

	return nil, fmt.Errorf(ErrIdentifierIsNotCallable, callee.Inspect())
}

func evalDExp(dexp *ast.DExp, env *object.Environment) (object.Object, error) {
	switch dexp.Operator.TkType {
	case tokens.DPlus:
		return evalBinary(
			ast.NewBinary(
				dexp.Left,
				ast.NewLiteral(
					tokens.NewToken(tokens.INTEGER, "1", int64(1)),
				),
				tokens.OPPlus,
			),
			env,
		)
	case tokens.DMinus:
		return evalBinary(
			ast.NewBinary(
				dexp.Left,
				ast.NewLiteral(
					tokens.NewToken(tokens.INTEGER, "1", int64(1)),
				),
				tokens.OPMinus,
			),
			env,
		)
	}

	return nil, fmt.Errorf("unsupported operator: %+v", dexp.Operator)
}

func evalCallClass(cls *object.Class, callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	return object.NewClassInstance(cls), nil
}

func evalCallFunction(fn *object.Function, callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	if len(fn.Parameters) != len(callExpr.Arguments) {
		return nil, fmt.Errorf("not engough params to function: %s, need %d arguments", fn.Inspect(), len(fn.Parameters))
	}

	var env = object.NewEnvWithOutter(fn.Env)
	if len(fn.Parameters) > 0 {
		for idx, param := range fn.Parameters {
			v, err := Eval(callExpr.Arguments[idx], globalEnv)
			if err != nil {
				return nil, err
			}

			fmt.Printf("setting %s with value: %+v\n", param.Name, v)
			env.Set(param.Name, v)
		}
	}

	obj, err := Eval(fn.Body, env)
	if err != nil {
		return nil, err
	}

	var result object.Object
	if obj != nil {
		result = obj
		for result.Type() == object.OBJ_RETURN {
			v := result.(*object.Return)
			result = v.Value
		}
	}

	return result, nil
}

func evalGetStmt(get *ast.Get, env *object.Environment) (object.Object, error) {
	instanceObj, err := Eval(get.Expr, env)
	if err != nil {
		return nil, err
	}

	clsInst, ok := instanceObj.(*object.ClassInstance)
	if !ok {
		return nil, fmt.Errorf(ErrOnlyClassInstanceHaveProperty, instanceObj.Inspect())
	}

	v, err := clsInst.Get(get.Name)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func evalSetStmt(set *ast.Set, env *object.Environment) (object.Object, error) {
	obj, err := Eval(set.Expr, env)
	if err != nil {
		return nil, err
	}

	clsInstance, ok := obj.(*object.ClassInstance)
	if !ok {
		return nil, fmt.Errorf("obj: %s is not class instance, only instance has fields", obj.Inspect())
	}

	value, err := Eval(set.Value, env)
	if err != nil {
		return nil, err
	}

	clsInstance.Set(set.Name, value)
	return value, nil
}

func evalThisExpr(kw *ast.ThisExpr, env *object.Environment) (object.Object, error) {
	// this should called in the method body
	// then can get the class instance from the env
	expr, ok := env.Get(kw.TokenLiteral())
	if !ok {
		return nil, fmt.Errorf(ErrThisNotFoundClassInstance)
	}
	return expr, nil
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
