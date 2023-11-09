package lox

import (
	"fmt"
	"io"
	"unicode"

	"github.com/rs/zerolog/log"
)

type Lexer struct {
	scan runeScanner
}

func NewLexer(rd io.RuneReader) (*Lexer, error) {
	scanner := runeScanner{}
	if err := scanner.fill(rd); err != nil {
		return nil, err
	}
	return &Lexer{scan: scanner}, nil
}

func (l *Lexer) Scan() ([]Token, error) {
	log.Debug().Msgf("(lexer) scanning %d runes", len(l.scan.runes))
	tokens := []Token{}
	for {
		token, err := l.next()
		if err == nil {
			log.Debug().Msgf("(lexer) token: %s", token)
			tokens = append(tokens, token)
		} else if err == io.EOF {
			log.Debug().Msgf("(lexer) reached EOF")
			tokens = append(tokens, EofToken)
			break
		} else {
			log.Error().Msgf("(lexer) error: %s", err)
			return nil, err
		}
	}
	log.Debug().Msgf("(lexer) produced %d tokens", len(tokens))
	return tokens, nil
}

type runeMatchFunc func(rune) bool

func isRune(want rune) func(rune) bool {
	return func(r rune) bool {
		return r == want
	}
}

var isNewline = isRune('\n')
var isQuote = isRune('"')
var isDot = isRune('.')
var isEquals = isRune('=')
var isSlash = isRune('/')

var isNotDigit = func(r rune) bool {
	return !unicode.IsDigit(r)
}

var isNotWhitespace = func(r rune) bool {
	return !unicode.IsSpace(r)
}

var isNotLetterOrUnderscore = func(r rune) bool {
	return r != '_' && !unicode.IsLetter(r)
}

func (l *Lexer) next() (Token, error) {
	line, column := l.scan.position()
	token := Token{Line: line, Column: column}

	if _, err := l.scan.until(isNotWhitespace); err != nil {
		return NoneToken, err
	}

	next, err := l.scan.advance()
	if err != nil {
		return NoneToken, err
	}

	switch next {
	case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
		token.Lexem = string(next)
	case '!', '=', '<', '>':
		eq, ok, err := l.scan.match(isEquals)
		if err != nil && err != io.EOF {
			return NoneToken, err
		}
		if !ok || err == io.EOF {
			token.Lexem = string(next)
		} else {
			token.Lexem = string(next) + string(eq)
		}
	case '/':
		_, ok, err := l.scan.match(isSlash)
		if err != nil && err != io.EOF {
			return NoneToken, err
		}
		if !ok || err == io.EOF {
			token.Type = TokenSlash
			token.Lexem = string(next)
		} else {
			runes, err := l.scan.until(isNewline)
			if err != nil && err != io.EOF {
				return NoneToken, err
			}
			token.Type = TokenComment
			token.Lexem = string(runes)
		}
	case '"':
		runes, err := l.scan.through(isQuote)
		if err != nil && err != io.EOF {
			return NoneToken, err
		}
		if err == io.EOF {
			return NoneToken, &LexError{
				Err:    &UnterminatedStringError{},
				Line:   line,
				Column: column,
			}
		}
		token.Type = TokenString
		token.Lexem = string(runes)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		integral, err := l.scan.until(isNotDigit)
		if err != nil && err != io.EOF {
			return NoneToken, err
		}
		if err == io.EOF {
			token.Type = TokenNumber
			token.Lexem = string(next) + string(integral)
		} else {
			dec, ok, err := l.scan.match(isDot)
			if err != nil && err != io.EOF {
				return NoneToken, err
			}
			if !ok || err == io.EOF {
				token.Type = TokenNumber
				token.Lexem = string(next) + string(integral)
			} else {
				fractional, err := l.scan.until(isNotDigit)
				if err != nil && err != io.EOF {
					return NoneToken, err
				}
				token.Type = TokenNumber
				token.Lexem = string(next) + string(integral) + string(dec) + string(fractional)
			}
		}
	default:
		if isNotLetterOrUnderscore(next) {
			return NoneToken, &LexError{
				Err:    &UnexpectedCharacterError{next},
				Line:   line,
				Column: column,
			}
		} else {
			runes, err := l.scan.until(isNotLetterOrUnderscore)
			if err != nil && err != io.EOF {
				return NoneToken, err
			}
			lexem := string(next) + string(runes)
			if t := TokenTypeFor(lexem); t != TokenNone {
				token.Type = t
			} else {
				token.Type = TokenIdentifier
			}
			token.Lexem = lexem
		}
	}
	if token.Type == TokenNone {
		token.Type = TokenTypeFor(token.Lexem)
	}
	return token, nil
}

// TODO: replace with a bufio-based scanner
type runeScanner struct {
	runes  []rune
	offset int
	line   int
	column int
}

func (s *runeScanner) fill(rd io.RuneReader) error {
	for {
		r, _, err := rd.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		s.runes = append(s.runes, r)
	}
}

func (s *runeScanner) done() bool {
	return s.offset == len(s.runes)
}

func (s *runeScanner) peek() (rune, error) {
	if s.done() {
		return -1, io.EOF
	}
	return s.runes[s.offset], nil
}

func (s *runeScanner) advance() (rune, error) {
	if s.done() {
		return -1, io.EOF
	}
	rune := s.runes[s.offset]
	s.offset += 1
	return rune, nil
}

func (s *runeScanner) match(fn runeMatchFunc) (rune, bool, error) {
	r, err := s.peek()
	if err != nil {
		return -1, false, err
	}
	if fn(r) {
		if _, err = s.advance(); err != nil {
			return -1, false, err
		}
		return r, true, nil
	}
	return r, false, nil
}

func (s *runeScanner) through(fn runeMatchFunc) ([]rune, error) {
	runes := []rune{}
	for {
		r, err := s.peek()
		if err != nil {
			return runes, err
		}
		if fn(r) {
			if _, err = s.advance(); err != nil {
				return runes, err
			}
			return runes, nil
		}
		if _, err = s.advance(); err != nil {
			return runes, err
		}
		runes = append(runes, r)
	}
}

func (s *runeScanner) until(fn runeMatchFunc) ([]rune, error) {
	runes := []rune{}
	for {
		r, err := s.peek()
		if err != nil {
			return runes, err
		}
		if fn(r) {
			return runes, nil
		}
		if _, err = s.advance(); err != nil {
			return runes, err
		}
		runes = append(runes, r)
	}
}

func (s *runeScanner) position() (int, int) {
	return s.line + 1, s.column + 1
}

type LexError struct {
	Err    error
	Line   int
	Column int
}

func (e *LexError) Error() string {
	return fmt.Sprintf("LexError at (%d, %d): %s", e.Line, e.Column, e.Err)
}

func (e *LexError) Unwrap() error {
	return e.Err
}

type UnterminatedStringError struct{}

func (e *UnterminatedStringError) Error() string {
	return "unterminated string"
}

type UnexpectedCharacterError struct {
	Char rune
}

func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("unexpected character %q", e.Char)
}
