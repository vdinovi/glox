package lox

import (
	"testing"
)

func TestRuntime(t *testing.T) {
	rt := Runtime{}
	expr := BinaryExpression{
		op: Operator{Type: OpMultiply, Lexem: "*"},
		left: UnaryExpression{
			op:    Operator{Type: OpMinus, Lexem: "-"},
			right: NumericExpression(123),
		},
		right: GroupingExpression{
			expr: NumericExpression(45.67),
		},
	}
	val, typ, err := rt.Evaluate(expr)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	_ = val
	_ = typ

}
