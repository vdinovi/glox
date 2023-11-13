package lox

import (
	"fmt"
)

type Value interface {
	fmt.Stringer
	Unwrap() any
	Truthy() bool
	Type() Type
}

type ValueString string

func (v ValueString) String() string {
	return string(v)
}

func (ValueString) Type() Type {
	return TypeString
}

func (v ValueString) Unwrap() any {
	return string(v)
}

func (v ValueString) Truthy() bool {
	return true
}

type ValueNumeric float64

func (v ValueNumeric) String() string {
	return fmt.Sprint(float64(v))
}

func (e ValueNumeric) Type() Type {
	return TypeNumeric
}

func (v ValueNumeric) Unwrap() any {
	return float64(v)
}

func (v ValueNumeric) Truthy() bool {
	return true
}

type ValueBoolean bool

var True = ValueBoolean(true)

var False = ValueBoolean(false)

func (v ValueBoolean) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (e ValueBoolean) Type() Type {
	return TypeBoolean
}

func (v ValueBoolean) Unwrap() any {
	return bool(v)
}

func (v ValueBoolean) Truthy() bool {
	return bool(v)
}

type ValueNil struct{}

var Nil = ValueNil(struct{}{})

func (v ValueNil) String() string {
	return "nil"
}

func (e ValueNil) Type() Type {
	return TypeNil
}

func (v ValueNil) Unwrap() any {
	return struct{}{}
}

func (v ValueNil) Truthy() bool {
	return false
}
