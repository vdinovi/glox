package lox

import (
	"testing"
)

func TestSimpleExpression(t *testing.T) {
	tests := []struct {
		expr Expression
		val  Value
		err  error
	}{
		{expr: uSubExpr(strExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeString), Position{})},
		{val: ValueNumeric(3.14), expr: piExpr()},
		{val: ValueString("str"), expr: strExpr()},
		{val: ValueBoolean(true), expr: trueExpr()},
		{val: ValueBoolean(false), expr: falseExpr()},
		{val: ValueNil{}, expr: nilExpr()},
		{val: ValueNumeric(-1), expr: uSubExpr(oneExpr())()},
		{val: ValueNumeric(-3.14), expr: uSubExpr(piExpr())()},
		{val: ValueNumeric(3.14), expr: uSubExpr(uSubExpr(piExpr())())()},
		{expr: uSubExpr(strExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeString), Position{})},
		{expr: uSubExpr(trueExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeBoolean), Position{})},
		{expr: uSubExpr(nilExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpSubtract, TypeNil), Position{})},
		// unary negate
		{val: ValueBoolean(false), expr: uNegExpr(oneExpr())()},
		{val: ValueBoolean(false), expr: uNegExpr(piExpr())()},
		{val: ValueBoolean(false), expr: uNegExpr(strExpr())()},
		{val: ValueBoolean(false), expr: uNegExpr(trueExpr())()},
		{val: ValueBoolean(true), expr: uNegExpr(falseExpr())()},
		{val: ValueBoolean(true), expr: uNegExpr(nilExpr())()},
		{val: ValueBoolean(true), expr: uNegExpr(uNegExpr(oneExpr())())()},
		{val: ValueBoolean(true), expr: uNegExpr(uNegExpr(piExpr())())()},
		{val: ValueBoolean(true), expr: uNegExpr(uNegExpr(strExpr())())()},
		{val: ValueBoolean(true), expr: uNegExpr(uNegExpr(trueExpr())())()},
		{val: ValueBoolean(false), expr: uNegExpr(uNegExpr(falseExpr())())()},
		{val: ValueBoolean(false), expr: uNegExpr(uNegExpr(nilExpr())())()},
		// unary add
		{val: ValueNumeric(1), expr: uAddExpr(oneExpr())()},
		{val: ValueNumeric(3.14), expr: uAddExpr(piExpr())()},
		{val: ValueNumeric(3.14), expr: uAddExpr(uAddExpr(piExpr())())()},
		{expr: uAddExpr(strExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeString), Position{})},
		{expr: uAddExpr(trueExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeBoolean), Position{})},
		{expr: uAddExpr(nilExpr())(), err: NewRuntimeError(NewInvalidUnaryOperatorForTypeError(OpAdd, TypeNil), Position{})},
		// binary add
		{val: ValueNumeric(1 + 3.14), expr: bAddExpr(oneExpr())(piExpr())()},
		{val: ValueNumeric(3.14 + 1), expr: bAddExpr(piExpr())(oneExpr())()},
		{val: ValueString("strstr"), expr: bAddExpr(strExpr())(strExpr())()},
		{expr: bAddExpr(oneExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeString), Position{})},
		{expr: bAddExpr(oneExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeBoolean), Position{})},
		{expr: bAddExpr(oneExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeNil), Position{})},
		{expr: bAddExpr(strExpr())(oneExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeNumeric), Position{})},
		{expr: bAddExpr(strExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeBoolean), Position{})},
		{expr: bAddExpr(strExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeString, TypeNil), Position{})},
		{expr: bAddExpr(trueExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeBoolean, TypeBoolean), Position{})},
		{expr: bAddExpr(nilExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNil, TypeNil), Position{})},
		// binary subtract
		{val: ValueNumeric(1 - 3.14), expr: bSubExpr(oneExpr())(piExpr())()},
		{val: ValueNumeric(3.14 - 1), expr: bSubExpr(piExpr())(oneExpr())()},
		{expr: bSubExpr(oneExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeString), Position{})},
		{expr: bSubExpr(oneExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeBoolean), Position{})},
		{expr: bSubExpr(oneExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeNil), Position{})},
		{expr: bSubExpr(strExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeString, TypeString), Position{})},
		{expr: bSubExpr(trueExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeBoolean, TypeBoolean), Position{})},
		{expr: bSubExpr(nilExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNil, TypeNil), Position{})},
		// binary multiply
		{val: ValueNumeric(1 * 3.14), expr: bMulExpr(oneExpr())(piExpr())()},
		{val: ValueNumeric(3.14 * 1), expr: bMulExpr(piExpr())(oneExpr())()},
		{expr: bMulExpr(oneExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeString), Position{})},
		{expr: bMulExpr(oneExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeBoolean), Position{})},
		{expr: bMulExpr(oneExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNumeric, TypeNil), Position{})},
		{expr: bMulExpr(strExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeString, TypeString), Position{})},
		{expr: bMulExpr(trueExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeBoolean, TypeBoolean), Position{})},
		{expr: bMulExpr(nilExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpMultiply, TypeNil, TypeNil), Position{})},
		// binary divide
		{val: ValueNumeric(1 / 3.14), expr: bDivExpr(oneExpr())(piExpr())()},
		{val: ValueNumeric(3.14 / 1), expr: bDivExpr(piExpr())(oneExpr())()},
		{expr: bDivExpr(oneExpr())(zeroExpr())(), err: NewRuntimeError(NewDivideByZeroError(ValueNumeric(1), ValueNumeric(0)), Position{})},
		{expr: bDivExpr(oneExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeString), Position{})},
		{expr: bDivExpr(oneExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeBoolean), Position{})},
		{expr: bDivExpr(oneExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNumeric, TypeNil), Position{})},
		{expr: bDivExpr(strExpr())(strExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeString, TypeString), Position{})},
		{expr: bDivExpr(trueExpr())(trueExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeBoolean, TypeBoolean), Position{})},
		{expr: bDivExpr(nilExpr())(nilExpr())(), err: NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpDivide, TypeNil, TypeNil), Position{})},
		// grouping
		{val: ValueNumeric(1), expr: groupExpr(oneExpr())()},
		{val: ValueNumeric(3.14), expr: groupExpr(piExpr())()},
		{val: ValueString("str"), expr: groupExpr(strExpr())()},
		{val: ValueBoolean(true), expr: groupExpr(trueExpr())()},
		{val: ValueBoolean(false), expr: groupExpr(falseExpr())()},
		{val: ValueNil{}, expr: groupExpr(nilExpr())()},
		// complex
		// (1 + (3.14 / (-1) - 1)) + (+3.14)
		{
			val:  ValueNumeric(2.57),
			expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(groupExpr(bDivExpr(piExpr())(bSubExpr(groupExpr(uSubExpr(oneExpr())())())(oneExpr())())())())())())(groupExpr(uAddExpr(piExpr())())())(),
		},
		// (1 + (3.14 / (-1) - "str")) + (+3.14)
		{
			expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(groupExpr(bDivExpr(piExpr())(bSubExpr(groupExpr(uSubExpr(oneExpr())())())(strExpr())())())())())())(groupExpr(uAddExpr(piExpr())())())(),
			err:  NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpSubtract, TypeNumeric, TypeString), Position{}),
		},
		// "str" + ("str" + ("str" + "str"))
		{
			val:  ValueString("strstrstrstr"),
			expr: bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(strExpr())())())())())(),
		},
		// "str" + ("str" + (1 + "str"))
		{
			expr: bAddExpr(strExpr())(groupExpr(bAddExpr(strExpr())(groupExpr(bAddExpr(oneExpr())(strExpr())())())())())(),
			err:  NewRuntimeError(NewInvalidBinaryOperatorForTypeError(OpAdd, TypeNumeric, TypeString), Position{}),
		},
		// bAnd
		{val: ValueNumeric(3.14), expr: bAndExpr(oneExpr())(piExpr())()},
		{val: ValueBoolean(false), expr: bAndExpr(falseExpr())(oneExpr())()},
		{val: ValueString("str"), expr: bAndExpr(trueExpr())(strExpr())()},
		{val: ValueNil{}, expr: bAndExpr(nilExpr())(falseExpr())()},
		// bAnd
		{val: ValueNumeric(1), expr: bOrExpr(oneExpr())(piExpr())()},
		{val: ValueNumeric(1), expr: bOrExpr(falseExpr())(oneExpr())()},
		{val: ValueBoolean(true), expr: bOrExpr(trueExpr())(strExpr())()},
		{val: ValueBoolean(false), expr: bOrExpr(nilExpr())(falseExpr())()},
	}
	for _, test := range tests {
		x := NewExecutor(&PrintSpy{})
		// _, err := test.expr.TypeCheck(ctx)
		// if err != nil {
		// 	t.Errorf("Unexpected error while typechecking %q: %s", test.expr, err)
		// 	continue
		// }
		val, err := test.expr.Evaluate(x.ctx)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected evaluate(%v) to yield error %q, but got %q", test.expr, test.err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("Unexpected error while evaluating %q: %s", test.expr, err)
			continue
		}
		if val.Type() != test.val.Type() {
			t.Errorf("Expected evaluate(%q) yield value of type %s, but got %s", test.expr, test.val.Type(), val.Type())
			continue
		}
		if !test.val.Equals(val) {
			t.Errorf("Expected evaluate(%q) yield value %s, but got %s", test.expr, test.val, val)
		}
	}
}
