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
