package lox

import (
	"strings"
	"testing"
)

func TestTypecheckCustom(t *testing.T) {
	src := `
var a = 1;
var b = 2;
print a + b;
print a = "test";
	`
	tokens, err := Scan(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Unexpected error in %q: %s", src, err)
	}
	prog, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Unexpected error in %q: %s", src, err)
	}
	ctx := NewContext()
	err = ctx.TypeCheckProgram(prog)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTypecheckPrintStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt PrintStatement
		err  error
	}{
		{text: "print 1;", stmt: PrintStatement{expr: oneExpr()}},
		{text: "print \"str\";", stmt: PrintStatement{expr: strExpr()}},
		{text: "print true;", stmt: PrintStatement{expr: trueExpr()}},
		{text: "print false;", stmt: PrintStatement{expr: falseExpr()}},
		{text: "print nil;", stmt: PrintStatement{expr: nilExpr()}},
		{text: "print foo;", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 7})},
	}

	for _, test := range tests {
		program := typecheckTestParse(t, test.text, 1)
		var print *PrintStatement
		var ok bool
		if print, ok = program[0].(*PrintStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckPrintStatement(print)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck of statement %v to produce error %q, but got %q", test.stmt, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckExpressionStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt ExpressionStatement
		err  error
	}{
		{text: "1;", stmt: ExpressionStatement{expr: oneExpr()}},
		{text: "\"str\";", stmt: ExpressionStatement{expr: strExpr()}},
		{text: "true;", stmt: ExpressionStatement{expr: trueExpr()}},
		{text: "false;", stmt: ExpressionStatement{expr: falseExpr()}},
		{text: "nil;", stmt: ExpressionStatement{expr: nilExpr()}},
		{text: "foo;", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 1})},
	}

	for _, test := range tests {
		program := typecheckTestParse(t, test.text, 1)
		var expr *ExpressionStatement
		var ok bool
		if expr, ok = program[0].(*ExpressionStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckExpressionStatement(expr)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck of statement %v to produce error %q, but got %q", test.stmt, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeBlockStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt BlockStatement
		err  error
	}{
		{text: "{1;}", stmt: BlockStatement{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}}}},
		{text: "{\"str\";}", stmt: BlockStatement{stmts: []Statement{&ExpressionStatement{expr: strExpr()}}}},
		{text: "{true;}", stmt: BlockStatement{stmts: []Statement{&ExpressionStatement{expr: trueExpr()}}}},
		{text: "{false;}", stmt: BlockStatement{stmts: []Statement{&ExpressionStatement{expr: falseExpr()}}}},
		{text: "{nil;}", stmt: BlockStatement{stmts: []Statement{&ExpressionStatement{expr: nilExpr()}}}},
		{text: "{foo;}", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 2})},
		{text: "{var foo = 1; print foo;}", stmt: BlockStatement{stmts: []Statement{&DeclarationStatement{name: "foo", expr: oneExpr()}, &ExpressionStatement{expr: nilExpr()}}}},
	}

	for _, test := range tests {
		program := typecheckTestParse(t, test.text, 1)
		var block *BlockStatement
		var ok bool
		if block, ok = program[0].(*BlockStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckBlockStatement(block)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck of statement %v to produce error %q, but got %q", test.stmt, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckUnaryExpression(t *testing.T) {
	// TODO
}

func TypeCheckBinaryExpression(t *testing.T) {
	// TODO
}

func TypeCheckGroupingExpression(t *testing.T) {
	// TODO
}

func TypeCheckVariableExpression(t *testing.T) {
	// TODO
}

func typecheckTestParse(t *testing.T, text string, numStatements int) []Statement {
	t.Helper()
	tokens, err := Scan(strings.NewReader(text))
	if err != nil {
		t.Errorf("Unexpected error in %q: %s", text, err)
		return nil
	}
	program, err := Parse(tokens)
	if err != nil {
		t.Errorf("Unexpected error in %q: %s", text, err)
		return nil
	}
	if len(program) != numStatements {
		t.Errorf("%q should have parsed to %d statements, but got %d", text, numStatements, len(program))
		return nil
	}
	return program
}
