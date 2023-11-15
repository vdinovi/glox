package lox

import (
	"fmt"
	"strings"
)

type Program []Statement

type Statement interface {
	Position() Position
	Equals(Statement) bool
	TypeCheck(ctx *Context) error
	Execute(*Executor) error
	fmt.Stringer
}

type BlockStatement struct {
	stmts []Statement
	pos   Position
}

func (s *BlockStatement) Position() Position {
	return s.pos
}

func (s *BlockStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckBlockStatement(s)
}

func (s *BlockStatement) Execute(e *Executor) error {
	return e.ExecuteBlockStatement(s)
}

func (s *BlockStatement) String() string {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.stmts {
		fmt.Fprint(&sb, " ", stmt.String())
	}
	sb.WriteString(" } ;")
	return sb.String()
}

func (s *BlockStatement) Equals(other Statement) bool {
	block, ok := other.(*BlockStatement)
	if !ok || len(s.stmts) != len(block.stmts) {
		return false
	}
	for i, st := range s.stmts {
		if !st.Equals(block.stmts[i]) {
			return false
		}
	}
	return true
}

type ConditionalStatement struct {
	expr       Expression
	thenBranch Statement
	elseBranch Statement
	pos        Position
}

func (s *ConditionalStatement) Position() Position {
	return s.pos
}

func (s *ConditionalStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckConditionalStatement(s)
}

func (s *ConditionalStatement) Execute(ctx *Executor) error {
	return ctx.ExecuteConditionalStatement(s)
}

func (s *ConditionalStatement) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "if %s %s", s.expr, s.thenBranch)
	if s.elseBranch != nil {
		fmt.Fprintf(&sb, "else %s", s.elseBranch)
	}
	return sb.String()
}

func (s *ConditionalStatement) Equals(other Statement) bool {
	cond, ok := other.(*ConditionalStatement)
	if !ok {
		return false
	}
	return s.expr.Equals(cond.expr) &&
		s.thenBranch.Equals(cond.thenBranch) &&
		(s.elseBranch == nil && cond.elseBranch == nil || s.elseBranch.Equals(cond.elseBranch))
}

type ExpressionStatement struct {
	expr Expression
	pos  Position
}

func (s *ExpressionStatement) Position() Position {
	return s.pos
}

func (s *ExpressionStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckExpressionStatement(s)
}

func (s *ExpressionStatement) Execute(e *Executor) error {
	return e.ExecuteExpressionStatement(s)
}

func (s *ExpressionStatement) String() string {
	return fmt.Sprintf("%s ;", s.expr)
}

func (s *ExpressionStatement) Equals(other Statement) bool {
	expr, ok := other.(*ExpressionStatement)
	if !ok {
		return false
	}
	return s.expr.Equals(expr.expr)
}

type PrintStatement struct {
	expr Expression
	pos  Position
}

func (s *PrintStatement) Position() Position {
	return s.pos
}

func (s *PrintStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckPrintStatement(s)
}

func (s *PrintStatement) Execute(e *Executor) error {
	return e.ExecutePrintStatement(s)
}

func (s *PrintStatement) String() string {
	return fmt.Sprintf("print %s ;", s.expr)
}

func (s *PrintStatement) Equals(other Statement) bool {
	print, ok := other.(*PrintStatement)
	if !ok {
		return false
	}
	return s.expr.Equals(print.expr)
}

type DeclarationStatement struct {
	name string
	pos  Position
	expr Expression
}

func (s *DeclarationStatement) Position() Position {
	return s.pos
}

func (s *DeclarationStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckDeclarationStatement(s)
}

func (s *DeclarationStatement) Execute(ctx *Executor) error {
	return ctx.ExecuteDeclarationStatement(s)
}

func (s *DeclarationStatement) String() string {
	return fmt.Sprintf("var %s = %s ;", s.name, s.expr)
}

func (s *DeclarationStatement) Equals(other Statement) bool {
	decl, ok := other.(*DeclarationStatement)
	if !ok || s.name != decl.name {
		return false
	}
	return s.expr.Equals(decl.expr)
}
