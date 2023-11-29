package lox

import (
	"fmt"
	"strings"
)

type Environment[T fmt.Stringer] struct {
	name     string
	parent   *Environment[T]
	bindings map[string]T
}

func NewEnvironment[T fmt.Stringer](name string, parent *Environment[T]) *Environment[T] {
	if parent != nil {
		name = fmt.Sprintf("%s:%s", parent.name, name)
	}
	return &Environment[T]{
		name:     name,
		parent:   parent,
		bindings: make(map[string]T),
	}
}

func (env *Environment[T]) Lookup(name string) (*T, *Environment[T]) {
	val, ok := env.bindings[name]
	if ok {
		return &val, env
	}
	if env.parent != nil {
		return env.parent.Lookup(name)
	}
	return nil, nil
}

func (env *Environment[T]) Get(key string, def T) T {
	t, ok := env.bindings[key]
	if !ok {
		return def
	}
	return t
}

func (env *Environment[T]) Set(key string, val T) (prev *T, err error) {
	// if _, ok := env.bindings[key]; ok {
	// 	return NewVariableRedeclarationError(key)
	// }
	p, ok := env.bindings[key]
	if ok {
		prev = &p
	}
	env.bindings[key] = val
	return prev, nil
}

func (env *Environment[T]) String() string {
	bindings := make([]string, len(env.bindings))
	i := 0
	for key, val := range env.bindings {
		bindings[i] = fmt.Sprintf("%s=%s", key, val)
		i += 1
	}
	return fmt.Sprintf("Env(%s){%s}", env.name, strings.Join(bindings, ","))
}

type Context struct {
	values *Environment[Value]
	types  *Environment[Type]
}

func NewContext() *Context {
	return &Context{
		values: NewEnvironment[Value]("root", nil),
		types:  NewEnvironment[Type]("root", nil),
	}
}

func (ctx *Context) Copy() *Context {
	return &Context{
		values: ctx.values,
		types:  ctx.types,
	}
}

func (ctx *Context) PushEnvironment(name string) {
	ctx.types = NewEnvironment(name, ctx.types)
	ctx.values = NewEnvironment(name, ctx.values)
}

func (ctx *Context) PopEnvironment() {
	ctx.types = ctx.types.parent
	ctx.values = ctx.values.parent
}

type VariableRedeclarationError struct {
	Name string
}

func (e VariableRedeclarationError) Error() string {
	return fmt.Sprintf("variable %s already declared", e.Name)
}

func NewVariableRedeclarationError(name string) VariableRedeclarationError {
	return VariableRedeclarationError{Name: name}
}
