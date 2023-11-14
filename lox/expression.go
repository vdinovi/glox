package lox

import (
	"fmt"
)

type Evaluable interface {
	Evaluate(*Context) (Value, error)
}

type Expression interface {
	Position() Position
	Type(*Context) (Type, error)
	Evaluable
	fmt.Stringer
}

type UnaryExpression struct {
	op    Operator
	right Expression
	pos   Position
}

func (e UnaryExpression) Position() Position {
	return e.pos
}

func (e UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
}

func (e UnaryExpression) Type(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckUnaryExpression(e)
	return typ, err
}

func (e UnaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateUnaryExpression(e)
	return val, err
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
	pos   Position
}

func (e BinaryExpression) Position() Position {
	return e.pos
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

func (e BinaryExpression) Type(ctx *Context) (Type, error) {
	_, _, typ, err := ctx.TypeCheckBinaryExpression(e)
	return typ, err
}

func (e BinaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateBinaryExpression(e)
	return val, err
}

type GroupingExpression struct {
	expr Expression
	pos  Position
}

func (e GroupingExpression) Position() Position {
	return e.pos
}

func (e GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e GroupingExpression) Type(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckGroupingExpression(e)
	return typ, err
}

func (e GroupingExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateGroupingExpression(e)
	return val, err
}

type AssignmentExpression struct {
	name  string
	right Expression
	pos   Position
}

func (e AssignmentExpression) Position() Position {
	return e.pos
}

func (e AssignmentExpression) String() string {
	return fmt.Sprintf("(%s = %s)", e.name, e.right)
}

func (e AssignmentExpression) Type(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckAssignmentExpression(e)
	return typ, err
}

func (e AssignmentExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateAssignmentExpression(e)
	return val, err
}

type VariableExpression struct {
	name string
	pos  Position
}

func (e VariableExpression) Position() Position {
	return e.pos
}

func (e VariableExpression) String() string {
	return e.name
}

func (e VariableExpression) Type(ctx *Context) (Type, error) {
	return ctx.TypeCheckVariableExpression(e)
}

func (e VariableExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateVariableExpression(e)
	return val, err
}

type StringExpression struct {
	value string
	pos   Position
}

func (e StringExpression) Position() Position {
	return e.pos
}

func (e StringExpression) String() string {
	return e.value
}

func (e StringExpression) Type(ctx *Context) (Type, error) {
	return ctx.TypeCheckStringExpression(e)
}

func (e StringExpression) Evaluate(ctx *Context) (Value, error) {
	return ctx.EvaluateStringExpression(e)
}

type NumericExpression struct {
	value float64
	pos   Position
}

func (e NumericExpression) Position() Position {
	return e.pos
}

func (e NumericExpression) String() string {
	return fmt.Sprint(e.value)
}

func (e NumericExpression) Type(ctx *Context) (Type, error) {
	return ctx.TypeCheckNumericExpression(e)
}

func (e NumericExpression) Evaluate(ctx *Context) (Value, error) {
	return ctx.EvaluateNumericExpression(e)
}

type BooleanExpression struct {
	value bool
	pos   Position
}

func (e BooleanExpression) Position() Position {
	return e.pos
}

func (e BooleanExpression) String() string {
	if e.value {
		return "true"
	}
	return "false"
}

func (e BooleanExpression) Type(ctx *Context) (Type, error) {
	return ctx.TypeCheckBooleanExpression(e)
}
func (e BooleanExpression) Evaluate(ctx *Context) (Value, error) {
	return ctx.EvaluateBooleanExpression(e)
}

type NilExpression struct {
	pos Position
}

func (e NilExpression) Position() Position {
	return e.pos
}

func (e NilExpression) String() string {
	return "nil"
}

func (e NilExpression) Type(ctx *Context) (Type, error) {
	return ctx.TypeCheckNilExpression(e)
}
func (e NilExpression) Evaluate(ctx *Context) (Value, error) {
	return ctx.EvaluateNilExpression(e)
}
