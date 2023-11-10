package lox

import (
	"fmt"
)

type Evaluable interface {
	Evaluate() (Value, error)
}

type Expression interface {
	Typed
	Evaluable
	fmt.Stringer
}

type UnaryExpression struct {
	op    Operator
	right Expression
}

func (e UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
}

func (e UnaryExpression) Type() (Type, error) {
	return TypeCheckUnaryExpression(e)
}

func (e UnaryExpression) Evaluate() (Value, error) {
	return EvaluateUnaryExpression(e)
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

func (e BinaryExpression) Type() (Type, error) {
	return TypeCheckBinaryExpression(e)
}

func (e BinaryExpression) Evaluate() (Value, error) {
	return EvaluateBinaryExpression(e)
}

type GroupingExpression struct {
	expr Expression
}

func (e GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e GroupingExpression) Type() (Type, error) {
	return TypeCheckGroupingExpression(e)
}

func (e GroupingExpression) Evaluate() (Value, error) {
	return EvaluateGroupingExpression(e)
}

type LiteralExpression interface {
	literal()
}

type StringExpression string

func (e StringExpression) literal() {}

func (e StringExpression) String() string {
	return string(e)
}

func (e StringExpression) Type() (Type, error) {
	return TypeCheckStringExpression(e)
}

func (e StringExpression) Evaluate() (Value, error) {
	return EvaluateStringExpression(e)
}

type NumericExpression float64

func (e NumericExpression) literal() {}

func (e NumericExpression) String() string {
	return fmt.Sprint(float64(e))
}

func (e NumericExpression) Type() (Type, error) {
	return TypeCheckNumericExpression(e)
}

func (e NumericExpression) Evaluate() (Value, error) {
	return EvaluateNumericExpression(e)
}

type BooleanExpression bool

func (e BooleanExpression) literal() {}

func (e BooleanExpression) String() string {
	if bool(e) {
		return "(true)"
	}
	return "(false)"
}

func (e BooleanExpression) Type() (Type, error) {
	return TypeCheckBooleanExpression(e)
}
func (e BooleanExpression) Evaluate() (Value, error) {
	return EvaluateBooleanExpression(e)
}

type NilExpression struct{}

func (e NilExpression) literal() {}

func (e NilExpression) String() string {
	return "(nil)"
}

func (e NilExpression) Type() (Type, error) {
	return TypeCheckNilExpression(e)
}
func (e NilExpression) Evaluate() (Value, error) {
	return EvaluateNilExpression(e)
}
