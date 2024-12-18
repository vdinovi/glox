package lox

import (
	"bufio"
	"io"
	"unicode"

	"github.com/rs/zerolog/log"
)

func Scan(ctx *Context, rd io.Reader) ([]Token, error) {
	restore := ctx.StartPhase(PhaseLex)
	defer restore()
	lexer, err := NewLexer(ctx, bufio.NewReader(rd))
	if err != nil {
		return nil, err
	}
	return lexer.Scan()
}

type Lexer struct {
	ctx  *Context
	scan runeScanner
}

func NewLexer(ctx *Context, rd io.RuneReader) (*Lexer, error) {
	l := &Lexer{
		ctx:  ctx,
		scan: runeScanner{},
	}
	if err := l.scan.fill(rd); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Lexer) Scan() ([]Token, error) {
	log.Debug().Msgf("(%s) scanning %d runes", l.ctx.Phase(), len(l.scan.runes))
	tokens := []Token{}
	for {
		token, err := l.next()
		if err == nil {
			log.Debug().Msgf("(%s) token: %s", l.ctx.Phase(), token)
			tokens = append(tokens, *token)
		} else if err == io.EOF {
			tokens = append(tokens, eofToken)
			log.Debug().Msgf("(%s) token: %s", l.ctx.Phase(), eofToken)
			log.Debug().Msgf("(%s) done", l.ctx.Phase())
			break
		} else {
			log.Error().Msgf("(%s) error: %s", l.ctx.Phase(), err)
			return nil, err
		}
	}
	log.Debug().Msgf("(%s) produced %d tokens", l.ctx.Phase(), len(tokens))
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

func (l *Lexer) next() (*Token, error) {
	if _, err := l.scan.until(isNotWhitespace); err != nil {
		return nil, err
	}

	token := Token{Position: l.scan.position()}
	next, err := l.scan.advance()
	if err != nil {
		return nil, err
	}

	switch next {
	case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
		token.Lexem = string(next)
	case '!', '=', '<', '>':
		eq, ok, err := l.scan.match(isEquals)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if !ok || err == io.EOF {
			token.Lexem = string(next)
		} else {
			token.Lexem = string(next) + string(eq)
		}
	case '/':
		_, ok, err := l.scan.match(isSlash)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if !ok || err == io.EOF {
			token.Type = TokenSlash
			token.Lexem = string(next)
		} else {
			runes, err := l.scan.until(isNewline)
			if err != nil && err != io.EOF {
				return nil, err
			}
			token.Type = TokenComment
			token.Lexem = string(runes)
		}
	case '"':
		runes, err := l.scan.through(isQuote)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			return nil, NewSyntaxError(NewUnterminatedStringError(), token.Position)
		}
		token.Type = TokenString
		token.Lexem = string(runes)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		integral, err := l.scan.until(isNotDigit)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			token.Type = TokenNumber
			token.Lexem = string(next) + string(integral)
		} else {
			dec, ok, err := l.scan.match(isDot)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if !ok || err == io.EOF {
				token.Type = TokenNumber
				token.Lexem = string(next) + string(integral)
			} else {
				fractional, err := l.scan.until(isNotDigit)
				if err != nil && err != io.EOF {
					return nil, err
				}
				token.Type = TokenNumber
				token.Lexem = string(next) + string(integral) + string(dec) + string(fractional)
			}
		}
	default:
		if isNotLetterOrUnderscore(next) {
			return nil, NewSyntaxError(
				NewUnexpectedCharacterError("a letter or underscore character", next), token.Position,
			)
		} else {
			runes, err := l.scan.until(isNotLetterOrUnderscore)
			if err != nil && err != io.EOF {
				return nil, err
			}
			lexem := string(next) + string(runes)
			if t := tokenTypeFor(lexem); t != ErrToken {
				token.Type = t
			} else {
				token.Type = TokenIdentifier
			}
			token.Lexem = lexem
		}
	}
	if token.Type == ErrToken {
		token.Type = tokenTypeFor(token.Lexem)
	}
	return &token, nil
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
	if rune == '\n' {
		s.line += 1
		s.column = 0
	} else {
		s.column += 1
	}

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

func (s *runeScanner) position() Position {
	return Position{s.line + 1, s.column + 1}
}
