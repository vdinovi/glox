package lox

import (
	"fmt"
)

type Evaluable interface {
	Evaluate(*EvaluationContext) (Value, error)
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

func (e UnaryExpression) Type(syms Symbols) (Type, error) {
	_, typ, err := syms.TypeCheckUnaryExpression(e)
	return typ, err
}

func (e UnaryExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	val, _, err := ctx.EvaluateUnaryExpression(e)
	return val, err
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

func (e BinaryExpression) Type(syms Symbols) (Type, error) {
	_, _, typ, err := syms.TypeCheckBinaryExpression(e)
	return typ, err
}

func (e BinaryExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	val, _, err := ctx.EvaluateBinaryExpression(e)
	return val, err
}

type GroupingExpression struct {
	expr Expression
}

func (e GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e GroupingExpression) Type(syms Symbols) (Type, error) {
	_, typ, err := syms.TypeCheckGroupingExpression(e)
	return typ, err
}

func (e GroupingExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	val, _, err := ctx.EvaluateGroupingExpression(e)
	return val, err
}

type LiteralExpression interface {
	literal()
}

type StringExpression string

func (e StringExpression) literal() {}

func (e StringExpression) String() string {
	return string(e)
}

func (e StringExpression) Type(syms Symbols) (Type, error) {
	return syms.TypeCheckStringExpression(e)
}

func (e StringExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	return ctx.EvaluateStringExpression(e)
}

type NumericExpression float64

func (e NumericExpression) literal() {}

func (e NumericExpression) String() string {
	return fmt.Sprint(float64(e))
}

func (e NumericExpression) Type(syms Symbols) (Type, error) {
	return syms.TypeCheckNumericExpression(e)
}

func (e NumericExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	return ctx.EvaluateNumericExpression(e)
}

type BooleanExpression bool

func (e BooleanExpression) literal() {}

func (e BooleanExpression) String() string {
	if bool(e) {
		return "(true)"
	}
	return "(false)"
}

func (e BooleanExpression) Type(syms Symbols) (Type, error) {
	return syms.TypeCheckBooleanExpression(e)
}
func (e BooleanExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	return ctx.EvaluateBooleanExpression(e)
}

type NilExpression struct{}

func (e NilExpression) literal() {}

func (e NilExpression) String() string {
	return "(nil)"
}

func (e NilExpression) Type(syms Symbols) (Type, error) {
	return syms.TypeCheckNilExpression(e)
}
func (e NilExpression) Evaluate(ctx *EvaluationContext) (Value, error) {
	return ctx.EvaluateNilExpression(e)
}
