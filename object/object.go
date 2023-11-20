package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/forfd8960/simpleinterpreter/ast"
	"github.com/forfd8960/simpleinterpreter/tokens"
)

var (
	ErrPropertyNotFound = "property: %s not found for instance: %s"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Class struct {
	Name    string
	Methods map[string]*Function
	Env     *Environment
}

func (cls *Class) Inspect() string {
	return cls.Name
}

func (cls *Class) Type() ObjectType {
	return OBJ_CLASS
}

func NewClass(name string, methods map[string]*Function, env *Environment) *Class {
	return &Class{
		Name:    name,
		Methods: methods,
		Env:     env,
	}
}

func (cls *Class) findMethod(name string) (*Function, bool) {
	fn, ok := cls.Methods[name]
	return fn, ok
}

type ClassInstance struct {
	Cls    *Class
	Fields map[string]Object
}

func NewClassInstance(cls *Class) *ClassInstance {
	return &ClassInstance{
		Cls:    cls,
		Fields: make(map[string]Object),
	}
}

func (instance *ClassInstance) Get(name *tokens.Token) (Object, error) {
	v, ok := instance.Fields[name.Literal]
	if ok {
		return v, nil
	}

	method, methodOk := instance.Cls.findMethod(name.Literal)
	if methodOk {
		// We create a new environment nestled inside the method’s original closure.
		// When the method is called, that will become the parent of the method body’s environment.
		return method.bind(instance), nil
	}

	return nil, fmt.Errorf(ErrPropertyNotFound, name.Literal, instance.Inspect())
}

func (instance *ClassInstance) Set(name *tokens.Token, value Object) {
	instance.Fields[name.Literal] = value
}

func (instance *ClassInstance) Inspect() string {
	return instance.Cls.Name + ":instance"
}

func (instance *ClassInstance) Type() ObjectType {
	return OBJ_CLASS_INSTANCE
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.Block
	Env        *Environment
}

func (fn *Function) bind(instance *ClassInstance) *Function {
	newEnv := NewEnvWithOutter(fn.Env)
	newEnv.Set("this", instance)
	return &Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		Env:        newEnv,
	}
}

func (fn *Function) Inspect() string {
	sb := &strings.Builder{}
	sb.WriteString("fn(")

	var ps []string
	for _, p := range fn.Parameters {
		ps = append(ps, p.Name)
	}
	sb.WriteString(strings.Join(ps, ","))
	sb.WriteString(")")
	return sb.String()
}

func (fn *Function) Type() ObjectType {
	return OBJ_FUNCTION
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}
func (i *Integer) Type() ObjectType {
	return OBJ_INTEGER
}

type Bool struct {
	Value bool
}

func (b *Bool) Inspect() string {
	switch b.Value {
	case true:
		return "true"
	default:
		return "false"
	}
}

func (b *Bool) Type() ObjectType {
	return OBJ_BOOL
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return OBJ_STRING
}

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() ObjectType {
	return OBJ_NULL
}

type Return struct {
	Value Object
}

func (ret *Return) Inspect() string {
	return "return"
}
func (ret *Return) Type() ObjectType {
	return OBJ_RETURN
}

type Print struct {
}

func (ret *Print) Inspect() string {
	return "print"
}
func (ret *Print) Type() ObjectType {
	return OBJ_RETURN
}

type Error struct {
	Message string
}

func (err *Error) Inspect() string {
	return "Error: " + err.Message
}
func (err *Error) Type() ObjectType {
	return OBJ_ERROR
}
