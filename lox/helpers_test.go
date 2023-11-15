package lox

import (
	_ "embed"
	"os"
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

var addOp = Operator{Type: OpAdd, Lexem: "+"}
var subOp = Operator{Type: OpSubtract, Lexem: "-"}
var mulOp = Operator{Type: OpMultiply, Lexem: "*"}
var divOp = Operator{Type: OpDivide, Lexem: "/"}

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
var bAddExpr = makeBinaryExpr(addOp)
var bSubExpr = makeBinaryExpr(subOp)
var bMulExpr = makeBinaryExpr(mulOp)
var bDivExpr = makeBinaryExpr(divOp)
var groupExpr = makeGroupingExpr

// func pos(line, column int) Position {
// 	return Position{Line: line, Column: column}
// }

// var sub = Operator{Type: OpSubtract, Lexem: "-"}
// var add = Operator{Type: OpAdd, Lexem: "+"}

// func exprStmt(e Expression) func(Position) func() ExpressionStatement {
// 	return func(pos Position) func() ExpressionStatement {
// 		return func() ExpressionStatement {
// 			return ExpressionStatement{expr: e, pos: pos}
// 		}
// 	}
// }

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
