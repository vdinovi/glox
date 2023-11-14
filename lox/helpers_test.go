package lox

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	DisableLogger()
	//SetLogLevel("debug")
	os.Exit(m.Run())
}

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

func numExpr(n float64) func() *NumericExpression {
	return func() *NumericExpression {
		return &NumericExpression{value: n}
	}
}

func strExpr(s string) func() *StringExpression {
	return func() *StringExpression {
		return &StringExpression{value: s}
	}
}

func boolExpr(b bool) func() *BooleanExpression {
	return func() *BooleanExpression {
		return &BooleanExpression{value: b}
	}
}

func nilExpr() func() *NilExpression {
	return func() *NilExpression {
		return &NilExpression{}
	}
}

func varExpr(name string) func() *VariableExpression {
	return func() *VariableExpression {
		return &VariableExpression{name: name}
	}
}

func unaryExpr(op Operator) func(Expression) func() *UnaryExpression {
	return func(right Expression) func() *UnaryExpression {
		return func() *UnaryExpression {
			return &UnaryExpression{op: op, right: right}
		}
	}
}

func binaryExpr(op Operator) func(Expression) func(Expression) func() *BinaryExpression {
	return func(left Expression) func(Expression) func() *BinaryExpression {
		return func(right Expression) func() *BinaryExpression {
			return func() *BinaryExpression {
				return &BinaryExpression{op: op, left: left, right: right}
			}
		}
	}
}

func groupExpr(e Expression) func() *GroupingExpression {
	return func() *GroupingExpression {
		return &GroupingExpression{expr: e}
	}
}
