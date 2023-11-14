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
