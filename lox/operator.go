package lox

import "fmt"

type OperatorType int

//go:generate stringer -type OperatorType -trimprefix=Op
const (
	ErrOp OperatorType = iota
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpEqualTo
	OpNotEqualTo
	OpLessThan
	OpLessThanOrEqualTo
	OpGreaterThan
	OpGreaterThanOrEqualTo
)

type Operator struct {
	Type   OperatorType // type of the operator
	Lexem  string       // associated string
	Line   int          // originating line
	Column int          // originating column
}

func (o Operator) String() string {
	return fmt.Sprintf("%s(%s)", o.Type, o.Lexem)
}

func OperatorDefault(t OperatorType) Operator {
	op := Operator{Type: t, Lexem: ""}
	switch t {
	case OpAdd:
		op.Lexem = "+"
	case OpSubtract:
		op.Lexem = "-"
	case OpMultiply:
		op.Lexem = "*"
	case OpDivide:
		op.Lexem = "/"
	case OpEqualTo:
		op.Lexem = "="
	case OpNotEqualTo:
		op.Lexem = "!="
	case OpLessThan:
		op.Lexem = "<"
	case OpLessThanOrEqualTo:
		op.Lexem = "<="
	case OpGreaterThan:
		op.Lexem = ">"
	case OpGreaterThanOrEqualTo:
		op.Lexem = ">="
	}
	return op
}
