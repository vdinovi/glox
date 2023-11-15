package lox

import (
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		val    Value
		typ    Type
		truthy bool
		str    string
	}{
		{ValueString(""), TypeString, true, "\"\""},
		{ValueString("str"), TypeString, true, "\"str\""},
		{ValueNumeric(0), TypeNumeric, true, "0"},
		{ValueNumeric(1), TypeNumeric, true, "1"},
		{ValueNumeric(-1), TypeNumeric, true, "-1"},
		{ValueNumeric(1.23), TypeNumeric, true, "1.23"},
		{ValueNumeric(-1.23), TypeNumeric, true, "-1.23"},
		{ValueBoolean(false), TypeBoolean, false, "false"},
		{ValueBoolean(true), TypeBoolean, true, "true"},
		{ValueNil(struct{}{}), TypeNil, false, "nil"},
	}
	for _, test := range tests {
		// x := NewExecutor(io.Discard)
		if v := test.val.Type(); v != test.typ {
			t.Errorf("Expected %s to have type %s, but got %s", test.val, test.typ, v)
		}
		if b := test.val.Truthy(); b != test.truthy {
			t.Errorf("Expected %s to be %v, but was %v", test.val, test.truthy, b)
		}
		if s := test.val.String(); s != test.str {
			t.Errorf("Expected %s to have string %s, but got %s", test.val, test.str, s)
		}
	}
}

func TestValueUnwrap(t *testing.T) {
	var val Value

	val = ValueString("str")
	if s, ok := val.Unwrap().(string); !ok {
		t.Errorf("Failed to downcast %s to string", val)
	} else if s != "str" {
		t.Errorf("Expected %s to downcast to %v, but got %v", val, "str", s)
	}

	val = ValueNumeric(1.23)
	if n, ok := val.Unwrap().(float64); !ok {
		t.Errorf("Failed to downcast %s to float64", val)
	} else if n != 1.23 {
		t.Errorf("Expected %s to downcast to %v, but got %v", val, 1.23, n)
	}

	val = ValueBoolean(true)
	if n, ok := val.Unwrap().(bool); !ok {
		t.Errorf("Failed to downcast %s to bool", val)
	} else if n != true {
		t.Errorf("Expected %s to downcast to %v, but got %v", val, true, n)
	}

	val = ValueNil(struct{}{})
	if n, ok := val.Unwrap().(struct{}); !ok {
		t.Errorf("Failed to downcast %s to struct{}", val)
	} else if n != struct{}{} {
		t.Errorf("Expected %s to downcast to %v, but got %v", val, struct{}{}, n)
	}
}

func TestValueEquals(t *testing.T) {
	tests := []struct {
		eq bool
		a  Value
		b  Value
	}{
		// numeric
		{eq: true, a: ValueNumeric(0), b: ValueNumeric(0)},
		{eq: true, a: ValueNumeric(1), b: ValueNumeric(1)},
		{eq: true, a: ValueNumeric(-1), b: ValueNumeric(-1)},
		{eq: true, a: ValueNumeric(3.14), b: ValueNumeric(3.14)},
		{eq: true, a: ValueNumeric(-3.14), b: ValueNumeric(-3.14)},
		{eq: false, a: ValueNumeric(1), b: ValueNumeric(0)},
		{eq: false, a: ValueNumeric(1), b: ValueNumeric(3.14)},
		{eq: false, a: ValueNumeric(1), b: ValueString("")},
		{eq: false, a: ValueNumeric(1), b: ValueBoolean(true)},
		{eq: false, a: ValueNumeric(1), b: ValueBoolean(false)},
		{eq: false, a: ValueNumeric(1), b: ValueNil{}},
		// string
		{eq: true, a: ValueString(""), b: ValueString("")},
		{eq: true, a: ValueString("str"), b: ValueString("str")},
		{eq: false, a: ValueString("str"), b: ValueString("")},
		{eq: false, a: ValueString("str"), b: ValueNumeric(1)},
		{eq: false, a: ValueString("str"), b: ValueBoolean(true)},
		{eq: false, a: ValueString("str"), b: ValueBoolean(false)},
		{eq: false, a: ValueString("str"), b: ValueNil{}},
		// boolean
		{eq: true, a: ValueBoolean(true), b: ValueBoolean(true)},
		{eq: true, a: ValueBoolean(false), b: ValueBoolean(false)},
		{eq: false, a: ValueBoolean(true), b: ValueBoolean(false)},
		{eq: false, a: ValueBoolean(true), b: ValueNumeric(1)},
		{eq: false, a: ValueBoolean(true), b: ValueString("")},
		{eq: false, a: ValueBoolean(true), b: ValueNil{}},
		// nil
		{eq: true, a: ValueNil{}, b: ValueNil{}},
		{eq: false, a: ValueNil{}, b: ValueNumeric(1)},
		{eq: false, a: ValueNil{}, b: ValueString("")},
		{eq: false, a: ValueNil{}, b: ValueBoolean(true)},
		{eq: false, a: ValueNil{}, b: ValueBoolean(false)},
	}
	for _, test := range tests {
		if test.eq {
			if !test.a.Equals(test.b) {
				t.Errorf("Expected %s to equal %s but does not", test.a, test.b)
			}
		} else {
			if test.a.Equals(test.b) {
				t.Errorf("Expected %s to not equal %s but does", test.a, test.b)
			}
		}
	}
}

func TestValueSentinels(t *testing.T) {
	tests := []struct {
		val  Value
		want Value
	}{
		{Zero, ValueNumeric(0)},
		{True, ValueBoolean(true)},
		{False, ValueBoolean(false)},
		{Nil, ValueNil(struct{}{})},
	}
	for _, test := range tests {
		if test.val != test.want {
			t.Errorf("Expected %v to equal %v", test.val, test.want)
		}
	}
}

