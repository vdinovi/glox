package lox

import (
	"fmt"
)

type Value interface {
	Typed
	fmt.Stringer
	Unwrap() any
	Truthy() bool
}

type ValueString string

func (v ValueString) String() string {
	return string(v)
}

func (ValueString) Type() (Type, error) {
	return TypeString, nil
}

func (v ValueString) Unwrap() any {
	return string(v)
}

func (v ValueString) Truthy() bool {
	return true
}

func TypeCheckStringExpression(e StringExpression) (Type, error) {
	return TypeString, nil
}

func EvaluateStringExpression(e StringExpression) (Value, error) {
	return ValueString(e), nil
}

type ValueNumeric float64

func (v ValueNumeric) String() string {
	return fmt.Sprint(float64(v))
}
func (e ValueNumeric) Type() (Type, error) {
	return TypeNumeric, nil
}

func (v ValueNumeric) Unwrap() any {
	return float64(v)
}

func (v ValueNumeric) Truthy() bool {
	return true
}

func TypeCheckNumericExpression(NumericExpression) (Type, error) {
	return TypeNumeric, nil
}

func EvaluateNumericExpression(e NumericExpression) (Value, error) {
	return ValueNumeric(e), nil
}

type ValueBoolean bool

var True = ValueBoolean(true)

var False = ValueBoolean(false)

func (v ValueBoolean) String() string {
	if bool(v) {
		return "true"
	}
	return "false"
}

func (e ValueBoolean) Type() (Type, error) {
	return TypeBoolean, nil
}

func (v ValueBoolean) Unwrap() any {
	return bool(v)
}

func (v ValueBoolean) Truthy() bool {
	return bool(v)
}

func TypeCheckBooleanExpression(BooleanExpression) (Type, error) {
	return TypeBoolean, nil
}

func EvaluateBooleanExpression(e BooleanExpression) (Value, error) {
	return ValueBoolean(e), nil
}

type ValueNil struct{}

var Nil = ValueNil(struct{}{})

func (v ValueNil) String() string {
	return "nil"
}

func (e ValueNil) Type() (Type, error) {
	return TypeNil, nil
}

func (v ValueNil) Unwrap() any {
	return nil
}

func (v ValueNil) Truthy() bool {
	return false
}

func TypeCheckNilExpression(e NilExpression) (Type, error) {
	return TypeNil, nil
}

func EvaluateNilExpression(e NilExpression) (Value, error) {
	return Nil, nil
}
