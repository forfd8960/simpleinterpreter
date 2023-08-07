package object

import "strconv"

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integrer struct {
	Value int64
}

func (i *Integrer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}
func (i *Integrer) Type() ObjectType {
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

type Null struct{}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() ObjectType {
	return OBJ_NULL
}
