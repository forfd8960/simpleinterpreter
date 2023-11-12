package eval

import (
	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/object"
)

// Resolve Semantic Analysis
// tracks down which declaration it refers to—each and every time the variable expression is evaluated
// https://craftinginterpreters.com/resolving-and-binding.html#a-resolver-class
/*
 If we could ensure a variable lookup always walked the same number of links in the environment chain,
 that would ensure that it found the same variable in the same scope every time.
*/
func Resolve(node ast.Node, env *object.Environment, stack *Stack) {
}
