package lox

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseExpressionStatement(t *testing.T) {
	type exprStmts []ExpressionStatement
	one := numExpr(1)
	pi := numExpr(3.14)
	str := strExpr("str")
	yes := boolExpr(true)
	no := boolExpr(false)
	none := nilExpr()
	foo := varExpr("foo")
	neg := unaryExpr(Operator{Type: OpSubtract, Lexem: "-"})
	bAdd := binaryExpr(Operator{Type: OpAdd, Lexem: "+"})
	bSub := binaryExpr(Operator{Type: OpSubtract, Lexem: "-"})
	bMul := binaryExpr(Operator{Type: OpMultiply, Lexem: "*"})
	bDiv := binaryExpr(Operator{Type: OpDivide, Lexem: "/"})
	group := groupExpr
	tests := []struct {
		text  string
		stmts []ExpressionStatement
		err   error
	}{
		{text: "1;", stmts: exprStmts{{expr: one()}}},
		{text: "3.14;", stmts: exprStmts{{expr: pi()}}},
		{text: "\"str\";", stmts: exprStmts{{expr: str()}}},
		{text: "true;", stmts: exprStmts{{expr: yes()}}},
		{text: "false;", stmts: exprStmts{{expr: no()}}},
		{text: "nil;", stmts: exprStmts{{expr: none()}}},
		{text: "//comment\n1;", stmts: exprStmts{{expr: one()}}},
		{text: "1;//comment\n", stmts: exprStmts{{expr: one()}}},
		{text: "foo;", stmts: exprStmts{{expr: foo()}}},
		{text: "-1;", stmts: exprStmts{{expr: neg(one())()}}},
		{text: "--1;", stmts: exprStmts{{expr: neg(neg(one())())()}}},
		{text: "(1);", stmts: exprStmts{{expr: group(one())()}}},
		{text: "(-1);", stmts: exprStmts{{expr: group(neg(one())())()}}},
		{text: "1 + 3.14;", stmts: exprStmts{{expr: bAdd(one())(pi())()}}},
		{text: "1 - -3.14;", stmts: exprStmts{{expr: bSub(one())(neg(pi())())()}}},
		{text: "-1 * 3.14;", stmts: exprStmts{{expr: bMul(neg(one())())(pi())()}}},
		{text: "-1 / -3.14;", stmts: exprStmts{{expr: bDiv(neg(one())())(neg(pi())())()}}},
		{text: "(1 + 3.14);", stmts: exprStmts{{expr: group(bAdd(one())(pi())())()}}},
		{text: "1 + (1 + 3.14);", stmts: exprStmts{{expr: bAdd(one())(group(bAdd(one())(pi())())())()}}},
		{text: "(1 + 3.14) + 1;", stmts: exprStmts{{expr: bAdd(group(bAdd(one())(pi())())())(one())()}}},
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
				t.Errorf("Expected %q to be %v, but got %v", test.text, want, got)
				break
			}
		}
	}
}

func TestParsePrintStatement(t *testing.T) {
	type printStmts []PrintStatement
	one := numExpr(1)
	pi := numExpr(3.14)
	str := strExpr("str")
	yes := boolExpr(true)
	no := boolExpr(false)
	none := nilExpr()
	foo := varExpr("foo")
	neg := unaryExpr(Operator{Type: OpSubtract, Lexem: "-"})
	bAdd := binaryExpr(Operator{Type: OpAdd, Lexem: "+"})
	group := groupExpr

	tests := []struct {
		text  string
		stmts []PrintStatement
		err   error
	}{
		{text: "print 1;", stmts: printStmts{{expr: one()}}},
		{text: "print 3.14;", stmts: printStmts{{expr: pi()}}},
		{text: "print \"str\";", stmts: printStmts{{expr: str()}}},
		{text: "print true;", stmts: printStmts{{expr: yes()}}},
		{text: "print false;", stmts: printStmts{{expr: no()}}},
		{text: "print nil;", stmts: printStmts{{expr: none()}}},
		{text: "//comment\nprint 1;", stmts: printStmts{{expr: one()}}},
		{text: "print 1;//comment\n", stmts: printStmts{{expr: one()}}},
		{text: "print foo;", stmts: printStmts{{expr: foo()}}},
		{text: "print -1;", stmts: printStmts{{expr: neg(one())()}}},
		{text: "print 1 + 3.14;", stmts: printStmts{{expr: bAdd(one())(pi())()}}},
		{text: "print (1);", stmts: printStmts{{expr: group(one())()}}},
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
				t.Errorf("Expected %q to be %v, but got %v", test.text, want, got)
				break
			}
		}
	}
}

func TestParseDeclarationStatement(t *testing.T) {
	type declStmts []DeclarationStatement
	one := numExpr(1)
	pi := numExpr(3.14)
	str := strExpr("str")
	yes := boolExpr(true)
	no := boolExpr(false)
	none := nilExpr()
	bAdd := binaryExpr(Operator{Type: OpAdd, Lexem: "+"})
	group := groupExpr

	tests := []struct {
		text  string
		stmts []DeclarationStatement
		err   error
	}{
		{text: "var foo;", stmts: declStmts{{name: "foo", expr: none()}}},
		{text: "var foo = 1;", stmts: declStmts{{name: "foo", expr: one()}}},
		{text: "var foo = 3.14;", stmts: declStmts{{name: "foo", expr: pi()}}},
		{text: "var foo = \"str\";", stmts: declStmts{{name: "foo", expr: str()}}},
		{text: "var foo = true;", stmts: declStmts{{name: "foo", expr: yes()}}},
		{text: "var foo = false;", stmts: declStmts{{name: "foo", expr: no()}}},
		{text: "var foo = nil;", stmts: declStmts{{name: "foo", expr: none()}}},
		{text: "//comment\nvar foo;", stmts: declStmts{{name: "foo", expr: none()}}},
		// TODO: Fixme
		//{text: "var foo;//comment\n", stmts: declStmts{{name: "foo", expr: one()}}},
		{text: "var foo = 1 + 3.14;", stmts: declStmts{{name: "foo", expr: bAdd(one())(pi())()}}},
		{text: "var foo = (1);", stmts: declStmts{{name: "foo", expr: group(one())()}}},
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
				t.Errorf("Expected %q to be %v, but got %v", test.text, want, got)
				break
			}
		}
	}
}

func TestParseBlockStatement(t *testing.T) {
	// TODO: Tests
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
