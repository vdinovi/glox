package lox

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Environment[T fmt.Stringer] struct {
	parent   *Environment[T]
	bindings map[string]T
	depth    int
}

func NewEnvironment[T fmt.Stringer](parent *Environment[T]) *Environment[T] {
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
	var sb strings.Builder
	fmt.Fprintf(&sb, "Environment(%d){", env.depth)
	for key, val := range env.bindings {
		fmt.Fprintf(&sb, " %s=%s,", key, val)
	}
	sb.WriteString(" }")
	return sb.String()
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

func (ctx *Context) Enter(phase string) (exit func()) {
	ctx.PushEnvironment()
	log.Debug().Msgf("(%s) enter %s", phase, ctx.values.String())
	return func() {
		ctx.PopEnvironment()
		log.Debug().Msgf("(%s) enter %s", phase, ctx.values.String())
	}
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
