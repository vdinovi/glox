package lox

import (
	"testing"
)

func TestExpression(t *testing.T) {
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

	want := "(* (- 123) (group 45.67))"
	got := expr.String()
	if got != want {
		t.Errorf("expected %s, but got %s", want, got)
	}
}

func Test_Expression_String(t *testing.T) {
	one := NumericExpression(1)
	two := NumericExpression(2)
	asdf := StringExpression("asdf")
	yes := BooleanExpression(true)
	minus := Operator{Type: OpMinus, Lexem: "-"}
	plus := Operator{Type: OpPlus, Lexem: "+"}
	multiply := Operator{Type: OpMultiply, Lexem: "*"}

	tests := []struct {
		Expression
		want string
	}{
		{UnaryExpression{op: minus, right: one}, "(- 1)"},
		{BinaryExpression{op: plus, left: one, right: two}, "(+ 1 2)"},
		{GroupingExpression{
			expr: BinaryExpression{op: multiply, left: yes, right: asdf},
		}, "(group (* true asdf))"},
	}

	for _, test := range tests {
		got := test.Expression.String()
		if got != test.want {
			t.Errorf("Expected %T to yield %s, but got %s", test.Expression, test.want, got)
		}
	}
}
