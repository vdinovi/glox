package lox

import (
	"fmt"
	"math"
	"strconv"
)

var Zero = ValueNumeric(0)

var True = ValueBoolean(true)

var False = ValueBoolean(false)

var Nil = ValueNil(struct{}{})

type Value interface {
	fmt.Stringer
	Unwrap() any
	Truthy() bool
	Type() Type
	Equals(Value) bool
}

type ValueString string

func (v ValueString) String() string {
	//return fmt.Sprintf("\"%s\"", string(v))
	return string(v)
}

func (ValueString) Type() Type {
	return TypeString
}

func (v ValueString) Unwrap() any {
	return string(v)
}

func (v ValueString) Truthy() bool {
	return true
}

func (v ValueString) Equals(other Value) bool {
	var ok bool
	var str ValueString
	if str, ok = other.(ValueString); !ok {
		return false
	}
	return string(v) == string(str)
}

func (v ValueString) Concat(other Value) (ValueString, error) {
	var ok bool
	var str ValueString
	if str, ok = other.(ValueString); !ok {
		return v, NewValueError(fmt.Sprintf("can't concat %s and %s", v.String(), other.String()))
	}
	return ValueString(string(v) + string(str)), nil
}

type ValueNumeric float64

func (v ValueNumeric) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

func (e ValueNumeric) Type() Type {
	return TypeNumeric
}

func (v ValueNumeric) Unwrap() any {
	return float64(v)
}

func (v ValueNumeric) Truthy() bool {
	return true
}

func (v ValueNumeric) Equals(other Value) bool {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return false
	}
	return v.approxEqual(num, 1e-9)
}

func (v ValueNumeric) approxEqual(other ValueNumeric, err float64) bool {
	return math.Abs(float64(v)-float64(other)) <= err
}

func (v ValueNumeric) Negative() (ValueNumeric, error) {
	return ValueNumeric(-float64(v)), nil
}

func (v ValueNumeric) Add(other Value) (ValueNumeric, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return v, NewValueError(fmt.Sprintf("can't add %s and %s", v.String(), other.String()))
	}
	return ValueNumeric(float64(v) + float64(num)), nil
}

func (v ValueNumeric) Subtract(other Value) (ValueNumeric, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return v, NewValueError(fmt.Sprintf("can't subtract %s and %s", v.String(), other.String()))
	}
	return ValueNumeric(float64(v) - float64(num)), nil
}

func (v ValueNumeric) Multiply(other Value) (ValueNumeric, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return v, NewValueError(fmt.Sprintf("can't multiply %s and %s", v.String(), other.String()))
	}
	return ValueNumeric(float64(v) * float64(num)), nil
}

func (v ValueNumeric) Divide(other Value) (ValueNumeric, error) {
	var ok bool
	var denom ValueNumeric
	if denom, ok = other.(ValueNumeric); !ok {
		return v, NewValueError(fmt.Sprintf("can't divide %s and %s", v.String(), other.String()))
	}
	n := float64(v)
	d := float64(denom)
	if d == 0 {
		if n == 0 {
			return ValueNumeric(0), nil
		}
		return v, NewDivideByZeroError(v, denom)
	}
	return ValueNumeric(n / d), nil
}

type ValueBoolean bool

func (v ValueBoolean) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (e ValueBoolean) Type() Type {
	return TypeBoolean
}

func (v ValueBoolean) Unwrap() any {
	return bool(v)
}

func (v ValueBoolean) Truthy() bool {
	return bool(v)
}

func (v ValueBoolean) Equals(other Value) bool {
	var ok bool
	var b ValueBoolean
	if b, ok = other.(ValueBoolean); !ok {
		return false
	}
	return bool(v) == bool(b)
}

type ValueNil struct{}

func (v ValueNil) String() string {
	return "nil"
}

func (e ValueNil) Type() Type {
	return TypeNil
}

func (v ValueNil) Unwrap() any {
	return struct{}{}
}

func (v ValueNil) Truthy() bool {
	return false
}

func (v ValueNil) Equals(other Value) bool {
	_, ok := other.(ValueNil)
	return ok
}
