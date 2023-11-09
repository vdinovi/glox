package lox

import "fmt"

type Value interface {
	Typed
	fmt.Stringer
	Unwrap() any
	Truthy() bool
}

type StringValue string

func (v StringValue) String() string {
	return string(v)
}

func (StringValue) Type() (Type, error) {
	return TypeString, nil
}

func (v StringValue) Unwrap() any {
	return string(v)
}

func (v StringValue) Truthy() bool {
	return true
}

type NumericValue float64

func (v NumericValue) String() string {
	return fmt.Sprint(float64(v))
}
func (e NumericValue) Type() (Type, error) {
	return TypeNumeric, nil
}

func (v NumericValue) Unwrap() any {
	return float64(v)
}

func (v NumericValue) Truthy() bool {
	return true
}

type BooleanValue bool

func (v BooleanValue) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (e BooleanValue) Type() (Type, error) {
	return TypeBoolean, nil
}

func (v BooleanValue) Unwrap() any {
	return bool(v)
}

func (v BooleanValue) Truthy() bool {
	return bool(v)
}

type NilValue struct{}

func (v NilValue) String() string {
	return "nil"
}

func (e NilValue) Type() (Type, error) {
	return TypeNil, nil
}

func (v NilValue) Unwrap() any {
	return nil
}

func (v NilValue) Truthy() bool {
	return false
}
