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
		{ValueString(""), TypeString, true, "Value(\"\")"},
		{ValueString("str"), TypeString, true, "Value(\"str\")"},
		{ValueNumeric(0), TypeNumeric, true, "Value(0.000)"},
		{ValueNumeric(1), TypeNumeric, true, "Value(1.000)"},
		{ValueNumeric(-1), TypeNumeric, true, "Value(-1.000)"},
		{ValueNumeric(1.23), TypeNumeric, true, "Value(1.230)"},
		{ValueNumeric(-1.23), TypeNumeric, true, "Value(-1.230)"},
		{ValueBoolean(false), TypeBoolean, false, "Value(false)"},
		{ValueBoolean(true), TypeBoolean, true, "Value(true)"},
		{ValueNil(struct{}{}), TypeNil, false, "Value(nil)"},
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
