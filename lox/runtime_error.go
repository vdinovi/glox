package lox

import "fmt"

// Container for all runtime errors
type RuntimeError struct {
	Err      error // the wrapped error
	Position       // the originating location
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error on line %d: %s", e.Position.Line, e.Err)
}

func (e RuntimeError) Unwrap() error {
	return e.Err
}

func NewRuntimeError(err error, pos Position) RuntimeError {
	return RuntimeError{
		Err:      err,
		Position: pos,
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
