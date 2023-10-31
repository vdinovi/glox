package lox

import (
	"fmt"
)

type LexerError struct {
	Err    error
	Line   int
	Column int
}

func (e LexerError) Error() string {
	return fmt.Sprintf("LexerError at (%d,%d): %s", e.Line, e.Column, e.Err.Error())
}

func (e LexerError) Unwrap() error {
	return e.Err
}

type UnterminatedStringError struct{}

func (e UnterminatedStringError) Error() string {
	return "unterminated string"
}

type InvaldIdentifierError struct {
	ID string
}

func (e InvaldIdentifierError) Error() string {
	return fmt.Sprintf("invalid identifier %q", e.ID)
}
