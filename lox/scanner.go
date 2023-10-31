package lox

import (
	"io"
)

type MatchFunc func(rune) bool

// TODO: Simplify interface and pull these out into auxilary methods
type Scanner interface {
	Next() (rune, error)
	Position() (int, int)
	MatchRune(MatchFunc) (rune, error)
	MatchUntil(MatchFunc) ([]rune, error)
	MatchThrough(MatchFunc) ([]rune, error)
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
	s.advance(1)
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

func (s *stringScanner) MatchRune(mf MatchFunc) (rune, error) {
	chars, err := s.Peek(1)
	if err != nil {
		return -1, err
	}
	ch := chars[0]
	if mf(ch) {
		if err := s.advance(1); err != nil {
			return ch, err
		}
	}
	return ch, nil
}

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

func (s *stringScanner) MatchUntil(mf MatchFunc) ([]rune, error) {
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

func (s *stringScanner) MatchThrough(mf MatchFunc) ([]rune, error) {
	chars, err := s.MatchUntil(mf)
	if err != nil {
		return chars, err
	}
	next, err := s.Next()
	if err != nil {
		return chars, err
	}
	return append(chars, next), nil
}

func (s *stringScanner) advance(size int) error {
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
