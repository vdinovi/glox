package lox

import (
	"fmt"
)

type Expression interface {
	Position() Position
	Type() Type
	Equals(Expression) bool
	TypeCheck(*Context) (Type, error)
	Evaluate(*Context) (Value, error)
	fmt.Stringer
}

type UnaryExpression struct {
	op    Operator
	right Expression
	pos   Position
	typ   Type
}

func (e *UnaryExpression) Position() Position {
	return e.pos
}

func (e *UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
}

func (e *UnaryExpression) Type() Type {
	return e.typ
}

func (e *UnaryExpression) TypeCheck(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckUnaryExpression(e)
	if err != nil {
		return ErrType, err
	}
	e.typ = typ
	return typ, nil
}

func (e *UnaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateUnaryExpression(e)
	return val, err
}

func (e *UnaryExpression) Equals(other Expression) bool {
	unary, ok := other.(*UnaryExpression)
	if !ok || e.op != unary.op {
		return false
	}
	return e.right.Equals(unary.right)
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
	pos   Position
	typ   Type
}

func (e *BinaryExpression) Position() Position {
	return e.pos
}

func (e *BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

func (e *BinaryExpression) Type() Type {
	return e.typ
}

func (e *BinaryExpression) TypeCheck(ctx *Context) (Type, error) {
	_, _, typ, err := ctx.TypeCheckBinaryExpression(e)
	if err != nil {
		return ErrType, err
	}
	e.typ = typ
	return typ, nil
}

func (e *BinaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateBinaryExpression(e)
	return val, err
}

func (e *BinaryExpression) Equals(other Expression) bool {
	binary, ok := other.(*BinaryExpression)
	if !ok || e.op != binary.op {
		return false
	}
	return e.left.Equals(binary.left) && e.right.Equals(binary.right)
}

type GroupingExpression struct {
	expr Expression
	pos  Position
	typ  Type
}

func (e *GroupingExpression) Position() Position {
	return e.pos
}

func (e *GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e *GroupingExpression) Type() Type {
	return e.typ
}

func (e *GroupingExpression) TypeCheck(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckGroupingExpression(e)
	if err != nil {
		return ErrType, err
	}
	e.typ = typ
	return typ, nil
}

func (e *GroupingExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateGroupingExpression(e)
	return val, err
}

func (e *GroupingExpression) Equals(other Expression) bool {
	group, ok := other.(*GroupingExpression)
	if !ok {
		return false
	}
	return e.expr.Equals(group.expr)
}

type AssignmentExpression struct {
	name  string
	right Expression
	pos   Position
	typ   Type
}

func (e *AssignmentExpression) Position() Position {
	return e.pos
}

func (e *AssignmentExpression) String() string {
	return fmt.Sprintf("(%s = %s)", e.name, e.right)
}

func (e *AssignmentExpression) Type() Type {
	return e.typ
}

func (e *AssignmentExpression) TypeCheck(ctx *Context) (Type, error) {
	_, typ, err := ctx.TypeCheckAssignmentExpression(e)
	if err != nil {
		return ErrType, err
	}
	e.typ = typ
	return typ, nil
}

func (e *AssignmentExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateAssignmentExpression(e)
	return val, err
}

func (e *AssignmentExpression) Equals(other Expression) bool {
	assign, ok := other.(*AssignmentExpression)
	if !ok || e.name != assign.name {
		return false
	}
	return e.right.Equals(assign.right)
}

type VariableExpression struct {
	name string
	pos  Position
	typ  Type
}

func (e *VariableExpression) Position() Position {
	return e.pos
}

func (e *VariableExpression) String() string {
	return e.name
}

func (e *VariableExpression) Type() Type {
	return e.typ
}

func (e *VariableExpression) TypeCheck(ctx *Context) (Type, error) {
	typ, err := ctx.TypeCheckVariableExpression(e)
	if err != nil {
		return ErrType, err
	}
	e.typ = typ
	return typ, nil
}

func (e *VariableExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateVariableExpression(e)
	return val, err
}

func (e *VariableExpression) Equals(other Expression) bool {
	variable, ok := other.(*VariableExpression)
	return ok && e.name == variable.name
}

type StringExpression struct {
	value string
	pos   Position
}

func (e *StringExpression) Position() Position {
	return e.pos
}

func (e *StringExpression) String() string {
	return e.value
}

func (e *StringExpression) Type() Type {
	return TypeString
}

func (e *StringExpression) TypeCheck(*Context) (Type, error) {
	return e.Type(), nil
}

func (e *StringExpression) Evaluate(ctx *Context) (Value, error) {
	return ValueString(e.value), nil
}

func (e *StringExpression) Equals(other Expression) bool {
	str, ok := other.(*StringExpression)
	return ok && e.value == str.value
}

type NumericExpression struct {
	value float64
	pos   Position
}

func (e *NumericExpression) Position() Position {
	return e.pos
}

func (e *NumericExpression) String() string {
	return fmt.Sprint(e.value)
}

func (e *NumericExpression) Type() Type {
	return TypeNumeric
}

func (e *NumericExpression) TypeCheck(*Context) (Type, error) {
	return e.Type(), nil
}

func (e *NumericExpression) Evaluate(ctx *Context) (Value, error) {
	return ValueNumeric(e.value), nil
}

func (e *NumericExpression) Equals(other Expression) bool {
	num, ok := other.(*NumericExpression)
	return ok && e.value == num.value
}

type BooleanExpression struct {
	value bool
	pos   Position
}

func (e *BooleanExpression) Position() Position {
	return e.pos
}

func (e *BooleanExpression) String() string {
	if e.value {
		return "true"
	}
	return "false"
}

func (e *BooleanExpression) Type() Type {
	return TypeBoolean
}

func (e *BooleanExpression) TypeCheck(*Context) (Type, error) {
	return e.Type(), nil
}

func (e *BooleanExpression) Evaluate(ctx *Context) (Value, error) {
	return ValueBoolean(e.value), nil
}

func (e *BooleanExpression) Equals(other Expression) bool {
	boolean, ok := other.(*BooleanExpression)
	return ok && e.value == boolean.value
}

type NilExpression struct {
	pos Position
}

func (e *NilExpression) Position() Position {
	return e.pos
}

func (e *NilExpression) String() string {
	return "nil"
}

func (e *NilExpression) Type() Type {
	return TypeNil
}

func (e *NilExpression) TypeCheck(*Context) (Type, error) {
	return e.Type(), nil
}

func (e *NilExpression) Evaluate(ctx *Context) (Value, error) {
	return Nil, nil
}

func (e *NilExpression) Equals(other Expression) bool {
	_, ok := other.(*NilExpression)
	return ok
}
