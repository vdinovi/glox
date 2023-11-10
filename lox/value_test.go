package lox

import (
	"testing"
)

func TestValue(t *testing.T) {
	symbols := make(Symbols)
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
		if actual, err := test.val.Type(symbols); err != nil {
			t.Errorf("Unexpected error: %s", err)
			continue
		} else if actual != test.typ {
			t.Errorf("Expected %s to have type %s, but got %s", test.val, test.typ, actual)
		}

		if actual := test.val.Truthy(); actual != test.truthy {
			if test.truthy {
				t.Errorf("Expected %s to be truthy, but was falsey", test.val)
			} else {
				t.Errorf("Expected %s to be falsey, but was truthy", test.val)
			}
		}

		if actual := test.val.String(); actual != test.str {
			t.Errorf("Expected %s to have string %s, but got %s", test.val, test.str, actual)
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
