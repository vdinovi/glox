package lox

import (
	"testing"
)

func TestSimpleExpression(t *testing.T) {
	tests := []struct {
		expr Expression
		val  Value
		typ  Type
		err  error
	}{
		{expr: oneExpr(), val: ValueNumeric(1), typ: TypeNumeric},
		{expr: piExpr(), val: ValueNumeric(3.14), typ: TypeNumeric},
		{expr: strExpr(), val: ValueString("str"), typ: TypeString},
		{expr: trueExpr(), val: ValueBoolean(true), typ: TypeBoolean},
		{expr: falseExpr(), val: ValueBoolean(false), typ: TypeBoolean},
		{expr: nilExpr(), val: ValueNil{}, typ: TypeNil},
		{expr: fooExpr(), val: ValueNil{}, typ: TypeNil},
	}
	for _, test := range tests {
		ctx := NewContext()
		typ, err := test.expr.TypeCheck(ctx)
		if err != nil {
			t.Errorf("Unexpected error while typechecking %q: %s", test.expr, err)
			continue
		}
		if typ != test.typ {
			t.Errorf("Expected typecheck(%q) to be of type %s, but got %s", test.expr, test.typ, typ)
			continue
		}
		val, err := test.expr.Evaluate(ctx)
		if err != nil {
			t.Errorf("Unexpected error while evaluating %q: %s", test.expr, err)
			continue
		}
		if val.Type() != test.typ {
			t.Errorf("Expected evaluate(%q) yield value of type %s, but got %s", test.expr, test.typ, val.Type())
			continue
		}
		if val != test.val {
			t.Errorf("Expected evaluate(%q) yield value %s, but got %s", test.expr, test.typ, val)
		}
	}
}
