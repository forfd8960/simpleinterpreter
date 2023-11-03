package object

type Environment struct {
	kv     map[string]Object
	outter *Environment // the outter env for a function
}

func NewEnvWithOutter(outter *Environment) *Environment {
	return &Environment{
		kv:     map[string]Object{},
		outter: outter,
	}
}

func NewEnvironment() *Environment {
	return &Environment{
		kv:     make(map[string]Object, 10),
		outter: &Environment{},
	}
}

func (env *Environment) Get(identifier string) (Object, bool) {
	obj, ok := env.kv[identifier]
	if !ok && env.outter != nil {
		return env.outter.Get(identifier)
	}
	return obj, ok
}

func (env *Environment) Set(identifier string, value Object) Object {
	env.kv[identifier] = value
	return value
}
