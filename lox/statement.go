package lox

import (
	"fmt"
)

type Statement interface {
	TypeCheck() error
	Execute(*Executor) error
	fmt.Stringer
}

type ExpressionStatement struct {
	expr Expression
}

func (s ExpressionStatement) TypeCheck() error {
	return TypeCheckExpressionStatement(s)
}

func (s ExpressionStatement) Execute(e *Executor) error {
	return e.ExecuteExpressionStatement(s)
}

func (s ExpressionStatement) String() string {
	return fmt.Sprintf("%s ;", s.expr)
}

type PrintStatement struct {
	expr Expression
}

func (s PrintStatement) TypeCheck() error {
	return TypeCheckPrintStatement(s)
}

func (s PrintStatement) Execute(e *Executor) error {
	return e.ExecutePrintStatement(s)
}

func (s PrintStatement) String() string {
	return fmt.Sprintf("print %s ;", s.expr)
}
