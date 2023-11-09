package lox

import (
	"fmt"
)

type Expression interface {
	fmt.Stringer
}

type UnaryExpression struct {
	op    Operator
	right Expression
}

func (e UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

type GroupingExpression struct {
	expr Expression
}

func (e GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

type LiteralExpression interface {
	literal()
}

type StringExpression string

func (e StringExpression) String() string {
	return string(e)
}

func (e StringExpression) literal() {}

type NumberExpression float64

func (e NumberExpression) String() string {
	return fmt.Sprint(float64(e))
}

func (e NumberExpression) literal() {}

type BoolExpression bool

func (e BoolExpression) String() string {
	if bool(e) {
		return "true"
	} else {
		return "false"
	}
}

func (e BoolExpression) literal() {}

type NilExpression struct{}

func (e NilExpression) String() string {
	return "nil"
}

func (e NilExpression) literal() {}
