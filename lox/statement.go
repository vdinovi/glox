package lox

import (
	"fmt"
	"strings"
)

type Program []Statement

type Statement interface {
	Position() Position
	TypeCheck(ctx *Context) error
	Execute(*Executor) error
	fmt.Stringer
}

type BlockStatement struct {
	stmts []Statement
	pos   Position
}

func (s BlockStatement) Position() Position {
	return s.pos
}

func (s BlockStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckBlockStatement(s)
}

func (s BlockStatement) Execute(e *Executor) error {
	return e.ExecuteBlockStatement(s)
}

func (s BlockStatement) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.stmts {
		fmt.Fprint(&sb, " ", stmt.String())
	}
	sb.WriteString(" } ;")
	return sb.String()
}

type ExpressionStatement struct {
	expr Expression
	pos  Position
}

func (s ExpressionStatement) Position() Position {
	return s.pos
}

func (s ExpressionStatement) TypeCheck(ctx *Context) error {
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

func (s PrintStatement) TypeCheck(ctx *Context) error {
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

func (s DeclarationStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckDeclarationStatement(s)
}

func (s DeclarationStatement) Execute(ctx *Executor) error {
	return ctx.ExecuteDeclarationStatement(s)
}

func (s DeclarationStatement) String() string {
	return fmt.Sprintf("var %s = %s ;", s.name, s.expr)
}
