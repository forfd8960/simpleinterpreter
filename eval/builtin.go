package eval

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
)

const (
	builtInPrint = "print"
)

var (
	ErrLackParameter        = fmt.Errorf("lack parameters for print")
	ErrInvalidParameterType = fmt.Errorf("invalid parameters for print(first paramter msut be string)")
)

var builtInfunctions = map[string]struct{}{
	builtInPrint: {},
}

func IsBuiltInFunction(fnName string) bool {
	_, ok := builtInfunctions[fnName]
	return ok
}

func buildPrint(callExpr *ast.Call, env *object.Environment) *object.Function {
	identifiers := make([]*ast.Identifier, 0, len(callExpr.Arguments))
	identifiers = append(identifiers, ast.NewIdentifier1("format"))

	for i := 1; i < len(callExpr.Arguments); i++ {
		identifiers = append(identifiers, ast.NewIdentifier1(fmt.Sprintf("val%d", i)))
	}

	return &object.Function{
		Parameters: identifiers,
	}
}

func evalBuildInFunctions(callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	switch callExpr.TokenLiteral() {
	case builtInPrint:
		return evalBuiltInPrint(callExpr, globalEnv)
	}

	return nil, nil
}

func evalBuiltInPrint(callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	fn := buildPrint(callExpr, globalEnv)

	var format string
	var values []any
	for idx, param := range fn.Parameters {
		v, err := Eval(callExpr.Arguments[idx], globalEnv)
		if err != nil {
			return nil, err
		}

		if param.Name == "format" {
			format = v.(*object.String).Value
			continue
		}

		values = append(values, v)
	}

	_, err := fmt.Printf(format, values...)
	return nil, err
}
