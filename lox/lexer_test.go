package lox

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLexerExpression(t *testing.T) {
	input := "(1.23 + (2*3) / -4) + !\"test\" * (false)"
	want := []Token{
		{Type: TokenLeftParen, Lexem: "("},
		{Type: TokenNumber, Lexem: "1.23"},
		{Type: TokenPlus, Lexem: "+"},
		{Type: TokenLeftParen, Lexem: "("},
		{Type: TokenNumber, Lexem: "2"},
		{Type: TokenStar, Lexem: "*"},
		{Type: TokenNumber, Lexem: "3"},
		{Type: TokenRightParen, Lexem: ")"},
		{Type: TokenSlash, Lexem: "/"},
		{Type: TokenMinus, Lexem: "-"},
		{Type: TokenNumber, Lexem: "4"},
		{Type: TokenRightParen, Lexem: ")"},
		{Type: TokenPlus, Lexem: "+"},
		{Type: TokenBang, Lexem: "!"},
		{Type: TokenString, Lexem: "test"},
		{Type: TokenStar, Lexem: "*"},
		{Type: TokenLeftParen, Lexem: "("},
		{Type: TokenFalse, Lexem: "false"},
		{Type: TokenRightParen, Lexem: ")"},
		{Type: TokenEOF, Lexem: ""},
	}
	lexer := NewLexer(strings.NewReader(input))
	got, err := lexer.ScanTokens()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if reflect.DeepEqual(got, want) {
		t.Fatalf("Expected %q to yield tokens %v, but got %v", input, want, got)
	}
}

func TestLexerBasicTokens(t *testing.T) {
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
		test.token.Line = 1
		test.token.Column = 1

		t.Run(fmt.Sprintf("%q yields %s", test.text, test.token.Type.String()), func(t *testing.T) {
			lexer := NewLexer(strings.NewReader(test.text))
			tokens, err := lexer.ScanTokens()
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if len(tokens) != 2 {
				t.Fatalf("expected %q to yield %d tokens, got %d", test.text, 1, len(tokens))
			}
			if token := tokens[0]; token != test.token {
				t.Fatalf("expected %q to yield token %+v, got %+v ", test.text, test.token, token)
			}
			if token := tokens[1]; token != EofToken {
				t.Errorf("expected %q to yield implicit token %+v, got %+v ", test.text, EofToken, token)
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
	}

	for _, test := range tests {
		lexer := NewLexer(strings.NewReader(test.text))
		tokens, err := lexer.ScanTokens()
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
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

func TestUnterminatedString(t *testing.T) {
	var lexerError *LexerError
	var unterminatedStringError *UnterminatedStringError

	tests := []string{
		"\"string",
	}
	for _, text := range tests {
		lexer := NewLexer(strings.NewReader(text))

		_, err := lexer.ScanTokens()

		if err == nil {
			t.Error("expected error but got none")
			continue
		}
		if !errors.As(err, &lexerError) {
			t.Errorf("unexpected error: %s", err.Error())
			continue
		}

		if !errors.As(err, &unterminatedStringError) {
			t.Errorf("unexpected error: %s", err.Error())
		}
	}
}

func TestUnexpectedCharacter(t *testing.T) {
	var lexerError *LexerError
	var unexpectedCharacterError *UnexpectedCharacterError
	tests := []string{
		"?",
		"foo?",
	}

	for _, text := range tests {
		lexer := NewLexer(strings.NewReader(text))
		_, err := lexer.ScanTokens()

		if err == nil {
			t.Error("expected error but got none")
			continue
		}

		if !errors.As(err, &lexerError) {
			t.Errorf("unexpected error %s", err.Error())
			continue
		}

		if !errors.As(err, &unexpectedCharacterError) {
			t.Errorf("unexpected error %s", err.Error())
			continue
		}
		if unexpectedCharacterError.Char != '?' {
			t.Errorf("expected char %c, got %c", unexpectedCharacterError.Char, '?')
		}
	}
}

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestLexerAll(t *testing.T) {
	t.Skip("Broken atm")
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	matches, err := filepath.Glob("../test/**/*.lox")
	if err != nil {
		t.Fatal(err)
	}
	var tests map[string]struct {
		Tokens []Token
		Error  error
	}
	data, err := os.ReadFile("./lexer_test/TestLexerAll.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &tests); err != nil {
		t.Fatal(err)
	}

	for _, path := range matches {
		test := tests[path]
		err := func() error {
			t.Helper()
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			lexer := NewLexer(bufio.NewReader(file))
			lexer.SetFilename(filepath.Base(path))
			tokens, err := lexer.ScanTokens()
			if update {
				tests[path] = struct {
					Tokens []Token
					Error  error
				}{
					Tokens: tokens,
					Error:  err,
				}
				test = tests[path]
			}
			if err != test.Error {
				t.Errorf("(%s) unexpected error %s", path, err.Error())
			}
			for i, token := range tokens {
				if expected := test.Tokens[i]; token != expected {
					t.Errorf("(%s) expected token %v, got %v", path, expected, token)
					break
				}
			}
			return nil
		}()
		if err != nil {
			t.Error(err)
		}
	}

	if update {
		data, err := json.Marshal(tests)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile("./lexer_test/TestLexerAll.json", data, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
