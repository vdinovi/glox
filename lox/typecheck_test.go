package lox

import (
	"reflect"
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
	start := Position{1, 1}
	tests := []struct {
		err  error
		stmt PrintStatement
	}{
		{stmt: PrintStatement{expr: &NumericExpression{value: 1, pos: start}, pos: start}},
		{stmt: PrintStatement{expr: &StringExpression{value: "str", pos: start}, pos: start}},
		{stmt: PrintStatement{expr: &BooleanExpression{value: true, pos: start}, pos: start}},
		{stmt: PrintStatement{expr: &BooleanExpression{value: false, pos: start}, pos: start}},
		{stmt: PrintStatement{expr: &NilExpression{pos: start}, pos: start}},
		//{stmt: PrintStatement{expr: VariableExpression{name: "foo", pos: start}, pos: start}},
	}

	for _, test := range tests {
		ctx := NewContext()
		err := ctx.TypeCheckPrintStatement(&test.stmt)
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
	// TODO: Update to use new test helpers
	start := Position{1, 1}
	tests := []struct {
		err  error
		stmt ExpressionStatement
	}{
		{stmt: ExpressionStatement{expr: &NumericExpression{value: 1, pos: start}, pos: start}},
		{stmt: ExpressionStatement{expr: &StringExpression{value: "str", pos: start}, pos: start}},
		{stmt: ExpressionStatement{expr: &BooleanExpression{value: true, pos: start}, pos: start}},
		{stmt: ExpressionStatement{expr: &BooleanExpression{value: false, pos: start}, pos: start}},
		{stmt: ExpressionStatement{expr: &NilExpression{pos: start}, pos: start}},
		//{stmt: ExpressionStatement{expr: VariableExpression{name: "foo", pos: start}, pos: start}},
	}

	for _, test := range tests {
		ctx := NewContext()
		err := ctx.TypeCheckExpressionStatement(&test.stmt)
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
	add := Operator{Type: OpAdd, Lexem: "+"}
	sub := Operator{Type: OpSubtract, Lexem: "+"}
	pos := Position{1, 1}
	one := numExpr(1)
	pi := numExpr(3.14)
	tests := []struct {
		typ  Type
		err  error
		expr UnaryExpression
	}{
		{typ: TypeNumeric, expr: UnaryExpression{op: sub, pos: pos, right: one()}},
		{typ: TypeNumeric, expr: UnaryExpression{op: sub, pos: pos, right: pi()}},
		{typ: TypeNumeric, expr: UnaryExpression{op: add, pos: pos, right: one()}},
		{typ: TypeNumeric, expr: UnaryExpression{op: add, pos: pos, right: pi()}},
		// panic: runtime error: comparing uncomparable type lox.InvalidOperatorForTypeError
		// 	{err: NewTypeError(NewInvalidOperatorForTypeError(sub.Type, TypeString), pos), expr: UnaryExpression{op: sub, pos: pos, right: str}},
		// 	{err: NewTypeError(NewInvalidOperatorForTypeError(add.Type, TypeString), pos), expr: UnaryExpression{op: add, pos: pos, right: str}},
		// 	{err: NewTypeError(NewInvalidOperatorForTypeError(sub.Type, TypeBoolean), pos), expr: UnaryExpression{op: sub, pos: pos, right: yes}},
		// 	{err: NewTypeError(NewInvalidOperatorForTypeError(add.Type, TypeBoolean), pos), expr: UnaryExpression{op: add, pos: pos, right: yes}},
	}

	for _, test := range tests {
		ctx := NewContext()
		_, typ, err := ctx.TypeCheckUnaryExpression(&test.expr)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected typecheck of %v to produce error %q, but got %q", test.expr, test.err, err)
				continue
			}
		} else if err != nil {
			t.Error(err)
			continue
		}

		if !reflect.DeepEqual(typ, test.typ) {
			t.Errorf("Expected typecheck of %v to produce type %s, but got %s", test.expr, test.typ, typ)
		}
	}

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
