package lox

// import (
// 	"fmt"
// 	"testing"
// )

// var operators = []Operator{
// 	{Type: OpPlus, Lexem: "+"},
// 	{Type: OpMinus, Lexem: "-"},
// 	{Type: OpMultiply, Lexem: "*"},
// 	{Type: OpDivide, Lexem: "/"},
// 	{Type: OpEqualEquals, Lexem: "=="},
// 	{Type: OpNotEquals, Lexem: "!="},
// 	{Type: OpLess, Lexem: "<"},
// 	{Type: OpLessEquals, Lexem: "<="},
// 	{Type: OpGreater, Lexem: ">"},
// 	{Type: OpGreaterEquals, Lexem: ">="},
// }

// func Test_Operator_String(t *testing.T) {
// 	for _, operator := range operators {
// 		want := fmt.Sprintf("%s(%s)", operator.Type, operator.Lexem)
// 		got := operator.String()
// 		if got != want {
// 			t.Errorf("Expected %s to yield %s, but got %s", operator.Type, want, got)
// 		}
// 	}
// }
