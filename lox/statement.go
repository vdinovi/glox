package lox

type Program []Statement

type Statement interface {
	String() string
	Print(Printer) (string, error)
	Position() Position
	Equals(Statement) bool
	Execute(*Context) error
	Typecheck(*Context) error
	Resolve(*Context) error
}

type BlockStatement struct {
	stmts []Statement
	pos   Position
}

func (s *BlockStatement) Position() Position {
	return s.pos
}

func (s *BlockStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *ConditionalStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *WhileStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *ForStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *ExpressionStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *PrintStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
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

func (s *DeclarationStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *DeclarationStatement) Equals(other Statement) bool {
	decl, ok := other.(*DeclarationStatement)
	if !ok || s.name != decl.name {
		return false
	}
	return s.expr.Equals(decl.expr)
}

type FunctionDefinitionStatement struct {
	name   string
	params []string
	body   []Statement
	rtype  Type
	pos    Position
}

func (s *FunctionDefinitionStatement) Position() Position {
	return s.pos
}

func (s *FunctionDefinitionStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *FunctionDefinitionStatement) Equals(other Statement) bool {
	o, ok := other.(*FunctionDefinitionStatement)
	if !ok || s.name != o.name || len(s.params) != len(o.params) {
		return false
	}
	for i, p := range s.params {
		if p != o.params[i] {
			return false
		}
	}
	for i, st := range s.body {
		if !st.Equals(o.body[i]) {
			return false
		}
	}
	return true
}

type ReturnStatement struct {
	expr Expression
	typ  Type
	pos  Position
}

func (s *ReturnStatement) Position() Position {
	return s.pos
}

func (s *ReturnStatement) String() string {
	str, err := s.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (s *ReturnStatement) Equals(other Statement) bool {
	ret, ok := other.(*ReturnStatement)
	if !ok {
		return false
	}
	if !s.expr.Equals(ret.expr) {
		return false
	}
	return true
}