func TestValueNumericBinaryOps(t *testing.T) {
	tests := []struct {
		op  string
		a   float64
		b   float64
		val ValueNumeric
		err error
	}{
		{op: "Add", a: 0, b: 1, val: ValueNumeric(1)},
		{op: "Add", a: 1, b: 1, val: ValueNumeric(2)},
		{op: "Add", a: 0, b: -1, val: ValueNumeric(-1)},
		{op: "Add", a: -1, b: -1, val: ValueNumeric(-2)},
		{op: "Add", a: 0, b: 3.14, val: ValueNumeric(3.14)},
		{op: "Add", a: 3.14, b: 3.14, val: ValueNumeric(6.28)},
		{op: "Add", a: 0, b: -3.14, val: ValueNumeric(-3.14)},
		{op: "Add", a: -3.14, b: -3.14, val: ValueNumeric(-6.28)},
		{op: "Subtract", a: 0, b: 1, val: ValueNumeric(-1)},
		{op: "Subtract", a: 1, b: 1, val: ValueNumeric(0)},
		{op: "Subtract", a: 0, b: -1, val: ValueNumeric(1)},
		{op: "Subtract", a: -1, b: -1, val: ValueNumeric(0)},
		{op: "Subtract", a: 0, b: 3.14, val: ValueNumeric(-3.14)},
		{op: "Subtract", a: 3.14, b: 3.14, val: ValueNumeric(0)},
		{op: "Subtract", a: 0, b: -3.14, val: ValueNumeric(3.14)},
		{op: "Subtract", a: -3.14, b: -3.14, val: ValueNumeric(0)},
		{op: "Multiply", a: 0, b: 1, val: ValueNumeric(0)},
		{op: "Multiply", a: 1, b: 1, val: ValueNumeric(1)},
		{op: "Multiply", a: 0, b: -1, val: ValueNumeric(0)},
		{op: "Multiply", a: -1, b: -1, val: ValueNumeric(1)},
		{op: "Multiply", a: 0, b: 3.14, val: ValueNumeric(0)},
		{op: "Multiply", a: 3.14, b: 3.14, val: ValueNumeric(9.8596)},
		{op: "Multiply", a: 0, b: -3.14, val: ValueNumeric(0)},
		{op: "Multiply", a: -3.14, b: -3.14, val: ValueNumeric(9.8596)},
		{op: "Divide", a: 0, b: 1, val: ValueNumeric(0)},
		{op: "Divide", a: 1, b: 1, val: ValueNumeric(1)},
		{op: "Divide", a: 0, b: -1, val: ValueNumeric(0)},
		{op: "Divide", a: -1, b: -1, val: ValueNumeric(1)},
		{op: "Divide", a: 0, b: 3.14, val: ValueNumeric(0)},
		{op: "Divide", a: 3.14, b: 3.14, val: ValueNumeric(1)},
		{op: "Divide", a: 0, b: -3.14, val: ValueNumeric(0)},
		{op: "Divide", a: -3.14, b: -3.14, val: ValueNumeric(1)},
		{op: "Divide", a: 0, b: 0, val: ValueNumeric(0)},
		{op: "Divide", a: 1, b: 0, err: NewDivideByZeroError(ValueNumeric(1), ValueNumeric(0))},
		{op: "Divide", a: 3.14, b: 0, err: NewDivideByZeroError(ValueNumeric(3.14), ValueNumeric(0))},
		{op: "Divide", a: -3.14, b: 0, err: NewDivideByZeroError(ValueNumeric(-3.14), ValueNumeric(0))},
	}
	for _, test := range tests {
		t.Log(test.op, test.a, test.b)
		var err error
		var val ValueNumeric
		a := ValueNumeric(test.a)
		b := ValueNumeric(test.b)
		switch test.op {
		case "Add":
			val, err = a.Add(b)
		case "Subtract":
			val, err = a.Subtract(b)
		case "Multiply":
			val, err = a.Multiply(b)
		case "Divide":
			val, err = a.Divide(b)
		default:
			t.Errorf("Unexpected operation %s", test.op)
			continue
		}
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected (%s).%s(%s) to yield error %q, but got %q", a, test.op, b, test.err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("Unexpected error in (%s).%s(%s): %s", a, test.op, b, err)
			continue
		}
		if val != test.val {
			t.Errorf("Expected (%s).%s(%s) to yield value %s, but got %s", a, test.op, b, test.val, val)
		}
	}
}

func TestValueStringBinaryOps(t *testing.T) {
	tests := []struct {
		op  string
		a   string
		b   string
		val ValueString
		err error
	}{
		{op: "Concat", a: "", b: "", val: ValueString("")},
		{op: "Concat", a: "x", b: "", val: ValueString("x")},
		{op: "Concat", a: "x", b: "y", val: ValueString("xy")},
	}
	for _, test := range tests {
		t.Log(test.op, test.a, test.b)
		var err error
		var val ValueString
		a := ValueString(test.a)
		b := ValueString(test.b)
		switch test.op {
		case "Concat":
			val, err = a.Concat(b)
		default:
			t.Errorf("Unexpected operation %s", test.op)
			continue
		}
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected (%s).%s(%s) to yield error %q, but got %q", a, test.op, b, test.err, err)
			}
			continue
		} else if err != nil {
			t.Errorf("Unexpected error in (%s).%s(%s): %s", a, test.op, b, err)
			continue
		}
		if val != test.val {
			t.Errorf("Expected (%s).%s(%s) to yield value %s, but got %s", a, test.op, b, test.val, val)
		}
	}
}
