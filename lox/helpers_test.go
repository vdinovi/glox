package lox

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	DisableLogger()
	//SetLogLevel("debug")
	os.Exit(m.Run())
}

//go:embed fixtures/program.lox
var fixtureProgram string

//go:embed fixtures/program_tokens.json
var fixtureProgramTokens string

var negOp = Operator{Type: OpNegate, Lexem: "!"}
var eqOp = Operator{Type: OpEqualTo, Lexem: "=="}
var neqOp = Operator{Type: OpNotEqualTo, Lexem: "!="}
var ltOp = Operator{Type: OpLessThan, Lexem: "<"}
var lteOp = Operator{Type: OpLessThanOrEqualTo, Lexem: "<="}
var gtOp = Operator{Type: OpGreaterThan, Lexem: ">"}
var gteOp = Operator{Type: OpGreaterThanOrEqualTo, Lexem: ">="}
var addOp = Operator{Type: OpAdd, Lexem: "+"}
var subOp = Operator{Type: OpSubtract, Lexem: "-"}
var mulOp = Operator{Type: OpMultiply, Lexem: "*"}
var divOp = Operator{Type: OpDivide, Lexem: "/"}
var andOp = Operator{Type: OpAnd, Lexem: "and"}
var orOp = Operator{Type: OpOr, Lexem: "or"}

var zeroExpr = makeNumericExpr(0)
var oneExpr = makeNumericExpr(1)
var piExpr = makeNumericExpr(3.14)
var strExpr = makeStringExpr("str")
var trueExpr = makeBooleanExpr(true)
var falseExpr = makeBooleanExpr(false)
var nilExpr = makeNilExpr()
var fooExpr = makeVarExpr("foo")

var uSubExpr = makeUnaryExpr(subOp)
var uAddExpr = makeUnaryExpr(addOp)
var uNegExpr = makeUnaryExpr(negOp)

var bAddExpr = makeBinaryExpr(addOp)
var bSubExpr = makeBinaryExpr(subOp)
var bMulExpr = makeBinaryExpr(mulOp)
var bDivExpr = makeBinaryExpr(divOp)

var groupExpr = makeGroupingExpr

var eqExpr = makeBinaryExpr(eqOp)
var neqExpr = makeBinaryExpr(neqOp)
var ltExpr = makeBinaryExpr(ltOp)
var lteExpr = makeBinaryExpr(lteOp)
var gtExpr = makeBinaryExpr(gtOp)
var gteExpr = makeBinaryExpr(gteOp)

var bAndExpr = makeBinaryExpr(andOp)
var bOrExpr = makeBinaryExpr(orOp)

var fooCallExpr = makeCallExpression(fooExpr())

func makeNumericExpr(n float64) func() *NumericExpression {
	return func() *NumericExpression {
		return &NumericExpression{value: n}
	}
}

func makeStringExpr(s string) func() *StringExpression {
	return func() *StringExpression {
		return &StringExpression{value: s}
	}
}

func makeBooleanExpr(b bool) func() *BooleanExpression {
	return func() *BooleanExpression {
		return &BooleanExpression{value: b}
	}
}

func makeNilExpr() func() *NilExpression {
	return func() *NilExpression {
		return &NilExpression{}
	}
}

func makeVarExpr(name string) func() *VariableExpression {
	return func() *VariableExpression {
		return &VariableExpression{name: name}
	}
}

func makeUnaryExpr(op Operator) func(Expression) func() *UnaryExpression {
	return func(right Expression) func() *UnaryExpression {
		return func() *UnaryExpression {
			return &UnaryExpression{op: op, right: right}
		}
	}
}

func makeBinaryExpr(op Operator) func(Expression) func(Expression) func() *BinaryExpression {
	return func(left Expression) func(Expression) func() *BinaryExpression {
		return func(right Expression) func() *BinaryExpression {
			return func() *BinaryExpression {
				return &BinaryExpression{op: op, left: left, right: right}
			}
		}
	}
}

func makeGroupingExpr(e Expression) func() *GroupingExpression {
	return func() *GroupingExpression {
		return &GroupingExpression{expr: e}
	}
}

func makeCallExpression(callee Expression) func(args ...Expression) func() *CallExpression {
	return func(args ...Expression) func() *CallExpression {
		return func() *CallExpression {
			return &CallExpression{callee: callee, args: args}
		}
	}
}

type testDriverPhase int

const (
	initializing = iota
	lexing
	parsing
	typechecking
	executing
)

type TestDriver struct {
	Text    string
	Tokens  []Token
	Program []Statement
	Printer PrintSpy
	Exec    *Executor
	Err     error
	t       *testing.T
	phase   testDriverPhase
}

func NewTestDriver(t *testing.T, text string) *TestDriver {
	return &TestDriver{
		Text: text,
		t:    t,
	}
}

func (td *TestDriver) Fatal() {
	td.t.Helper()
	if td.Err != nil {
		td.t.Fatalf("Unexpected error while %s: %s", td.Phase(), td.Err)
	}
}

func (td *TestDriver) Error() {
	td.t.Helper()
	if td.Err != nil {
		td.t.Errorf("Unexpected error while %s: %s", td.Phase(), td.Err)
	}
}

func (td *TestDriver) Lex() {
	td.phase = lexing
	td.Tokens, td.Err = Scan(strings.NewReader(td.Text))
}

func (td *TestDriver) Parse() {
	if td.Err != nil {
		return
	}
	if len(td.Tokens) < 1 {
		td.Err = fmt.Errorf("no tokens to parse (ensure Lex has been called)")
	}
	td.phase = parsing
	td.Program, td.Err = Parse(td.Tokens)
}

func (td *TestDriver) TypeCheck() {
	if td.Err != nil {
		return
	}
	if len(td.Program) < 1 {
		td.Err = fmt.Errorf("no program to typecheck (ensure Parse has been called)")
	}
	if td.Exec == nil {
		td.Exec = NewExecutor(&td.Printer)
	}
	td.phase = parsing
	td.Err = td.Exec.TypeCheckProgram(td.Program)
}

func (td *TestDriver) Execute() {
	if td.Err != nil {
		return
	}
	if len(td.Program) < 1 {
		td.Err = fmt.Errorf("no program to execute (ensure Parse has been called)")
	}
	if td.Exec == nil {
		td.Exec = NewExecutor(&td.Printer)
	}
	td.phase = executing
	td.Err = td.Exec.ExecuteProgram(td.Program)
}

func (td *TestDriver) Phase() string {
	switch td.phase {
	case initializing:
		return "initializing"
	case lexing:
		return "lexing"
	case parsing:
		return "parsing"
	case typechecking:
		return "typechecking"
	case executing:
		return "executing"
	}
	panic(td.phase)
}

type PrintSpy struct {
	Buffer strings.Builder
	Prints []string
}

func (s *PrintSpy) Write(p []byte) (n int, err error) {
	for _, b := range p {
		switch b {
		case '\n':
			s.Prints = append(s.Prints, s.Buffer.String())
			s.Buffer.Reset()
		default:
			if err := s.Buffer.WriteByte(b); err != nil {
				return n, err
			}
			n += 1
		}
	}
	return n, nil
}
