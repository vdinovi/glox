package lox

import (
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

	_, err := l.scan.MatchUntil(isNotWhitespace)
	if err != nil {
		return nil, err
	}

	ch, err := l.scan.Next()
	if err != nil {
		return nil, err
	}

	switch ch {
	case '(', ')', '{', '}', ',', '.', '-', '+', ';', '*':
		token.Lexem = string(ch)
		token.Type = TokenTypeFor(token.Lexem)
	case '!', '=', '<', '>':
		nextCh, err := l.scan.MatchRune(isEquals)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF || nextCh != '=' {
			token.Lexem = string(ch)
		} else {
			token.Lexem = string(ch) + "="
		}
		token.Type = TokenTypeFor(token.Lexem)
	case '/':
		nextCh, err := l.scan.MatchRune(isSlash)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF || nextCh != '/' {
			token.Lexem = string(ch)
			token.Type = Slash
		} else {
			_, err := l.scan.MatchThrough(isNewline)
			if err != nil {
				return nil, err
			}
			return l.next()
		}
	case '"':
		chars, err := l.scan.MatchThrough(isQuote)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			return nil, &LexerError{
				Err:    &UnterminatedStringError{},
				Line:   line,
				Column: column,
			}
		} else {
			token.Lexem = string(chars[:len(chars)-1])
			token.Type = String
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		chars, err := l.scan.MatchUntil(isNotDigit)
		if err != nil && err != io.EOF {
			return nil, err
		}
		token.Lexem = string(ch) + string(chars)
		token.Type = Number
		if err != io.EOF {
			nextCh, err := l.scan.MatchRune(isDot)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if nextCh == '.' {
				chars, err := l.scan.MatchUntil(isNotDigit)
				if err != nil && err != io.EOF {
					return nil, err
				}
				token.Lexem += "." + string(chars)
				token.Type = Number
			}
		}
	default:
		chars, err := l.scan.MatchThrough(isNotLetterOrUnderscore)
		if err != nil && err != io.EOF {
			return nil, err
		}
		token.Lexem = string(ch) + string(chars)
		if tt := TokenTypeFor(token.Lexem); tt != None {
			token.Type = tt
		} else {
			for _, ch := range token.Lexem {
				if ch != '_' && !unicode.IsLetter(ch) {
					return nil, &LexerError{
						Err:    InvaldIdentifierError{token.Lexem},
						Line:   line,
						Column: column,
					}
				}
			}
			token.Type = Identifier
		}
	}
	return &token, nil
}
