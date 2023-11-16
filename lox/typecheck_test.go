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
		err  error
	}{
		{text: "print 1;"},
		{text: "print \"str\";"},
		{text: "print true;"},
		{text: "print false;"},
		{text: "print nil;"},
		{text: "print foo;", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 7})},
	}

	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
			continue
		}
		var print *PrintStatement
		var ok bool
		if print, ok = td.Program[0].(*PrintStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckPrintStatement(print)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", print, test.err, err)
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
		err  error
	}{
		{text: "1;"},
		{text: "\"str\";"},
		{text: "true;"},
		{text: "false;"},
		{text: "nil;"},
		{text: "foo;", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 1})},
	}

	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
			continue
		}
		var expr *ExpressionStatement
		var ok bool
		if expr, ok = td.Program[0].(*ExpressionStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckExpressionStatement(expr)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", expr, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckBlockStatement(t *testing.T) {
	tests := []struct {
		text string
		err  error
	}{
		{text: "{1;}"},
		{text: "{\"str\";}"},
		{text: "{true;}"},
		{text: "{false;}"},
		{text: "{nil;}"},
		{text: "{var foo = 1; print foo;}"},
		{text: "{foo;}", err: NewTypeError(NewUndefinedVariableError("foo"), Position{1, 2})},
	}

	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
			continue
		}
		var block *BlockStatement
		var ok bool
		if block, ok = td.Program[0].(*BlockStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckBlockStatement(block)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", block, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckConditionalStatement(t *testing.T) {
	tests := []struct {
		text string
		err  error
	}{
		{text: "if (true) 1; else 3.14;"},
		{text: "if (true) 1;"},
		{text: "if (true) {1;}"},
		{text: "if (true) if (false) 1; else 3.14; else 1;"},
		{text: "if (true) { if (false) 1; else 3.14;} else {1;}"},
	}

	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
			continue
		}
		var cond *ConditionalStatement
		var ok bool
		if cond, ok = td.Program[0].(*ConditionalStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckConditionalStatement(cond)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", cond, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckWhileStatement(t *testing.T) {
	tests := []struct {
		text string
		err  error
	}{
		{text: "while (true) 1;"},
		{text: "while (true) { 1; }"},
		{text: "while (true) { while (false) 1;}"},
	}

	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
		}
		var cond *WhileStatement
		var ok bool
		if cond, ok = td.Program[0].(*WhileStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckWhileStatement(cond)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", cond, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckForStatement(t *testing.T) {
	tests := []struct {
		text string
		err  error
	}{
		{text: "for (;;) 1;"},
		{text: "for (;;) { 1; }"},
		{text: "for (var x = 1;;) 1;"},
		{text: "for (var x = 1; x;) 1;"},
		{text: "for (var x = 1; x; x = x + 1) 1;"},
	}
	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		if td.Err != nil {
			t.Errorf("Unexpected error in %q while %s: %s", test.text, td.Phase(), td.Err)
			continue
		}
		if len(td.Program) != 1 {
			t.Errorf("Expected %q to generate %d statements but got %d", test.text, 1, len(td.Program))
			continue
		}
		var cond *ForStatement
		var ok bool
		if cond, ok = td.Program[0].(*ForStatement); !ok {
			t.Errorf("%q yielded unexpected statment %v", test.text, td.Program[0])
			continue
		}
		ctx := NewContext()
		err := ctx.TypeCheckForStatement(cond)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%s) to produce error %q, but got %q", cond, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestTypeCheckExpression(t *testing.T) {
	tests := []struct {
		typ  Type
		expr Expression
		err  error
	}{
		// primitives
		{typ: TypeNumeric, expr: oneExpr()},
		{typ: TypeNumeric, expr: piExpr()},
		{typ: TypeString, expr: strExpr()},
		{typ: TypeBoolean, expr: trueExpr()},
		{typ: TypeBoolean, expr: falseExpr()},
		{typ: TypeNil, expr: nilExpr()},
		// unary sub
		{typ: TypeNumeric, expr: uSubExpr(oneExpr())()},
		{typ: TypeNumeric, expr: uSubExpr(piExpr())()},
		{typ: TypeNumeric, expr: uSubExpr(uSubExpr(piExpr())())()},
		{expr: uSubExpr(strExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeString), Position{})},
		{expr: uSubExpr(trueExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeBoolean), Position{})},
		{expr: uSubExpr(nilExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeNil), Position{})},
		// unary add
		{typ: TypeNumeric, expr: uAddExpr(oneExpr())()},
		{typ: TypeNumeric, expr: uAddExpr(piExpr())()},
		{typ: TypeNumeric, expr: uAddExpr(uAddExpr(piExpr())())()},
		{expr: uAddExpr(strExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeString), Position{})},
		{expr: uAddExpr(trueExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeBoolean), Position{})},
		{expr: uAddExpr(nilExpr())(), err: NewTypeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeNil), Position{})},
		// binary add
		{typ: TypeNumeric, expr: bAddExpr(oneExpr())(piExpr())()},
		{typ: TypeNumeric, expr: bAddExpr(piExpr())(oneExpr())()},
		{typ: TypeString, expr: bAddExpr(strExpr())(strExpr())()},
		{expr: bAddExpr(oneExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeString), Position{})},
		{expr: bAddExpr(oneExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeBoolean), Position{})},
		{expr: bAddExpr(oneExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeNil), Position{})},
		{expr: bAddExpr(strExpr())(oneExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeNumeric), Position{})},
		{expr: bAddExpr(strExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeBoolean), Position{})},
		{expr: bAddExpr(strExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeNil), Position{})},
		{expr: bAddExpr(trueExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeBoolean, TypeBoolean), Position{})},
		{expr: bAddExpr(nilExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNil, TypeNil), Position{})},
		// binary subtract
		{typ: TypeNumeric, expr: bSubExpr(oneExpr())(piExpr())()},
		{typ: TypeNumeric, expr: bSubExpr(piExpr())(oneExpr())()},
		{expr: bSubExpr(oneExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeString), Position{})},
		{expr: bSubExpr(oneExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeBoolean), Position{})},
		{expr: bSubExpr(oneExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeNil), Position{})},
		{expr: bSubExpr(strExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeString, TypeString), Position{})},
		{expr: bSubExpr(trueExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeBoolean, TypeBoolean), Position{})},
		{expr: bSubExpr(nilExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNil, TypeNil), Position{})},
		// binary multiply
		{typ: TypeNumeric, expr: bMulExpr(oneExpr())(piExpr())()},
		{typ: TypeNumeric, expr: bMulExpr(piExpr())(oneExpr())()},
		{expr: bMulExpr(oneExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeString), Position{})},
		{expr: bMulExpr(oneExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeBoolean), Position{})},
		{expr: bMulExpr(oneExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeNil), Position{})},
		{expr: bMulExpr(strExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeString, TypeString), Position{})},
		{expr: bMulExpr(trueExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeBoolean, TypeBoolean), Position{})},
		{expr: bMulExpr(nilExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNil, TypeNil), Position{})},
		// binary divide
		{typ: TypeNumeric, expr: bDivExpr(oneExpr())(piExpr())()},
		{typ: TypeNumeric, expr: bDivExpr(piExpr())(oneExpr())()},
		{expr: bDivExpr(oneExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeString), Position{})},
		{expr: bDivExpr(oneExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeBoolean), Position{})},
		{expr: bDivExpr(oneExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeNil), Position{})},
		{expr: bDivExpr(strExpr())(strExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeString, TypeString), Position{})},
		{expr: bDivExpr(trueExpr())(trueExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeBoolean, TypeBoolean), Position{})},
		{expr: bDivExpr(nilExpr())(nilExpr())(), err: NewTypeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNil, TypeNil), Position{})},
		// grouping
		{typ: TypeNumeric, expr: groupExpr(oneExpr())()},
		{typ: TypeNumeric, expr: groupExpr(piExpr())()},
		{typ: TypeString, expr: groupExpr(strExpr())()},
		{typ: TypeBoolean, expr: groupExpr(trueExpr())()},
		{typ: TypeBoolean, expr: groupExpr(falseExpr())()},
		{typ: TypeNil, expr: groupExpr(nilExpr())()},
		// complex arith
		// (1 + (3.14 / (-1) - -1)) + (+3.14)
		{
			typ:  TypeNumeric,
			expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(groupExpr(bDivExpr(piExpr())(bSubExpr(groupExpr(uSubExpr(oneExpr())())())(uSubExpr(oneExpr())())())())())())())(groupExpr(uAddExpr(piExpr())())())(),
		},
		// (1 + (3.14 / (-1) - "str")) + (+3.14)
		{
			expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(groupExpr(bDivExpr(piExpr())(bSubExpr(groupExpr(uSubExpr(oneExpr())())())(strExpr())())())())())())(groupExpr(uAddExpr(piExpr())())())(),
			err:  NewTypeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeString), Position{}),
		},
		// "str" + ("str" + ("str" + "str"))
		{
			typ:  TypeString,
			expr: bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(strExpr())())())())())(),
		},
		// "str" + ("str" + (1 + "str"))
		{
			expr: bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(groupExpr(bAddExpr(oneExpr())(strExpr())())())())())(),
			err:  NewTypeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeString), Position{}),
		},
		// and
		{typ: TypeNumeric, expr: bAndExpr(oneExpr())(piExpr())()},
		{typ: TypeString, expr: bAndExpr(strExpr())(strExpr())()},
		{typ: TypeBoolean, expr: bAndExpr(trueExpr())(falseExpr())()},
		{typ: TypeNil, expr: bAndExpr(nilExpr())(nilExpr())()},
		{typ: TypeAny, expr: bAndExpr(oneExpr())(strExpr())()},
		// or
		{typ: TypeNumeric, expr: bOrExpr(oneExpr())(piExpr())()},
		{typ: TypeString, expr: bOrExpr(strExpr())(strExpr())()},
		{typ: TypeBoolean, expr: bOrExpr(trueExpr())(falseExpr())()},
		{typ: TypeNil, expr: bOrExpr(nilExpr())(nilExpr())()},
		{typ: TypeAny, expr: bOrExpr(oneExpr())(strExpr())()},
	}
	for _, test := range tests {
		ctx := NewContext()
		typ, err := test.expr.TypeCheck(ctx)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%v) to yield error %q, but got %q", test.expr, test.err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("Unexpected error while typechecking %q: %s", test.expr, err)
			continue
		}
		if typ != test.typ {
			t.Errorf("Expected typecheck(%q) to be of type %s, but got %s", test.expr, test.typ, typ)
		}
	}
}

