package lox

import (
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
		test.token.Line = 1
		test.token.Column = 1

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
	var SyntaxError *SyntaxError
	var unterminatedStringError *UnterminatedStringError

	tests := []string{
		"\"string",
	}
	for _, text := range tests {
		lexer, err := NewLexer(strings.NewReader(text))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		_, err = lexer.Scan()
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
	var syntaxError *SyntaxError
	var unexpectedCharacterError *UnexpectedCharacterError
	tests := []struct {
		text string
	}{
		{"?"},
		{"foo?"},
	}

	for _, test := range tests {
		lexer, err := NewLexer(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error %s", err)
			continue
		}
		_, err = lexer.Scan()
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

func TestLexerComplex(t *testing.T) {
	input := `
var one = 1;
var str = "str";
var null = nil;
var yes = true
var undefined;

print str;
print one + 2 ;
(1.23 + (one*3) / -4) + !\"test\" * (false);


// performs arithmetic on stuff
fun arith(a, b, c, d) {
	return (a + (b - c)) * d / a;
} 

arith(one, 2, yes, str)

// compares stuff
fun compare(a, b, c, d) {
	return (a > b) >= c < d <= (a + b) < c != (a == c);
}

compare(-1.23, yes, nil, undefined)

// does conditional stuff
fun conditional(a, b, c) {
	while (c < 5) {
		print c;
		c = c + 1;
	}

	for (d = 0; d < 5; d = d + 1) {
		print d;
	}

	if a < 1 {
		return a;
	} else if a >= 100 {
		return b;
	} else {
		return nil;
	}
}

class Foo {
	init(x) {
		this.x = x;
	}

	print() {
		print this.x;
	}
}

class Bar < Foo {
	init(y) {
		super().init("foo");
		this.y = y;
	}

	print() {
		super().print()
		print this.y;
	}
}

var foo = Foo("foo");
foo.print()

var bar = Bar("bar");
bar.print()
	`
	// Generated with
	// for _, tok := range tokens {
	// 	switch tok.Type {
	// 	case lox.TokenString, lox.TokenNumber, lox.TokenIdentifier, lox.TokenComment:
	// 		fmt.Printf("\t\t{Type: Token%s, Lexem: %q},\n", tok.Type, tok.Lexem)
	// 	default:
	// 		fmt.Printf("\t\ttokenDefault(Token%s),\n", tok.Type)
	// 	}
	// }
	want := []Token{
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "one"},
		tokenDefault(TokenEqual),
		{Type: TokenNumber, Lexem: "1"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "str"},
		tokenDefault(TokenEqual),
		{Type: TokenString, Lexem: "str"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "null"},
		tokenDefault(TokenEqual),
		tokenDefault(TokenNil),
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "yes"},
		tokenDefault(TokenEqual),
		tokenDefault(TokenTrue),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "undefined"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenPrint),
		{Type: TokenIdentifier, Lexem: "str"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenPrint),
		{Type: TokenIdentifier, Lexem: "one"},
		tokenDefault(TokenPlus),
		{Type: TokenNumber, Lexem: "2"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenLeftParen),
		{Type: TokenNumber, Lexem: "1.23"},
		tokenDefault(TokenPlus),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "one"},
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
		tokenDefault(TokenSemicolon),
		{Type: TokenComment, Lexem: " performs arithmetic on stuff"},
		tokenDefault(TokenFun),
		{Type: TokenIdentifier, Lexem: "arith"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenReturn),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenPlus),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenMinus),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenStar),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenSlash),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		{Type: TokenIdentifier, Lexem: "arith"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "one"},
		tokenDefault(TokenComma),
		{Type: TokenNumber, Lexem: "2"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "yes"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "str"},
		tokenDefault(TokenRightParen),
		{Type: TokenComment, Lexem: " compares stuff"},
		tokenDefault(TokenFun),
		{Type: TokenIdentifier, Lexem: "compare"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenReturn),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenGreater),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenGreaterEqual),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenLess),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenLessEqual),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenPlus),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLess),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenBangEqual),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenEqualEqual),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		{Type: TokenIdentifier, Lexem: "compare"},
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenMinus),
		{Type: TokenNumber, Lexem: "1.23"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "yes"},
		tokenDefault(TokenComma),
		tokenDefault(TokenNil),
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "undefined"},
		tokenDefault(TokenRightParen),
		{Type: TokenComment, Lexem: " does conditional stuff"},
		tokenDefault(TokenFun),
		{Type: TokenIdentifier, Lexem: "conditional"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenComma),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenWhile),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenLess),
		{Type: TokenNumber, Lexem: "5"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenPrint),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenSemicolon),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "c"},
		tokenDefault(TokenPlus),
		{Type: TokenNumber, Lexem: "1"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenFor),
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenEqual),
		{Type: TokenNumber, Lexem: "0"},
		tokenDefault(TokenSemicolon),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenLess),
		{Type: TokenNumber, Lexem: "5"},
		tokenDefault(TokenSemicolon),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenPlus),
		{Type: TokenNumber, Lexem: "1"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenPrint),
		{Type: TokenIdentifier, Lexem: "d"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenIf),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenLess),
		{Type: TokenNumber, Lexem: "1"},
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenReturn),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenElse),
		tokenDefault(TokenIf),
		{Type: TokenIdentifier, Lexem: "a"},
		tokenDefault(TokenGreaterEqual),
		{Type: TokenNumber, Lexem: "100"},
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenReturn),
		{Type: TokenIdentifier, Lexem: "b"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenElse),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenReturn),
		tokenDefault(TokenNil),
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenClass),
		{Type: TokenIdentifier, Lexem: "Foo"},
		tokenDefault(TokenLeftBrace),
		{Type: TokenIdentifier, Lexem: "init"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "x"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenThis),
		tokenDefault(TokenDot),
		{Type: TokenIdentifier, Lexem: "x"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "x"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenPrint),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenPrint),
		tokenDefault(TokenThis),
		tokenDefault(TokenDot),
		{Type: TokenIdentifier, Lexem: "x"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenClass),
		{Type: TokenIdentifier, Lexem: "Bar"},
		tokenDefault(TokenLess),
		{Type: TokenIdentifier, Lexem: "Foo"},
		tokenDefault(TokenLeftBrace),
		{Type: TokenIdentifier, Lexem: "init"},
		tokenDefault(TokenLeftParen),
		{Type: TokenIdentifier, Lexem: "y"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenSuper),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenDot),
		{Type: TokenIdentifier, Lexem: "init"},
		tokenDefault(TokenLeftParen),
		{Type: TokenString, Lexem: "foo"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenThis),
		tokenDefault(TokenDot),
		{Type: TokenIdentifier, Lexem: "y"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "y"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenPrint),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenLeftBrace),
		tokenDefault(TokenSuper),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenDot),
		tokenDefault(TokenPrint),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenPrint),
		tokenDefault(TokenThis),
		tokenDefault(TokenDot),
		{Type: TokenIdentifier, Lexem: "y"},
		tokenDefault(TokenSemicolon),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenRightBrace),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "foo"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "Foo"},
		tokenDefault(TokenLeftParen),
		{Type: TokenString, Lexem: "foo"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenSemicolon),
		{Type: TokenIdentifier, Lexem: "foo"},
		tokenDefault(TokenDot),
		tokenDefault(TokenPrint),
		tokenDefault(TokenLeftParen),
		tokenDefault(TokenRightParen),
		tokenDefault(TokenVar),
		{Type: TokenIdentifier, Lexem: "bar"},
		tokenDefault(TokenEqual),
		{Type: TokenIdentifier, Lexem: "Bar"},
		tokenDefault(TokenLeftParen),
		{Type: TokenString, Lexem: "bar"},
		tokenDefault(TokenRightParen),
		tokenDefault(TokenSemicolon),
		{Type: TokenIdentifier, Lexem: "bar"},
		tokenDefault(TokenDot),
		tokenDefault(TokenPrint),
		tokenDefault(TokenLeftParen),
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
	if len(tokens) != len(want) {
		t.Fatalf("Expected program to yield %d tokens, but got %d", len(want), len(tokens))
	}
	for i, w := range want {
		if tok := tokens[i]; tok != w {
			t.Fatalf("Expected tokens[%d] to be %v, but got %v", i, w, tok)
		}
	}
}
