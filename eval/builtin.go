package eval

import (
	"fmt"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
)

const (
	builtInPrint = "print"
	argFormat    = "format"
)

var (
	ErrLackParameter        = fmt.Errorf("lack parameters for print")
	ErrInvalidParameterType = fmt.Errorf("invalid parameters for print(first paramter msut be string)")
)

var builtInfunctions = map[string]fnBuilder{
	builtInPrint: buildPrint,
}

type fnBuilder func(callExpr *ast.Call, env *object.Environment) *object.Function

func IsBuiltInFunction(fnName string) bool {
	_, ok := builtInfunctions[fnName]
	return ok
}

var buildPrint = func(callExpr *ast.Call, env *object.Environment) *object.Function {
	identifiers := make([]*ast.Identifier, 0, len(callExpr.Arguments))
	identifiers = append(identifiers, ast.NewIdentifier1(argFormat))

	for i := 1; i < len(callExpr.Arguments); i++ {
		identifiers = append(identifiers, ast.NewIdentifier1(fmt.Sprintf("val%d", i)))
	}

	return &object.Function{
		Parameters: identifiers,
	}
}

func evalBuiltInPrint(callExpr *ast.Call, globalEnv *object.Environment) (object.Object, error) {
	if len(callExpr.Arguments) < 1 {
		return nil, ErrLackParameter
	}

	var format string
	var values []any

	fn := buildPrint(callExpr, globalEnv)
	for idx, param := range fn.Parameters {
		v, err := Eval(callExpr.Arguments[idx], globalEnv)
		if err != nil {
			return nil, err
		}

		if param.Name == argFormat {
			fmtVal, ok := v.(*object.String)
			if !ok {
				return nil, ErrInvalidParameterType
			}

			format = fmtVal.Value
			continue
		}

		values = append(values, getValueLiteral(v))
	}

	_, err := fmt.Printf(format, values...)
	return nil, err
}

func getValueLiteral(v object.Object) any {
	switch data := v.(type) {
	case *object.Integer:
		return data.Value
	case *object.String:
		return data.Value
	case *object.Bool:
		return data.Value
	case *object.Null:
		return nil
	default:
		return data
	}
}
