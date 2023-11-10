package lox

import (
	"fmt"
)

type Value interface {
	Typed
	fmt.Stringer
	Unwrap() any
	Truthy() bool
}

type ValueString string

func (v ValueString) String() string {
	return fmt.Sprintf("Value(%q)", string(v))
}

func (ValueString) Type(Symbols) (Type, error) {
	return TypeString, nil
}

func (v ValueString) Unwrap() any {
	return string(v)
}

func (v ValueString) Truthy() bool {
	return true
}

type ValueNumeric float64

func (v ValueNumeric) String() string {
	return fmt.Sprintf("Value(%.3f)", float64(v))
}

func (e ValueNumeric) Type(Symbols) (Type, error) {
	return TypeNumeric, nil
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
		return "Value(true)"
	}
	return "Value(false)"
}

func (e ValueBoolean) Type(Symbols) (Type, error) {
	return TypeBoolean, nil
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
	return "Value(nil)"
}

func (e ValueNil) Type(Symbols) (Type, error) {
	return TypeNil, nil
}

func (v ValueNil) Unwrap() any {
	return struct{}{}
}

func (v ValueNil) Truthy() bool {
	return false
}
