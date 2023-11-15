package lox

import "fmt"

// Container for all syntax-related errors
type SyntaxError struct {
	Err      error // the wrapped error
	Position       // the originating location
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("Syntax Error on line %d: %s", e.Position.Line, e.Err)
}

func (e SyntaxError) Unwrap() error {
	return e.Err
}

func NewSyntaxError(err error, pos Position) SyntaxError {
	return SyntaxError{
		Err:      err,
		Position: pos,
	}
}

// Error indicating that an unexpected character (rune) was encountered
type UnexpectedCharacterError struct {
	Expected string
	Actual   rune
}

func (e UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("expected %s but got character %q", e.Expected, e.Actual)
}

func NewUnexpectedCharacterError(expected string, actual rune) UnexpectedCharacterError {
	return UnexpectedCharacterError{
		Expected: expected,
		Actual:   actual,
	}
}

type UnexpectedTokenError struct {
	Expected string
	Actual   Token
}

// Error indicating that an unexpected token was encountered
func (e UnexpectedTokenError) Error() string {
	return fmt.Sprintf("expected %s but got token %q", e.Expected, e.Actual)
}

func NewUnexpectedTokenError(expected string, actual Token) UnexpectedTokenError {
	return UnexpectedTokenError{
		Expected: expected,
		Actual:   actual,
	}
}

// Error indicating that a string went unterminated
type UnterminatedStringError struct{}

func (e UnterminatedStringError) Error() string {
	return "unterminated string"
}

func NewUnterminatedStringError() UnterminatedStringError {
	return UnterminatedStringError{}
}

// Error indicating that a token went unmatched
type UnmatchedTokenError struct {
	Token
}

func (e UnmatchedTokenError) Error() string {
	return fmt.Sprintf("unmatched token %s", e.Token)
}

func NewUnmatchedTokenError(token Token) UnmatchedTokenError {
	return UnmatchedTokenError{
		Token: token,
	}
}

// Error indicating that an expected terminal was not found
type MissingTerminalError struct{}

func (e MissingTerminalError) Error() string {
	return "missing terminal"
}

func NewMissingTerminalError(token Token) MissingTerminalError {
	return MissingTerminalError{}
}

// Error indicating an error in numeric conversion from token
type NumberConversionError struct {
	Err error
	Token
}

func (e NumberConversionError) Error() string {
	return fmt.Sprintf("failed to convert %s to number: %s", e.Token, e.Err)
}

func (e NumberConversionError) Unwrap() error {
	return e.Err
}

func NewNumberConversionError(err error, token Token) NumberConversionError {
	return NumberConversionError{
		Err:   err,
		Token: token,
	}
}

// Error indicating an invalid target in assignment expression
type InvalidAssignmentTargetError struct {
	Name string
}

func (e InvalidAssignmentTargetError) Error() string {
	return fmt.Sprintf("invalid assignment target %s", e.Name)
}

func NewInvalidAssignmentTargetError(name string) InvalidAssignmentTargetError {
	return InvalidAssignmentTargetError{Name: name}
}

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

// Error indicating that the variable is undefined
type UndefinedVariableError struct {
	Name string
}

func (e UndefinedVariableError) Error() string {
	return fmt.Sprintf("variable %s is not defined", e.Name)
}

func NewUndefinedVariableError(name string) UndefinedVariableError {
	return UndefinedVariableError{Name: name}
}

// Error indicating division by zero
type DivideByZeroError struct {
	Numerator   ValueNumeric
	Denominator ValueNumeric
}

func (e DivideByZeroError) Error() string {
	return fmt.Sprintf("Divide by zero (%s / %s)", e.Numerator, e.Denominator)
}

func NewDivideByZeroError(num, denom ValueNumeric) DivideByZeroError {
	return DivideByZeroError{Numerator: num, Denominator: denom}
}

type ValueError struct {
	Message string
}

func (v ValueError) Error() string {
	return fmt.Sprintf("Value Error: %s", v.Message)
}

func NewValueError(message string) ValueError {
	return ValueError{Message: message}
}
