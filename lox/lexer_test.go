package lox

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexerBasicTokens(t *testing.T) {
	tests := []struct {
		text  string
		token Token
	}{
		{"(", Token{Type: LeftParen, Lexem: "("}},
		{")", Token{Type: RightParen, Lexem: ")"}},
		{"{", Token{Type: LeftBrace, Lexem: "{"}},
		{"}", Token{Type: RightBrace, Lexem: "}"}},
		{",", Token{Type: Comma, Lexem: ","}},
		{".", Token{Type: Dot, Lexem: "."}},
		{"-", Token{Type: Minus, Lexem: "-"}},
		{"+", Token{Type: Plus, Lexem: "+"}},
		{";", Token{Type: Semicolon, Lexem: ";"}},
		{"*", Token{Type: Star, Lexem: "*"}},
		{"!", Token{Type: Bang, Lexem: "!"}},
		{"=", Token{Type: Equal, Lexem: "="}},
		{"!=", Token{Type: BangEqual, Lexem: "!="}},
		{"==", Token{Type: EqualEqual, Lexem: "=="}},
		{"<=", Token{Type: LessEqual, Lexem: "<="}},
		{">=", Token{Type: GreaterEqual, Lexem: ">="}},
		{"/", Token{Type: Slash, Lexem: "/"}},
		{"\"string\"", Token{Type: String, Lexem: "string"}},
		{"1234", Token{Type: Number, Lexem: "1234"}},
		{"1.234", Token{Type: Number, Lexem: "1.234"}},
		{"and", Token{Type: And, Lexem: "and"}},
		{"class", Token{Type: Class, Lexem: "class"}},
		{"else", Token{Type: Else, Lexem: "else"}},
		{"false", Token{Type: False, Lexem: "false"}},
		{"fun", Token{Type: Fun, Lexem: "fun"}},
		{"for", Token{Type: For, Lexem: "for"}},
		{"if", Token{Type: If, Lexem: "if"}},
		{"nil", Token{Type: Nil, Lexem: "nil"}},
		{"or", Token{Type: Or, Lexem: "or"}},
		{"print", Token{Type: Print, Lexem: "print"}},
		{"return", Token{Type: Return, Lexem: "return"}},
		{"super", Token{Type: Super, Lexem: "super"}},
		{"this", Token{Type: This, Lexem: "this"}},
		{"true", Token{Type: True, Lexem: "true"}},
		{"var", Token{Type: Var, Lexem: "var"}},
		{"while", Token{Type: While, Lexem: "while"}},
		{"foo", Token{Type: Identifier, Lexem: "foo"}},
	}

	for _, test := range tests {
		test.token.Line = 1
		test.token.Column = 1

		t.Run(fmt.Sprintf("%q yields %s", test.text, test.token.Type.String()), func(t *testing.T) {
			lexer := NewLexer(strings.NewReader(test.text))
			tokens, err := lexer.ScanTokens()
			if err != nil {
				t.Fatalf("unexpected error: %q", err)
			}
			if len(tokens) != 2 {
				t.Fatalf("expected %q to yield %d tokens but got %d", test.text, 1, len(tokens))
			}
			if token := tokens[0]; token != test.token {
				t.Fatalf("expected %q to yield token %+v but got %+v ", test.text, test.token, token)
			}
			if token := tokens[1]; token != EofToken {
				t.Errorf("expected %q to yield implicit token %+v but got %+v ", test.text, EofToken, token)
			}
		})
	}
}

func TestLexerIgnore(t *testing.T) {
	tests := []struct {
		text string
	}{
		{" "},
		{"\t"},
		{"\n"},
		{"//comment"},
	}

	for _, test := range tests {
		lexer := NewLexer(strings.NewReader(test.text))
		tokens, err := lexer.ScanTokens()
		if err != nil {
			t.Errorf("unexpected error: %q", err)
			continue
		}
		if len(tokens) != 1 {
			t.Errorf("expected %q to yield %d tokens but got %d", test.text, 1, len(tokens))
			continue
		}
		if token := tokens[0]; token != EofToken {
			t.Errorf("expected %q to yield implicit token %+v but got %+v ", test.text, EofToken, token)
		}
	}
}

func TestLexerError(t *testing.T) {
	tests := []struct {
		text  string
		error error
	}{
		{"\"unterminated string", UnterminatedStringError{}},
	}

	for _, test := range tests {
		lexer := NewLexer(strings.NewReader(test.text))
		_, err := lexer.ScanTokens()

		lexErr, ok := err.(*LexerError)
		if !ok {
			t.Errorf("unexpected error %s", err)
			continue
		}

		_, ok = lexErr.Err.(*UnterminatedStringError)
		if !ok {
			t.Errorf("unexpected error %s", lexErr)
			continue
		}
	}
}
