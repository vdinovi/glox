package lox

import (
	"errors"
	"testing"
)

// func TestParse(t *testing.T) {
// 	input := "(1.23 + (2*3) / -4) + \"test\" * (false)"
// 	tokens := []Token{
// 		{Type: TokenLeftParen, Lexem: "("},
// 		{Type: TokenNumber, Lexem: "1.23"},
// 		{Type: TokenPlus, Lexem: "+"},
// 		{Type: TokenLeftParen, Lexem: "("},
// 		{Type: TokenNumber, Lexem: "2"},
// 		{Type: TokenStar, Lexem: "*"},
// 		{Type: TokenNumber, Lexem: "3"},
// 		{Type: TokenRightParen, Lexem: ")"},
// 		{Type: TokenSlash, Lexem: "/"},
// 		{Type: TokenMinus, Lexem: "-"},
// 		{Type: TokenNumber, Lexem: "4"},
// 		{Type: TokenRightParen, Lexem: ")"},
// 		{Type: TokenPlus, Lexem: "+"},
// 		{Type: TokenString, Lexem: "test"},
// 		{Type: TokenStar, Lexem: "*"},
// 		{Type: TokenLeftParen, Lexem: "("},
// 		{Type: TokenFalse, Lexem: "false"},
// 		{Type: TokenRightParen, Lexem: ")"},
// 		{Type: TokenEOF, Lexem: ""},
// 	}
// 	parser := NewParser(tokens)
// 	expr, err := parser.Parse()
// 	if err != nil {
// 		t.Fatalf("Unexpected error: %s", err)
// 	}
// 	got := expr.String()
// 	want := "(+ (group (+ 1.23 (/ (group (* 2 3)) (- 4)))) (* test (group false)))"

// 	if got != want {
// 		t.Fatalf("Expected %q to yield expression %q, but got %q", input, want, got)
// 	}
// }

func Test_Parse_UnmatchedParen(t *testing.T) {
	tokens := []Token{
		{Type: TokenLeftParen, Lexem: "("},
		{Type: TokenNumber, Lexem: "1"},
		{Type: TokenPlus, Lexem: "+"},
		{Type: TokenNumber, Lexem: "2"},
		{Type: TokenEOF, Lexem: ""},
	}
	parser := NewParser(tokens)
	expr, err := parser.Parse()
	_ = expr
	if err == nil {
		t.Fatal("Expected SyntaxError but got none")
	}
	var syntaxErr SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("Expected SyntaxError but got %v", err)
	}
	var utErr UnmatchedTokenError
	if !errors.As(err, &utErr) {
		t.Fatal("Expected SyntaxError to wrap UnmatchedTokenError but did not")
	}
	if utErr.Token.Type != TokenLeftParen {
		t.Fatalf("Expected UnmatchedTokenError to be for %s but got %s", TokenLeftParen, utErr.Token.Type)
	}
}
