package lox

import (
	"errors"
	"fmt"
	"math"
)

type Value interface {
	fmt.Stringer
	Printable
	Truthy() bool
	Type() Type
	Equals(Value) bool
}

var Zero = ValueNumeric(0)

var True = ValueBoolean(true)

var False = ValueBoolean(false)

var Nil = ValueNil(struct{}{})

var ErrInvalidType = errors.New("invalid type")

type ValueString string

func (v ValueString) String() string {
	str, err := v.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (ValueString) Type() Type {
	return TypeString
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
		return v, ErrInvalidType
	}
	return ValueString(string(v) + string(str)), nil
}

type ValueNumeric float64

func (v ValueNumeric) String() string {
	str, err := v.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e ValueNumeric) Type() Type {
	return TypeNumeric
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
		return v, ErrInvalidType
	}
	return ValueNumeric(float64(v) + float64(num)), nil
}

func (v ValueNumeric) Subtract(other Value) (ValueNumeric, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return v, ErrInvalidType
	}
	return ValueNumeric(float64(v) - float64(num)), nil
}

func (v ValueNumeric) Multiply(other Value) (ValueNumeric, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return v, ErrInvalidType
	}
	return ValueNumeric(float64(v) * float64(num)), nil
}

func (v ValueNumeric) Divide(other Value) (ValueNumeric, error) {
	var ok bool
	var denom ValueNumeric
	if denom, ok = other.(ValueNumeric); !ok {
		return v, NewInvalidBinaryOperatorForTypeError(OpDivide, v.Type(), other.Type())
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

func (v ValueNumeric) Compare(other Value) (int, error) {
	var ok bool
	var num ValueNumeric
	if num, ok = other.(ValueNumeric); !ok {
		return 0, ErrInvalidType
	}
	a := float64(v)
	b := float64(num)
	if a == b {
		return 0, nil
	}
	if a < b {
		return -1, nil
	}
	return 1, nil
}

type ValueBoolean bool

func (v ValueBoolean) String() string {
	str, err := v.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e ValueBoolean) Type() Type {
	return TypeBoolean
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
	str, err := v.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e ValueNil) Type() Type {
	return TypeNil
}

func (v ValueNil) Truthy() bool {
	return false
}

func (v ValueNil) Equals(other Value) bool {
	_, ok := other.(ValueNil)
	return ok
}

type ValueCallable struct {
	name string
	fn   Function
}

func (v ValueCallable) String() string {
	str, err := v.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e ValueCallable) Type() Type {
	return TypeCallable
}

func (v ValueCallable) Truthy() bool {
	return true
}

func (v ValueCallable) Equals(other Value) bool {
	call, ok := other.(ValueCallable)
	if !ok || v.name != call.name {
		return false
	}
	return true
}

func (v ValueCallable) Call(ctx *Context, args ...Value) (Value, error) {
	return v.fn.Execute(ctx, v.name, args...)
}
