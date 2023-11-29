package lox

import (
	"fmt"
	"strings"
)

type Env struct {
	nesting []string
	parent  *Env
	values  environment[Value]
	types   environment[Type]
}

type environment[T fmt.Stringer] map[string]T

func NewEnv(name string, parent *Env) *Env {
	env := &Env{
		parent: parent,
		values: make(environment[Value], 0),
		types:  make(environment[Type], 0),
	}
	if parent == nil {
		env.nesting = []string{name}
	} else {
		env.nesting = make([]string, len(parent.nesting)+1)
		copy(env.nesting, parent.nesting)
		env.nesting[len(parent.nesting)] = name
	}
	return env
}

func (e *Env) Name() string {
	return strings.Join(e.nesting, ":")
}

func (e *Env) String() string {
	values := make([]string, len(e.values))
	i := 0
	for key, val := range e.values {
		values[i] = fmt.Sprintf("%s=%s", key, val)
		i += 1
	}
	types := make([]string, len(e.types))
	i = 0
	for key, val := range e.types {
		types[i] = fmt.Sprintf("%s=%s", key, val)
		i += 1
	}
	return fmt.Sprintf("Env(%s){values:%s, types:%s}", e.Name(), strings.Join(values, ","), strings.Join(types, ", "))
}

func (e *Env) Value(name string) Value {
	if val, ok := e.values[name]; ok {
		return val
	}
	return nil
}

func (e *Env) ResolveValue(name string) (Value, *Env) {
	val := e.Value(name)
	if val == nil && e.parent != nil {
		return e.parent.ResolveValue(name)
	}
	return val, e
}

func (e *Env) SetValue(name string, val Value) (prev Value) {
	prev = e.values[name]
	e.values[name] = val
	return prev
}

func (e *Env) Type(name string) Type {
	if typ, ok := e.types[name]; ok {
		return typ
	}
	return TypeNone
}

func (e *Env) ResolveType(name string) (Type, *Env) {
	typ := e.Type(name)
	if typ == TypeNone && e.parent != nil {
		return e.parent.ResolveType(name)
	}
	return typ, e
}

func (e *Env) SetType(name string, typ Type) (prev Type) {
	prev = e.types[name]
	e.types[name] = typ
	return prev
}
