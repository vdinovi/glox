package lox

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestParseVariableDeclaration(t *testing.T) {
	start := Position{1, 1}
	tests := []struct {
		text  string
		stmts []DeclarationStatement
		err   error
	}{
		{
			text: "var one = 1;",
			stmts: []DeclarationStatement{
				{name: "one", pos: start, expr: NumericExpression{value: 1, pos: start}},
			},
		},
		{
			text: "var pi = 3.14;",
			stmts: []DeclarationStatement{
				{name: "pi", pos: start, expr: NumericExpression{value: 3.14, pos: start}},
			},
		},
		{
			text: "var neg_pi = -3.14;",
			stmts: []DeclarationStatement{
				{name: "neg_pi", pos: start,
					expr: UnaryExpression{
						op:    Operator{Type: OpSubtract, Lexem: "-"},
						pos:   start,
						right: NumericExpression{value: 3.14, pos: start},
					}},
			},
		},
		{
			text: "var str = \"string\";",
			stmts: []DeclarationStatement{
				{name: "str", pos: start, expr: StringExpression{value: "string", pos: start}},
			},
		},
		{
			text: "var yes = true;",
			stmts: []DeclarationStatement{
				{name: "yes", pos: start, expr: BooleanExpression{value: true, pos: start}},
			},
		},
		{
			text: "var no = false;",
			stmts: []DeclarationStatement{
				{name: "no", pos: start, expr: BooleanExpression{value: false, pos: start}},
			},
		},
		{
			text: "var null = nil;",
			stmts: []DeclarationStatement{
				{name: "null", pos: start, expr: NilExpression{pos: start}},
			},
		},
		{
			text: "var undefined;",
			stmts: []DeclarationStatement{
				{name: "undefined", pos: start, expr: NilExpression{pos: start}},
			},
		},
		{
			text: "var ;",
			err:  NewSyntaxError(NewUnexpectedTokenError("Identifier", Token{Type: TokenSemicolon, Lexem: ";", Position: start}), start),
		},
		{
			text: "var x = 1",
			// TODO: why is this reporting Position{0,0} and not Position{1,1}?
			err: NewSyntaxError(NewUnexpectedTokenError("Semicolon", Token{Type: TokenEOF, Position: Position{}}), Position{}),
		},
		// TODO: panics
		// {
		// 	text: "x = 1",
		// 	err:  NewSyntaxError(NewUnexpectedTokenError("Identifier", Token{Type: TokenSemicolon, Lexem: ";", Position: start}), start),
		// },
	}
	for _, test := range tests {
		lexer, err := NewLexer(strings.NewReader(test.text))
		if err != nil {
			t.Error(err)
			continue
		}
		tokens, err := lexer.Scan()
		if err != nil {
			t.Error(err)
			continue
		}
		parser := NewParser(tokens)
		program, err := parser.Parse()
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error (%v), but got (%v)", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; got != want {
				t.Errorf("Expected %s to have declaration %s, but got %s", test.text, want, got)
				continue
			}
		}
	}
}

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

func TestParserProgram(t *testing.T) {
	var tokens []Token
	err := json.Unmarshal([]byte(program_tokens), &tokens)
	if err != nil {
		t.Fatalf("Failed to deserialize tokens")
	}

	parser := NewParser(tokens)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	_, err = parser.Parse()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
