package lox

import (
	"fmt"
)

type Program []Statement

type Statement interface {
	fmt.Stringer
	Printable
	Executable
	Located
	Typecheckable
	Equals(Statement) bool
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
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *BlockStatement) Print(p Printer) (string, error) {
	return p.Print(s)
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
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *ConditionalStatement) Print(p Printer) (string, error) {
	return p.Print(s)
}

func (s *ConditionalStatement) Equals(other Statement) bool {
	cond, ok := other.(*ConditionalStatement)
	if !ok {
		return false
	}
	return s.expr.Equals(cond.expr) &&
		s.thenBranch.Equals(cond.thenBranch) &&
		(s.elseBranch == nil && cond.elseBranch == nil ||
			s.elseBranch.Equals(cond.elseBranch))
}

type WhileStatement struct {
	expr Expression
	body Statement
	pos  Position
}

func (s *WhileStatement) Position() Position {
	return s.pos
}

func (s *WhileStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckWhileStatement(s)
}

func (s *WhileStatement) Execute(ctx *Executor) error {
	return ctx.ExecuteWhileStatement(s)
}

func (s *WhileStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *WhileStatement) Print(p Printer) (string, error) {
	return p.Print(s)
}

func (s *WhileStatement) Equals(other Statement) bool {
	while, ok := other.(*WhileStatement)
	if !ok {
		return false
	}
	return s.expr.Equals(while.expr) &&
		s.body.Equals(while.body)
}

type ForStatement struct {
	init Statement
	cond Expression
	incr Expression
	body Statement
	pos  Position
}

func (s *ForStatement) Position() Position {
	return s.pos
}

func (s *ForStatement) TypeCheck(ctx *Context) error {
	return ctx.TypeCheckForStatement(s)
}

func (s *ForStatement) Execute(ctx *Executor) error {
	return ctx.ExecuteForStatement(s)
}

func (s *ForStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *ForStatement) Print(p Printer) (string, error) {
	return p.Print(s)
}

func (s *ForStatement) Equals(other Statement) bool {
	for_, ok := other.(*ForStatement)
	if !ok {
		return false
	}
	return s.body.Equals(for_.body) &&
		(s.init == nil && for_.init == nil || s.init.Equals(for_.init)) &&
		(s.cond == nil && for_.cond == nil || s.cond.Equals(for_.cond)) &&
		(s.incr == nil && for_.incr == nil || s.incr.Equals(for_.incr))
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
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *ExpressionStatement) Print(p Printer) (string, error) {
	return p.Print(s)
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
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *PrintStatement) Print(p Printer) (string, error) {
	return p.Print(s)
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
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *DeclarationStatement) Print(p Printer) (string, error) {
	return p.Print(s)
}

func (s *DeclarationStatement) Equals(other Statement) bool {
	decl, ok := other.(*DeclarationStatement)
	if !ok || s.name != decl.name {
		return false
	}
	return s.expr.Equals(decl.expr)
}
