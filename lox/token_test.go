package lox

import (
	"errors"
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	tests := []struct {
		token Token
		lexem string
	}{
		{tokenDefault(TokenLeftParen), "("},
		{tokenDefault(TokenRightParen), ")"},
		{tokenDefault(TokenLeftBrace), "{"},
		{tokenDefault(TokenRightBrace), "}"},
		{tokenDefault(TokenComma), ","},
		{tokenDefault(TokenDot), "."},
		{tokenDefault(TokenMinus), "-"},
		{tokenDefault(TokenPlus), "+"},
		{tokenDefault(TokenSemicolon), ";"},
		{tokenDefault(TokenSlash), "/"},
		{tokenDefault(TokenStar), "*"},
		{tokenDefault(TokenBang), "!"},
		{tokenDefault(TokenBangEqual), "!="},
		{tokenDefault(TokenEqual), "="},
		{tokenDefault(TokenEqualEqual), "=="},
		{tokenDefault(TokenGreater), ">"},
		{tokenDefault(TokenGreaterEqual), ">="},
		{tokenDefault(TokenLess), "<"},
		{tokenDefault(TokenLessEqual), "<="},
		{tokenDefault(TokenAnd), "and"},
		{tokenDefault(TokenClass), "class"},
		{tokenDefault(TokenElse), "else"},
		{tokenDefault(TokenFalse), "false"},
		{tokenDefault(TokenFun), "fun"},
		{tokenDefault(TokenFor), "for"},
		{tokenDefault(TokenIf), "if"},
		{tokenDefault(TokenNil), "nil"},
		{tokenDefault(TokenOr), "or"},
		{tokenDefault(TokenPrint), "print"},
		{tokenDefault(TokenReturn), "return"},
		{tokenDefault(TokenSuper), "super"},
		{tokenDefault(TokenThis), "this"},
		{tokenDefault(TokenTrue), "true"},
		{tokenDefault(TokenVar), "var"},
		{tokenDefault(TokenWhile), "while"},
		{tokenDefault(TokenEOF), ""},
		{Token{Type: TokenString, Lexem: "str"}, "str"},
		{Token{Type: TokenNumber, Lexem: "1.234"}, "1.234"},
		{Token{Type: TokenIdentifier, Lexem: "foo"}, "foo"},
		{Token{Type: TokenComment, Lexem: "comment"}, "comment"},
	}
	for _, test := range tests {
		want := fmt.Sprintf("%s(%q)", test.token.Type, test.lexem)
		if got := test.token.String(); got != want {
			t.Errorf("Expected %v to have string %s, but got %s", test.token, want, got)
		}
	}
}

func TestTokenOperator(t *testing.T) {
	tests := []struct {
		token Token
		op    OperatorType
	}{
		{tokenDefault(TokenPlus), OpAdd},
		{tokenDefault(TokenMinus), OpSubtract},
		{tokenDefault(TokenSlash), OpDivide},
		{tokenDefault(TokenStar), OpMultiply},
		{tokenDefault(TokenBangEqual), OpNotEqualTo},
		{tokenDefault(TokenEqualEqual), OpEqualTo},
		{tokenDefault(TokenGreater), OpGreaterThan},
		{tokenDefault(TokenGreaterEqual), OpGreaterThanOrEqualTo},
		{tokenDefault(TokenLess), OpLessThan},
		{tokenDefault(TokenLessEqual), OpLessThanOrEqualTo},
	}
	for _, test := range tests {
		op, err := test.token.Operator()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if op.Type != test.op {
			t.Errorf("Expected %s to yield operator %s, but got %s", test.token, test.op, op.Type)
		}
	}

	_, err := tokenDefault(TokenEOF).Operator()
	if err == nil {
		t.Fatal("Expected error but got none")
	}

	want := NoOperatorForTokenError{TokenEOF}
	if !errors.Is(err, want) {
		t.Fatalf("Expected error %v, but got %v", want, err)
	}
}
