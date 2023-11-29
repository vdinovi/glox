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

type TestDriver struct {
	Text    string
	Tokens  []Token
	Program []Statement
	Printer PrintSpy
	ctx     *Context
	Err     error
	t       *testing.T
	phase   Phase
}

func NewTestDriver(t *testing.T, text string) *TestDriver {
	td := &TestDriver{
		phase: PhaseInit,
		Text:  text,
		t:     t,
	}
	td.ctx = NewContext(&td.Printer)
	return td
}

func (td *TestDriver) Phase() Phase {
	return td.phase
}

func (td *TestDriver) Fatal() {
	td.t.Helper()
	if td.Err != nil {
		td.unexpectedError(td.Err)
	}
}

func (td *TestDriver) Error() {
	td.t.Helper()
	if td.Err != nil {
		td.unexpectedError(td.Err)
	}
}

func (td *TestDriver) unexpectedError(err error) {
	td.t.Helper()
	td.t.Errorf("Unexpected error in phase %s: %s", td.Phase(), err)
}

func (td *TestDriver) Lex() {
	td.phase = PhaseLex
	td.Tokens, td.Err = Scan(td.ctx, strings.NewReader(td.Text))
}

func (td *TestDriver) Parse() {
	if td.Err != nil {
		return
	}
	if len(td.Tokens) < 1 {
		td.Err = fmt.Errorf("no tokens to parse (ensure Lex has been called)")
	}
	td.phase = PhaseParse
	td.Program, td.Err = Parse(td.ctx, td.Tokens)
}

func (td *TestDriver) TypeCheck() {
	if td.Err != nil {
		return
	}
	if len(td.Program) < 1 {
		td.Err = fmt.Errorf("no program to typecheck (ensure Parse has been called)")
	}
	td.phase = PhaseTypecheck
	td.Err = Typecheck(td.ctx, td.Program)
}

func (td *TestDriver) Execute() {
	if td.Err != nil {
		return
	}
	if len(td.Program) < 1 {
		td.Err = fmt.Errorf("no program to execute (ensure Parse has been called)")
	}
	td.phase = PhaseExecute
	td.Err = Execute(td.ctx, td.Program)
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
