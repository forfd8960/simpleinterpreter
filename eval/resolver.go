package eval

import (
	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
)

// Resolve Semantic Analysis
// tracks down which declaration it refers to—each and every time the variable expression is evaluated
// https://craftinginterpreters.com/resolving-and-binding.html#a-resolver-class
func Resolve(node ast.Node, env *object.Environment, stack *Stack) {
}
