package lox

import (
	"fmt"
)

//go:generate stringer -type Type  -trimprefix=Type
type Type int

const (
	ErrType Type = iota
	TypeNil
	TypeNumeric
	TypeString
	TypeBoolean
)

type Typed interface {
	Type() (Type, error)
}

type TypeError struct {
	Err error
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type Error: %s", e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}

type TypeMismatchError struct {
	Types []Type
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("types do not match %v", e.Types)
}

func TypeMismatch(types ...Type) TypeMismatchError {
	return TypeMismatchError{
		Types: types,
	}
}
