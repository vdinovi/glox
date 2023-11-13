package lox

import (
	"fmt"
)

type Program []Statement

type Statement interface {
	Position() Position
	TypeCheck(*EvaluationContext) error
	Execute(*Executor) error
	fmt.Stringer
}

type ExpressionStatement struct {
	expr Expression
	pos  Position
}

func (s ExpressionStatement) Position() Position {
	return s.pos
}

func (s ExpressionStatement) TypeCheck(ctx *EvaluationContext) error {
	return ctx.TypeCheckExpressionStatement(s)
}

func (s ExpressionStatement) Execute(e *Executor) error {
	return e.ExecuteExpressionStatement(s)
}

func (s ExpressionStatement) String() string {
	return fmt.Sprintf("%s ;", s.expr)
}

type PrintStatement struct {
	expr Expression
	pos  Position
}

func (s PrintStatement) Position() Position {
	return s.pos
}

func (s PrintStatement) TypeCheck(ctx *EvaluationContext) error {
	return ctx.TypeCheckPrintStatement(s)
}

func (s PrintStatement) Execute(e *Executor) error {
	return e.ExecutePrintStatement(s)
}

func (s PrintStatement) String() string {
	return fmt.Sprintf("print %s ;", s.expr)
}

type DeclarationStatement struct {
	name string
	pos  Position
	expr Expression
}

func (s DeclarationStatement) Position() Position {
	return s.pos
}

func (d DeclarationStatement) TypeCheck(*EvaluationContext) error {
	return nil
}

func (d DeclarationStatement) Execute(*Executor) error {
	return nil
}

func (d DeclarationStatement) String() string {
	return ""
}