func TestTypeCheckVariableExpression(t *testing.T) {
	tests := []struct {
		typ      Type
		bindings map[string]Type
		expr     Expression
		err      error
	}{
		{expr: fooExpr(), err: NewTypeError(NewUndefinedVariableError("foo"), Position{})},
		{typ: TypeNumeric, expr: fooExpr(), bindings: map[string]Type{"foo": TypeNumeric}},
		{typ: TypeString, expr: fooExpr(), bindings: map[string]Type{"foo": TypeString}},
		{typ: TypeBoolean, expr: fooExpr(), bindings: map[string]Type{"foo": TypeBoolean}},
		{typ: TypeNil, expr: fooExpr(), bindings: map[string]Type{"foo": TypeNil}},
	}
	for _, test := range tests {
		ctx := NewContext()
		for name, typ := range test.bindings {
			_, err := ctx.types.Set(name, typ)
			if err != nil {
				t.Errorf("Unexpected error while setting bindings for %q: %s", test.expr, err)
			}
		}

		typ, err := test.expr.TypeCheck(ctx)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck(%v) to yield error %q, but got %q", test.expr, test.err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("Unexpected error while typechecking %q: %s", test.expr, err)
			continue
		}
		if typ != test.typ {
			t.Errorf("Expected typecheck(%q) to be of type %s, but got %s", test.expr, test.typ, typ)
		}
	}
}

func BenchmarkTypecheckFixtureProgram(b *testing.B) {
	tokens, err := Scan(strings.NewReader(fixtureProgram))
	if err != nil {
		b.Errorf("Unexpected error lexing fixture 'program': %s", err)
	}
	program, err := Parse(tokens)
	if err != nil {
		b.Errorf("Unexpected error parsing fixture 'program': %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := NewContext()
		err := ctx.TypeCheckProgram(program)
		if err != nil {
			b.Errorf("Unexpected error typechecking fixture 'program': %s", err)
		}
	}
}
