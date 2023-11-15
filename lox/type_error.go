package lox

import (
	"fmt"
)

// Container for all type-related errors
type TypeError struct {
	Err      error // the wrapped error
	Position       // the originating location
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type Error on line %d: %s", e.Position.Line, e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}

func NewTypeError(err error, pos Position) TypeError {
	return TypeError{
		Err:      err,
		Position: pos,
	}
}

// Error indicating that the types do not match
type TypeMismatchError struct {
	Left  Type
	Right Type
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("types %s does not match %s", e.Left, e.Right)
}

func NewTypeMismatchError(left, right Type) TypeMismatchError {
	return TypeMismatchError{
		Left:  left,
		Right: right,
	}
}

// Error indicating that the operation cannot be applied to the types
type InvalidUnaryOperatorForTypeError struct {
	OperatorType
	Right Type
}

func (e InvalidUnaryOperatorForTypeError) Error() string {
	return fmt.Sprintf("unary operator %s can't be applied to type %v", e.OperatorType, e.Right)
}

func NewInvalidUnaryOperatorForTypeError(opType OperatorType, right Type) InvalidUnaryOperatorForTypeError {
	return InvalidUnaryOperatorForTypeError{
		OperatorType: opType,
		Right:        right,
	}
}

// Error indicating that the binary operation cannot be applied to the types
type InvalidBinaryOperatorForTypeError struct {
	OperatorType
	Left  Type
	Right Type
}

func (e InvalidBinaryOperatorForTypeError) Error() string {
	return fmt.Sprintf("binary operator %s can't be applied to types %v and %v", e.OperatorType, e.Left, e.Right)
}

func NewInvalidBinaryOperatorForTypeError(opType OperatorType, left, right Type) InvalidBinaryOperatorForTypeError {
	return InvalidBinaryOperatorForTypeError{
		OperatorType: opType,
		Left:         left,
		Right:        right,
	}
}
