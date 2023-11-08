package lox

import "fmt"

type OperatorType int

//go:generate stringer -type OperatorType -trimprefix=Op
const (
	OpNone OperatorType = iota
	OpPlus
	OpMinus
	OpMultiply
	OpDivide
	OpEquals
	OpEqualEquals
	OpNotEquals
	OpLess
	OpLessEquals
	OpGreater
	OpGreaterEquals
)

type Operator struct {
	Type  OperatorType
	Lexem string
}

var NoneOp = Operator{Type: OpNone}

func (o Operator) String() string {
	return fmt.Sprintf("%s(%s)", o.Type, o.Lexem)
}
