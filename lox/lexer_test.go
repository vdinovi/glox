package lox

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestLexerBasic(t *testing.T) {
	tests := []struct {
		text  string
		token Token
	}{
		{"(", Token{Type: TokenLeftParen, Lexem: "("}},
		{")", Token{Type: TokenRightParen, Lexem: ")"}},
		{"{", Token{Type: TokenLeftBrace, Lexem: "{"}},
		{"}", Token{Type: TokenRightBrace, Lexem: "}"}},
		{",", Token{Type: TokenComma, Lexem: ","}},
		{".", Token{Type: TokenDot, Lexem: "."}},
		{"-", Token{Type: TokenMinus, Lexem: "-"}},
		{"+", Token{Type: TokenPlus, Lexem: "+"}},
		{";", Token{Type: TokenSemicolon, Lexem: ";"}},
		{"*", Token{Type: TokenStar, Lexem: "*"}},
		{"!", Token{Type: TokenBang, Lexem: "!"}},
		{"=", Token{Type: TokenEqual, Lexem: "="}},
		{"!=", Token{Type: TokenBangEqual, Lexem: "!="}},
		{"==", Token{Type: TokenEqualEqual, Lexem: "=="}},
		{"<=", Token{Type: TokenLessEqual, Lexem: "<="}},
		{">=", Token{Type: TokenGreaterEqual, Lexem: ">="}},
		{"/", Token{Type: TokenSlash, Lexem: "/"}},
		{"\"string\"", Token{Type: TokenString, Lexem: "string"}},
		{"1234", Token{Type: TokenNumber, Lexem: "1234"}},
		{"1.234", Token{Type: TokenNumber, Lexem: "1.234"}},
		{"and", Token{Type: TokenAnd, Lexem: "and"}},
		{"class", Token{Type: TokenClass, Lexem: "class"}},
		{"else", Token{Type: TokenElse, Lexem: "else"}},
		{"false", Token{Type: TokenFalse, Lexem: "false"}},
		{"fun", Token{Type: TokenFun, Lexem: "fun"}},
		{"for", Token{Type: TokenFor, Lexem: "for"}},
		{"if", Token{Type: TokenIf, Lexem: "if"}},
		{"nil", Token{Type: TokenNil, Lexem: "nil"}},
		{"or", Token{Type: TokenOr, Lexem: "or"}},
		{"print", Token{Type: TokenPrint, Lexem: "print"}},
		{"return", Token{Type: TokenReturn, Lexem: "return"}},
		{"super", Token{Type: TokenSuper, Lexem: "super"}},
		{"this", Token{Type: TokenThis, Lexem: "this"}},
		{"true", Token{Type: TokenTrue, Lexem: "true"}},
		{"var", Token{Type: TokenVar, Lexem: "var"}},
		{"while", Token{Type: TokenWhile, Lexem: "while"}},
		{"foo", Token{Type: TokenIdentifier, Lexem: "foo"}},
		{"//comment", Token{Type: TokenComment, Lexem: "comment"}},
	}

	for _, test := range tests {
		test.token.Position = Position{1, 1}

		t.Run(fmt.Sprintf("%q yields %s", test.text, test.token.Type.String()), func(t *testing.T) {
			lexer, err := NewLexer(strings.NewReader(test.text))
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			tokens, err := lexer.Scan()
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if len(tokens) != 2 {
				t.Fatalf("Expected %q to yield %d tokens, got %d", test.text, 1, len(tokens))
			}
			if token := tokens[0]; token != test.token {
				t.Fatalf("Expected %q to yield token %+v, got %+v ", test.text, test.token, token)
			}
			if token := tokens[1]; token != eofToken {
				t.Errorf("Expected %q to yield implicit token %+v, got %+v ", test.text, eofToken, token)
			}
		})
	}
}

func TestLexerExpression(t *testing.T) {
	input := "(1.23 + (2*3) / -4) + !\"test\" * (false)"
	want := []Token{
		tokenDefault(TokenLeftParen),
		{Type: TokenNumber, Lexem: "1.23"},
		tokenDefault(TokenPlus),
		tokenDefault(TokenLeftParen),
		{Type: TokenNumber, Lexem: "1.23"},
		tokenDefault(TokenStar),
		{Type: TokenNumber, Lexem: "3"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenSlash),
		tokenDefault(TokenMinus),
		{Type: TokenNumber, Lexem: "4"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenPlus),
		tokenDefault(TokenBang),
		{Type: TokenString, Lexem: "test"},
		tokenDefault(TokenStar),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenFalse),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenEOF),
	}
	lexer, err := NewLexer(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	tokens, err := lexer.Scan()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if reflect.DeepEqual(tokens, want) {
		t.Fatalf("Expected %q to yield tokens %v, but got %v", input, want, tokens)
	}
}

func TestLexerIgnore(t *testing.T) {
	tests := []struct {
		text string
	}{
		{" "},
		{"\t"},
		{"\n"},
	}

	for _, test := range tests {
		lexer, err := NewLexer(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		tokens, err := lexer.Scan()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if len(tokens) != 1 {
			t.Errorf("Expected %q to yield %d tokens but got %d", test.text, 1, len(tokens))
			continue
		}
		if token := tokens[0]; token != eofToken {
			t.Errorf("Expected %q to yield implicit token %+v but got %+v ", test.text, eofToken, token)
		}
	}
}

func TestUnterminatedString(t *testing.T) {
	var SyntaxError SyntaxError
	var unterminatedStringError UnterminatedStringError

	tests := []string{
		"\"string",
	}
	for _, text := range tests {
		_, err := Scan(strings.NewReader(text))
		if err == nil {
			t.Error("Expected error but got none")
			continue
		}
		if !errors.As(err, &SyntaxError) {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		if !errors.As(err, &unterminatedStringError) {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

func TestUnexpectedCharacter(t *testing.T) {
	var syntaxError SyntaxError
	var unexpectedCharacterError UnexpectedCharacterError
	tests := []struct {
		text string
	}{
		{"?"},
		{"foo?"},
	}

	for _, test := range tests {
		_, err := Scan(strings.NewReader(test.text))
		if err == nil {
			t.Error("Expected error but got none")
			continue
		}
		if !errors.As(err, &syntaxError) {
			t.Errorf("Unexpected error %s", err)
			continue
		}
		if !errors.As(err, &unexpectedCharacterError) {
			t.Errorf("Unexpected error %s", err)
			continue
		}
		if unexpectedCharacterError.Actual != '?' {
			t.Errorf("Expected %c, got %c", unexpectedCharacterError.Actual, '?')
		}
	}
}

func TestLexerFixtureProgram(t *testing.T) {
	var want []Token
	err := json.Unmarshal([]byte(fixtureProgramTokens), &want)
	if err != nil {
		t.Fatalf("Failed to deserialize tokens")
	}

	lexer, err := NewLexer(strings.NewReader(fixtureProgram))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	tokens, err := lexer.Scan()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if len(tokens) != len(want) {
		t.Fatalf("Expected program to yield %d tokens, but got %d", len(want), len(tokens))
	}
	for i, w := range want {
		if tok := tokens[i]; tok != w {
			t.Fatalf("Expected tokens[%d] to be %v, but got %v", i, w, tok)
		}
	}
}

func BenchmarkLexerFixtureProgram(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Scan(strings.NewReader(fixtureProgram))
		if err != nil {
			b.Errorf("Unexpected error lexing fixture 'program': %s", err)
		}
	}
}
