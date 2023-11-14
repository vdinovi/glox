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
