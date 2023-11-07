package lox

import (
	"fmt"
	"io"
	"unicode"
)

type Lexer struct {
	scan  Scanner
	fname string
}

func NewLexer(rd io.RuneReader) *Lexer {
	return &Lexer{
		scan: NewScanner(rd),
	}
}

func (l *Lexer) SetFilename(fname string) {
	l.fname = fname
}

func (l *Lexer) ScanTokens() ([]Token, error) {
	tokens := []Token{}
	for {
		token, err := l.next()
		if err == nil {
			tokens = append(tokens, *token)
		} else if err == io.EOF {
			tokens = append(tokens, EofToken)
			break
		} else {
			return nil, err
		}
	}
	return tokens, nil
}

var isNewline = IsChar('\n')
var isQuote = IsChar('"')
var isDot = IsChar('.')
var isEquals = IsChar('=')
var isSlash = IsChar('/')

var isNotDigit = func(ch rune) bool {
	return !unicode.IsDigit(ch)
}
var isNotWhitespace = func(ch rune) bool {
	return !unicode.IsSpace(ch)
}
var isNotLetterOrUnderscore = func(ch rune) bool {
	return ch != '_' && !unicode.IsLetter(ch)
}

func (l *Lexer) next() (*Token, error) {
	line, column := l.scan.Position()
	token := Token{
		Line:   line,
		Column: column,
	}

	_, err := MatchUntil(l.scan, isNotWhitespace)
	if err != nil {
		return nil, err
	}

	ch, err := l.scan.Next()
	if err != nil {
		return nil, err
	}

	err = nil
	switch ch {
	case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
		token.Lexem = string(ch)
		token.Type = TokenTypeFor(token.Lexem)
	case '!', '=', '<', '>':
		err = l.maybeEquals(&token, string(ch), line, column)
	case '/':
		err = l.commentOrSlash(&token, string(ch), line, column)
	case '"':
		err = l.string(&token, string(ch), line, column)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		err = l.number(&token, string(ch), line, column)
	default:
		if isNotLetterOrUnderscore(ch) {
			err = &LexerError{
				Err:    &UnexpectedCharacterError{ch},
				Line:   line,
				Column: column,
			}
		} else {
			err = l.identifier(&token, string(ch), line, column)
		}
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (l *Lexer) maybeEquals(token *Token, lexem string, line, column int) error {
	nextCh, err := MatchRune(l.scan, isEquals)
	if err != nil && err != io.EOF {
		return err
	}
	if err == io.EOF || nextCh != '=' {
		token.Lexem = lexem
	} else {
		token.Lexem = lexem + "="
	}
	token.Type = TokenTypeFor(token.Lexem)
	return nil
}

func (l *Lexer) commentOrSlash(token *Token, lexem string, line, column int) error {
	nextCh, err := MatchRune(l.scan, isSlash)
	if err != nil && err != io.EOF {
		return err
	}
	if err == io.EOF || nextCh != '/' {
		token.Lexem = lexem
		token.Type = TokenSlash
	} else {
		chars, err := MatchUntil(l.scan, isNewline)
		if err != nil && err != io.EOF {
			return err
		}
		token.Lexem = string(chars)
		token.Type = TokenComment
	}
	return nil
}

func (l *Lexer) string(token *Token, lexem string, line, column int) error {
	chars, err := MatchThrough(l.scan, isQuote)
	if err != nil && err != io.EOF {
		return err
	}
	if err == io.EOF {
		return &LexerError{
			Err:    &UnterminatedStringError{},
			Line:   line,
			Column: column,
		}
	} else {
		token.Lexem = string(chars[:len(chars)-1])
		token.Type = TokenString
	}
	return nil
}

func (l *Lexer) number(token *Token, lexem string, line, column int) error {
	chars, err := MatchUntil(l.scan, isNotDigit)
	if err != nil && err != io.EOF {
		return err
	}
	token.Lexem = lexem + string(chars)
	token.Type = TokenNumber
	if err != io.EOF {
		nextCh, err := MatchRune(l.scan, isDot)
		if err != nil && err != io.EOF {
			return err
		}
		if nextCh == '.' {
			chars, err := MatchUntil(l.scan, isNotDigit)
			if err != nil && err != io.EOF {
				return err
			}
			token.Lexem += "." + string(chars)
			token.Type = TokenNumber
		}
	}
	return nil
}

func (l *Lexer) identifier(token *Token, lexem string, line, column int) error {
	chars, err := MatchUntil(l.scan, isNotLetterOrUnderscore)
	if err != nil && err != io.EOF {
		return err
	}
	token.Lexem = lexem + string(chars)
	if tt := TokenTypeFor(token.Lexem); tt != TokenNone {
		token.Type = tt
	} else {
		for _, ch := range token.Lexem {
			if ch != '_' && !unicode.IsLetter(ch) {
				return &LexerError{
					Err:    &InvaldIdentifierError{token.Lexem},
					Line:   line,
					Column: column,
				}
			}
		}
		token.Type = TokenIdentifier
	}
	return nil
}

type LexerError struct {
	Err    error
	Line   int
	Column int
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("LexerError at (%d,%d): %s", e.Line, e.Column, e.Err.Error())
}

func (e *LexerError) Unwrap() error {
	return e.Err
}

type UnterminatedStringError struct{}

func (e *UnterminatedStringError) Error() string {
	return "unterminated string"
}

type InvaldIdentifierError struct {
	ID string
}

func (e *InvaldIdentifierError) Error() string {
	return fmt.Sprintf("invalid identifier %q", e.ID)
}

type UnexpectedCharacterError struct {
	Char rune
}

func (e *UnexpectedCharacterError) Error() string {
	return fmt.Sprintf("unexpected character %q", e.Char)
}

type MatchFunc func(rune) bool

func IsChar(delim rune) func(rune) bool {
	return func(ch rune) bool {
		return ch == delim
	}
}

func HasChar(delims ...rune) func(rune) bool {
	set := make(map[rune]struct{})
	for _, d := range delims {
		set[d] = struct{}{}
	}
	return func(ch rune) bool {
		_, ok := set[ch]
		return ok
	}
}

func NotChar(delims ...rune) func(rune) bool {
	set := make(map[rune]struct{})
	for _, d := range delims {
		set[d] = struct{}{}
	}
	return func(ch rune) bool {
		_, ok := set[ch]
		return !ok
	}
}

type Scanner interface {
	Next() (rune, error)
	Peek(int) ([]rune, error)
	Position() (int, int)
	Advance(int) error
}

func NewScanner(rd io.RuneReader) Scanner {
	sc := stringScanner{}
	for {
		ch, _, err := rd.ReadRune()
		if err != nil {
			break
		}
		sc.chars = append(sc.chars, ch)
	}
	return &sc
}

func MatchRune(s Scanner, mf MatchFunc) (rune, error) {
	chars, err := s.Peek(1)
	if err != nil {
		return -1, err
	}
	ch := chars[0]
	if mf(ch) {
		if err := s.Advance(1); err != nil {
			return ch, err
		}
	}
	return ch, nil
}

func MatchUntil(s Scanner, mf MatchFunc) ([]rune, error) {
	chars := []rune{}
	for {
		if next, err := s.Peek(1); err != nil {
			return chars, err
		} else if mf(next[0]) {
			return chars, nil
		} else {
			if ch, err := s.Next(); err != nil {
				return chars, err
			} else {
				chars = append(chars, ch)
			}
		}
	}
}

func MatchThrough(s Scanner, mf MatchFunc) ([]rune, error) {
	chars, err := MatchUntil(s, mf)
	if err != nil {
		return chars, err
	}
	next, err := s.Next()
	if err != nil {
		return chars, err
	}
	return append(chars, next), nil
}

type stringScanner struct {
	chars  []rune
	offset int
	line   int
	column int
}

func (s *stringScanner) Next() (rune, error) {
	if s.offset >= len(s.chars) {
		return -1, io.EOF
	}
	ch := s.chars[s.offset]
	if err := s.Advance(1); err != nil {
		return ch, err
	}
	return ch, nil
}

func (s *stringScanner) Peek(size int) ([]rune, error) {
	from, to := s.offset, s.offset+size
	if to > len(s.chars) {
		return nil, io.EOF
	}
	return s.chars[from:to], nil
}

func (s *stringScanner) Position() (int, int) {
	return s.line + 1, s.column + 1
}

func (s *stringScanner) Advance(size int) error {
	from, to := s.offset, s.offset+size
	if to > len(s.chars) {
		return io.EOF
	}
	for _, ch := range s.chars[from:to] {
		s.column += 1
		switch ch {
		case '\n':
			s.line += 1
			s.column = 0
		}
	}
	s.offset += size
	return nil
}
