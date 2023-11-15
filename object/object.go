package object

import (
	"strconv"
	"strings"

	"github.com/forfd8960/simpleinterpreter/ast"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Class struct {
	Name    string
	Methods []*Function
	Env     *Environment
}

func (cls *Class) Inspect() string {
	return cls.Name
}

func (cls *Class) Type() ObjectType {
	return OBJ_CLASS
}

type ClassInstance struct {
	Cls *Class
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
