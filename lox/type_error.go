package lox

import "fmt"

// Container for all type-related errors
type TypeError struct {
	Err    error // the wrapped error
	Line   int   // the originating line in which the error occurred
	Column int   // the originating column in which the error occurred
}

func (e TypeError) Error() string {
	return fmt.Sprintf("Type Error on line %d: %s", e.Line, e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}

func NewTypeError(err error, line, column int) TypeError {
	return TypeError{
		Err:    err,
		Line:   line,
		Column: column,
	}
}

// Error indicating that the types do not match
type TypeMismatchError struct {
	Types []Type
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("types do not match %v", e.Types)
}

func NewTypeMismatchError(types ...Type) TypeMismatchError {
	return TypeMismatchError{
		Types: types,
	}
}

// Error indicating that the operation cannot be applied to the types
type InvalidOperatorForTypeError struct {
	OperatorType
	Types []Type
}

func (e InvalidOperatorForTypeError) Error() string {
	return fmt.Sprintf("operator %s can't be applied to types %v", e.OperatorType, e.Types)
}

func NewInvalidOperatorForTypeError(opType OperatorType, types ...Type) InvalidOperatorForTypeError {
	return InvalidOperatorForTypeError{
		OperatorType: opType,
		Types:        types,
	}
}
