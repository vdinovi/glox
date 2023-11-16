package lox

import (
	"fmt"
)

type Environment[T any] struct {
	parent   *Environment[T]
	bindings map[string]T
	depth    int
}

func NewEnvironment[T any](parent *Environment[T]) *Environment[T] {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}
	return &Environment[T]{
		parent:   parent,
		bindings: make(map[string]T),
		depth:    depth,
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

func (env *Environment[T]) Set(key string, val T) error {
	// if _, ok := env.bindings[key]; ok {
	// 	return NewVariableRedeclarationError(key)
	// }
	env.bindings[key] = val
	return nil
}

func (env *Environment[T]) String() string {
	return fmt.Sprintf("(%d) %v", env.depth, env.bindings)
}

type Context struct {
	values *Environment[Value]
	types  *Environment[Type]
}

func NewContext() *Context {
	return &Context{
		values: NewEnvironment[Value](nil),
		types:  NewEnvironment[Type](nil),
	}
}

func (ctx *Context) PushEnvironment() {
	ctx.types = NewEnvironment(ctx.types)
	ctx.values = NewEnvironment(ctx.values)
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
