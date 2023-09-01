package object

import (
	"strconv"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
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

type Error struct {
	Message string
}

func (err *Error) Inspect() string {
	return "Error: " + err.Message
}
func (err *Error) Type() ObjectType {
	return OBJ_ERROR
}
