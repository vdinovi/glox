package lox

import (
	"fmt"
)

type Program []Statement

type Statement interface {
	TypeCheck(Symbols) error
	Execute(*Executor) error
	fmt.Stringer
}

type ExpressionStatement struct {
	expr Expression
}

func (s ExpressionStatement) TypeCheck(syms Symbols) error {
	return syms.TypeCheckExpressionStatement(s)
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

func (s PrintStatement) TypeCheck(syms Symbols) error {
	return syms.TypeCheckPrintStatement(s)
}

func (s PrintStatement) Execute(e *Executor) error {
	return e.ExecutePrintStatement(s)
}

func (s PrintStatement) String() string {
	return fmt.Sprintf("print %s ;", s.expr)
}

type DeclarationStatement struct {
	name string
	expr Expression
}

func (d DeclarationStatement) TypeCheck(Symbols) error {
	return nil
}

func (d DeclarationStatement) Execute(*Executor) error {
	return nil
}

func (d DeclarationStatement) String() string {
	return ""
}
