package lox

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseExpressionStatement(t *testing.T) {
	type exprStmts []ExpressionStatement
	tests := []struct {
		text  string
		stmts []ExpressionStatement
		err   error
	}{
		{text: "1;", stmts: exprStmts{{expr: oneExpr()}}},
		{text: "3.14;", stmts: exprStmts{{expr: piExpr()}}},
		{text: "\"str\";", stmts: exprStmts{{expr: strExpr()}}},
		{text: "true;", stmts: exprStmts{{expr: trueExpr()}}},
		{text: "false;", stmts: exprStmts{{expr: falseExpr()}}},
		{text: "nil;", stmts: exprStmts{{expr: nilExpr()}}},
		{text: "//comment\n1;", stmts: exprStmts{{expr: oneExpr()}}},
		{text: "1;//comment\n", stmts: exprStmts{{expr: oneExpr()}}},
		{text: "foo;", stmts: exprStmts{{expr: fooExpr()}}},
		{text: "-1;", stmts: exprStmts{{expr: uSubExpr(oneExpr())()}}},
		{text: "--1;", stmts: exprStmts{{expr: uSubExpr(uSubExpr(oneExpr())())()}}},
		{text: "(1);", stmts: exprStmts{{expr: groupExpr(oneExpr())()}}},
		{text: "(-1);", stmts: exprStmts{{expr: groupExpr(uSubExpr(oneExpr())())()}}},
		{text: "1 + 3.14;", stmts: exprStmts{{expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "1 - -3.14;", stmts: exprStmts{{expr: bSubExpr(oneExpr())(uSubExpr(piExpr())())()}}},
		{text: "-1 * 3.14;", stmts: exprStmts{{expr: bMulExpr(uSubExpr(oneExpr())())(piExpr())()}}},
		{text: "-1 / -3.14;", stmts: exprStmts{{expr: bDivExpr(uSubExpr(oneExpr())())(uSubExpr(piExpr())())()}}},
		{text: "(1 + 3.14);", stmts: exprStmts{{expr: groupExpr(bAddExpr(oneExpr())(piExpr())())()}}},
		{text: "1 + (1 + 3.14);", stmts: exprStmts{{expr: bAddExpr(oneExpr())(groupExpr(bAddExpr(oneExpr())(piExpr())())())()}}},
		{text: "(1 + 3.14) + 1;", stmts: exprStmts{{expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(piExpr())())())(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParsePrintStatement(t *testing.T) {
	type printStmts []PrintStatement
	tests := []struct {
		text  string
		stmts []PrintStatement
		err   error
	}{
		{text: "print 1;", stmts: printStmts{{expr: oneExpr()}}},
		{text: "print 3.14;", stmts: printStmts{{expr: piExpr()}}},
		{text: "print \"str\";", stmts: printStmts{{expr: strExpr()}}},
		{text: "print true;", stmts: printStmts{{expr: trueExpr()}}},
		{text: "print false;", stmts: printStmts{{expr: falseExpr()}}},
		{text: "print nil;", stmts: printStmts{{expr: nilExpr()}}},
		{text: "//comment\nprint 1;", stmts: printStmts{{expr: oneExpr()}}},
		{text: "print 1;//comment\n", stmts: printStmts{{expr: oneExpr()}}},
		{text: "print foo;", stmts: printStmts{{expr: fooExpr()}}},
		{text: "print -1;", stmts: printStmts{{expr: uSubExpr(oneExpr())()}}},
		{text: "print 1 + 3.14;", stmts: printStmts{{expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "print (1);", stmts: printStmts{{expr: groupExpr(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParseDeclarationStatement(t *testing.T) {
	type declStmts []DeclarationStatement
	tests := []struct {
		text  string
		stmts []DeclarationStatement
		err   error
	}{
		{text: "var foo;", stmts: declStmts{{name: "foo", expr: nilExpr()}}},
		{text: "var foo = 1;", stmts: declStmts{{name: "foo", expr: oneExpr()}}},
		{text: "var foo = 3.14;", stmts: declStmts{{name: "foo", expr: piExpr()}}},
		{text: "var foo = \"str\";", stmts: declStmts{{name: "foo", expr: strExpr()}}},
		{text: "var foo = true;", stmts: declStmts{{name: "foo", expr: trueExpr()}}},
		{text: "var foo = false;", stmts: declStmts{{name: "foo", expr: falseExpr()}}},
		{text: "var foo = nil;", stmts: declStmts{{name: "foo", expr: nilExpr()}}},
		{text: "//comment\nvar foo;", stmts: declStmts{{name: "foo", expr: nilExpr()}}},
		{text: "var foo;//comment\n", stmts: declStmts{{name: "foo", expr: nilExpr()}}},
		{text: "var foo = 1 + 3.14;", stmts: declStmts{{name: "foo", expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "var foo = (1);", stmts: declStmts{{name: "foo", expr: groupExpr(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParseBlockStatement(t *testing.T) {
	tests := []struct {
		text  string
		stmts []BlockStatement
		err   error
	}{
		{text: "{var foo;}", stmts: []BlockStatement{{stmts: []Statement{&DeclarationStatement{name: "foo", expr: nilExpr()}}}}},
		{text: "{1;}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}}}}},
		{text: "{print 1;}", stmts: []BlockStatement{{stmts: []Statement{&PrintStatement{expr: oneExpr()}}}}},
		{text: "{1; 1;}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}, &ExpressionStatement{expr: oneExpr()}}}}},
		{text: "{1; {1;}}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}, &BlockStatement{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}}}}}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %s, but got %s", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParserProgram(t *testing.T) {
	// TODO: Needs to serialize AST to golden file for this test to work
	t.Skip()
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
