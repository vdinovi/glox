package lox

import (
	"fmt"
)

//go:generate stringer -type TokenType -trimprefix=Token
type TokenType int

const (
	// ErrToken is a sentinel value representing an invalid type
	ErrToken TokenType = iota
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

// Represents a valid syntax token
type Token struct {
	Type     TokenType `json:"Type"`     // type of token
	Lexem    string    `json:"Lexem"`    // associated string
	Position Position  `json:"Position"` // originating location
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", t.Type, t.Lexem)
}

// Returns the operator associated with this token.
// If there is no associated operator, returns the error NoOperatorForTokenError
func (t Token) Operator() (Operator, error) {
	var err error
	op := Operator{Type: ErrOp, Lexem: t.Lexem}
	switch t.Type {
	case TokenPlus:
		op.Type = OpAdd
	case TokenMinus:
		op.Type = OpSubtract
	case TokenStar:
		op.Type = OpMultiply
	case TokenSlash:
		op.Type = OpDivide
	case TokenEqualEqual:
		op.Type = OpEqualTo
	case TokenBangEqual:
		op.Type = OpNotEqualTo
	case TokenLess:
		op.Type = OpLessThan
	case TokenLessEqual:
		op.Type = OpLessThanOrEqualTo
	case TokenGreater:
		op.Type = OpGreaterThan
	case TokenGreaterEqual:
		op.Type = OpGreaterThanOrEqualTo
	default:
		err = NoOperatorForTokenError{t.Type}
	}
	return op, err
}

type Position struct {
	Line   int `json:"Line"`   // originating line
	Column int `json:"Column"` // originating column
}

type NoOperatorForTokenError struct {
	TokenType
}

func (e NoOperatorForTokenError) Error() string {
	return fmt.Sprintf("no operator associated with token %s", e.TokenType)
}

var eofToken = Token{Type: TokenEOF}

// var invalidToken = Token{Type: ErrToken}

func tokenTypeFor(lexem string) TokenType {
	switch lexem {
	case "(":
		return TokenLeftParen
	case ")":
		return TokenRightParen
	case "{":
		return TokenLeftBrace
	case "}":
		return TokenRightBrace
	case ",":
		return TokenComma
	case ".":
		return TokenDot
	case "-":
		return TokenMinus
	case "+":
		return TokenPlus
	case ";":
		return TokenSemicolon
	case "*":
		return TokenStar
	case "!":
		return TokenBang
	case "!=":
		return TokenBangEqual
	case "=":
		return TokenEqual
	case "==":
		return TokenEqualEqual
	case "<":
		return TokenLess
	case "<=":
		return TokenLessEqual
	case ">":
		return TokenGreater
	case ">=":
		return TokenGreaterEqual
	case "/":
		return TokenSlash
	case "and":
		return TokenAnd
	case "class":
		return TokenClass
	case "else":
		return TokenElse
	case "false":
		return TokenFalse
	case "fun":
		return TokenFun
	case "for":
		return TokenFor
	case "if":
		return TokenIf
	case "nil":
		return TokenNil
	case "or":
		return TokenOr
	case "print":
		return TokenPrint
	case "return":
		return TokenReturn
	case "super":
		return TokenSuper
	case "this":
		return TokenThis
	case "true":
		return TokenTrue
	case "var":
		return TokenVar
	case "while":
		return TokenWhile
	}
	return ErrToken
}

func tokenDefault(typ TokenType) Token {
	t := Token{Type: typ, Lexem: ""}
	switch typ {
	case TokenLeftParen:
		t.Lexem = "("
	case TokenRightParen:
		t.Lexem = ")"
	case TokenLeftBrace:
		t.Lexem = "{"
	case TokenRightBrace:
		t.Lexem = "}"
	case TokenComma:
		t.Lexem = ","
	case TokenDot:
		t.Lexem = "."
	case TokenMinus:
		t.Lexem = "-"
	case TokenPlus:
		t.Lexem = "+"
	case TokenSemicolon:
		t.Lexem = ";"
	case TokenSlash:
		t.Lexem = "/"
	case TokenStar:
		t.Lexem = "*"
	case TokenBang:
		t.Lexem = "!"
	case TokenBangEqual:
		t.Lexem = "!="
	case TokenEqual:
		t.Lexem = "="
	case TokenEqualEqual:
		t.Lexem = "=="
	case TokenGreater:
		t.Lexem = ">"
	case TokenGreaterEqual:
		t.Lexem = ">="
	case TokenLess:
		t.Lexem = "<"
	case TokenLessEqual:
		t.Lexem = "<="
	case TokenAnd:
		t.Lexem = "and"
	case TokenClass:
		t.Lexem = "class"
	case TokenElse:
		t.Lexem = "else"
	case TokenFalse:
		t.Lexem = "false"
	case TokenFun:
		t.Lexem = "fun"
	case TokenFor:
		t.Lexem = "for"
	case TokenIf:
		t.Lexem = "if"
	case TokenNil:
		t.Lexem = "nil"
	case TokenOr:
		t.Lexem = "or"
	case TokenPrint:
		t.Lexem = "print"
	case TokenReturn:
		t.Lexem = "return"
	case TokenSuper:
		t.Lexem = "super"
	case TokenThis:
		t.Lexem = "this"
	case TokenTrue:
		t.Lexem = "true"
	case TokenVar:
		t.Lexem = "var"
	case TokenWhile:
		t.Lexem = "while"
	}
	return t
}
