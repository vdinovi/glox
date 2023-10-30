package lox

import (
	"bufio"
	"fmt"
	"io"
)

func Lex(r io.Reader) ([]Token, []LexError) {
	scanner := bufio.NewScanner(r)
	var toks []Token
	var errs []LexError

	for line := 1; scanner.Scan(); line += 1 {
		if err := scanner.Err(); err != nil {
			errs = append(errs, LexError{Line: line, Err: err})
		} else {
			toks = append(toks, Token{
				Type:  TokenLine,
				Lexem: scanner.Text(),
				Line:  line,
			})
		}
	}
	return toks, errs
}

type LexError struct {
	Line int
	Err  error
}

func (e LexError) Error() string {
	return fmt.Sprintf("[LexError on line %d] %s", e.Line, e.Err)
}

func (e LexError) Unwrap() error {
	return e.Err
}
