package lox

import (
	"fmt"
	"testing"
)

var tokens = []Token{
	{Type: TokenLeftParen, Lexem: "("},
	{Type: TokenRightParen, Lexem: ")"},
	{Type: TokenLeftBrace, Lexem: "{"},
	{Type: TokenRightBrace, Lexem: "}"},
	{Type: TokenComma, Lexem: ","},
	{Type: TokenDot, Lexem: "."},
	{Type: TokenMinus, Lexem: "-"},
	{Type: TokenPlus, Lexem: "+"},
	{Type: TokenSemicolon, Lexem: ";"},
	{Type: TokenStar, Lexem: "*"},
	{Type: TokenBang, Lexem: "!"},
	{Type: TokenBangEqual, Lexem: "!="},
	{Type: TokenEqual, Lexem: "="},
	{Type: TokenEqualEqual, Lexem: "=="},
	{Type: TokenLess, Lexem: "<"},
	{Type: TokenLessEqual, Lexem: "<="},
	{Type: TokenGreater, Lexem: ">"},
	{Type: TokenGreaterEqual, Lexem: ">="},
	{Type: TokenSlash, Lexem: "/"},
	{Type: TokenAnd, Lexem: "and"},
	{Type: TokenClass, Lexem: "class"},
	{Type: TokenElse, Lexem: "else"},
	{Type: TokenFalse, Lexem: "false"},
	{Type: TokenFun, Lexem: "fun"},
	{Type: TokenFor, Lexem: "for"},
	{Type: TokenIf, Lexem: "if"},
	{Type: TokenNil, Lexem: "nil"},
	{Type: TokenOr, Lexem: "or"},
	{Type: TokenPrint, Lexem: "print"},
	{Type: TokenReturn, Lexem: "return"},
	{Type: TokenSuper, Lexem: "super"},
	{Type: TokenThis, Lexem: "this"},
	{Type: TokenTrue, Lexem: "true"},
	{Type: TokenVar, Lexem: "var"},
	{Type: TokenWhile, Lexem: "while"},
}

func Test_TokenTypeFor(t *testing.T) {
	for _, token := range tokens {
		got := TokenTypeFor(token.Lexem)
		if got != token.Type {
			t.Errorf("Expected %s to yield %s, but got %s", token.Lexem, token.Type, got)
		}
	}
}

func Test_Token_String(t *testing.T) {
	for _, token := range tokens {
		want := fmt.Sprintf("%s(%q)", token.Type, token.Lexem)
		got := token.String()
		if got != want {
			t.Errorf("Expected %s to yield %s, but got %s", token.Type, want, got)
		}
	}
}

func Test_Token_Operator(t *testing.T) {

	tests := []struct {
		TokenType
		lexem string
		OperatorType
	}{
		{TokenPlus, "+", OpPlus},
		{TokenMinus, "-", OpMinus},
		{TokenStar, "*", OpMultiply},
		{TokenSlash, "/", OpDivide},
		{TokenEqualEqual, "==", OpEqualEquals},
		{TokenBangEqual, "!=", OpNotEquals},
		{TokenLess, "<", OpLess},
		{TokenLessEqual, "<=", OpLessEquals},
		{TokenGreater, ">", OpGreater},
		{TokenGreaterEqual, ">=", OpGreaterEquals},
	}
	for _, test := range tests {
		token := Token{Type: test.TokenType, Lexem: test.lexem}
		got, err := token.Operator()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}

		want := Operator{Type: test.OperatorType, Lexem: test.lexem}
		if *got != want {
			t.Errorf("Expected %s to yield %s, but got %s", token.Type, want, *got)
		}
	}
}
