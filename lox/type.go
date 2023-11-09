package lox

import "fmt"

//go:generate stringer -type Type  -trimprefix=Type
type Type int

const (
	TypeNil Type = iota
	TypeNumeric
	TypeString
	TypeBoolean
)

type Typed interface {
	Type() (Type, error)
}

type TypeError struct {
	Expression
	Err error
}

func (e TypeError) Error() string {
	return fmt.Sprintf("type error in %q: %s", e.Expression, e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}
