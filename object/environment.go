package object

type Environment struct {
	env map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{env: make(map[string]Object, 10)}
}

func (env *Environment) Get(identifier string) (Object, bool) {
	obj, ok := env.env[identifier]
	return obj, ok
}

func (env *Environment) Set(identifier string, value Object) Object {
	env.env[identifier] = value
	return value
}
