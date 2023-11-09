package lox

import "fmt"

type OperatorType int

//go:generate stringer -type OperatorType -trimprefix=Op
const (
	OpPlus OperatorType = iota
	OpMinus
	OpMultiply
	OpDivide
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

func (o Operator) String() string {
	return fmt.Sprintf("%s(%s)", o.Type, o.Lexem)
}
