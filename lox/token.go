package lox

import "fmt"

//go:generate stringer -type TokenType -trimprefix=Token
type TokenType int

const (
	TokenNone TokenType = iota
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar
	TokenBang
	TokenBangEqual
	TokenEqual
	TokenEqualEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual
	TokenIdentifier
	TokenString
	TokenNumber
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFun
	TokenFor
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenThis
	TokenTrue
	TokenVar
	TokenWhile
	TokenComment
	TokenEOF
)

var tokenTypeMap = map[string]TokenType{
	"(":      TokenLeftParen,
	")":      TokenRightParen,
	"{":      TokenLeftBrace,
	"}":      TokenRightBrace,
	",":      TokenComma,
	".":      TokenDot,
	"-":      TokenMinus,
	"+":      TokenPlus,
	";":      TokenSemicolon,
	"*":      TokenStar,
	"!":      TokenBang,
	"!=":     TokenBangEqual,
	"=":      TokenEqual,
	"==":     TokenEqualEqual,
	"<":      TokenLess,
	"<=":     TokenLessEqual,
	">":      TokenGreater,
	">=":     TokenGreaterEqual,
	"/":      TokenSlash,
	"and":    TokenAnd,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"fun":    TokenFun,
	"for":    TokenFor,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"this":   TokenThis,
	"true":   TokenTrue,
	"var":    TokenVar,
	"while":  TokenWhile,
}

func TokenTypeFor(lexem string) TokenType {
	if tt, ok := tokenTypeMap[lexem]; ok {
		return tt
	} else {
		return TokenNone
	}
}

type Token struct {
	Type   TokenType
	Lexem  string
	Line   int
	Column int
}

var EofToken = Token{Type: TokenEOF}
var NoneToken = Token{Type: TokenNone}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", t.Type, t.Lexem)
}

var operatorTypeMap = map[TokenType]OperatorType{
	TokenPlus:         OpPlus,
	TokenMinus:        OpMinus,
	TokenStar:         OpMultiply,
	TokenSlash:        OpDivide,
	TokenEqual:        OpEquals,
	TokenEqualEqual:   OpEqualEquals,
	TokenBangEqual:    OpNotEquals,
	TokenLess:         OpLess,
	TokenLessEqual:    OpLessEquals,
	TokenGreater:      OpGreater,
	TokenGreaterEqual: OpGreaterEquals,
}

func (t Token) Operator() (*Operator, error) {
	opType, ok := operatorTypeMap[t.Type]
	if !ok {
		return nil, InvalidOperatorError{t}
	}
	return &Operator{
		Type:  opType,
		Lexem: t.Lexem,
	}, nil
}

type InvalidOperatorError struct {
	Token
}

func (e InvalidOperatorError) Error() string {
	return fmt.Sprintf("%s is not a valid operator", e.Token)
}
