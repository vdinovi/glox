package lox

import "fmt"

// Container for all runtime errors
type RuntimeError struct {
	Err    error // the wrapped error
	Line   int   // the originating line in which the error occurred
	Column int   // the originating column in which the error occurred
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error on line %d: %s", e.Line, e.Err)
}

func (e RuntimeError) Unwrap() error {
	return e.Err
}

func NewRuntimeError(err error, line, column int) RuntimeError {
	return RuntimeError{
		Err:    err,
		Line:   line,
		Column: column,
	}
}

// Error indicating that the value cannot be downcasted to the specified type
type DowncastError struct {
	Value
	Type string
}

func (e DowncastError) Error() string {
	return fmt.Sprintf("failed to downcast value %s to type %s", e.Value, e.Type)
}

func NewDowncastError(val Value, typ string) DowncastError {
	return DowncastError{
		Value: val,
		Type:  typ,
	}
}

// Error indicating that the expression cannot be applied to the types
type InvalidExpressionForTypeError struct {
	Expression
	Types []Type
}

func (e InvalidExpressionForTypeError) Error() string {
	return fmt.Sprintf("expression %s can't be applied to types %v", e.Expression, e.Types)
}

func NewInvalidExpressionForTypeError(expr Expression, types ...Type) InvalidExpressionForTypeError {
	return InvalidExpressionForTypeError{
		Expression: expr,
		Types:      types,
	}
}
