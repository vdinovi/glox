package lox

import (
	"fmt"
)

type Expression interface {
	fmt.Stringer
	Type() (Type, error)
}

type UnaryExpression struct {
	op    Operator
	right Expression
}

func (e UnaryExpression) String() string {
	return fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
}

func (e UnaryExpression) Type() (Type, error) {
	rightType, err := e.right.Type()
	if err != nil {
		return -1, err
	}
	switch e.op.Type {
	case OpMinus:
		if rightType != TypeNumeric {
			return -1, TypeError{
				Expression: e,
				Err:        fmt.Errorf("can't apply %s to %s", e.op, rightType),
			}
		}
		return TypeNumeric, nil
		// case OpBang:
		// 	if rightType != TypeBoolean {
		// 		return -1, TypeError{
		// 			Expression: e,
		// 			Err:  fmt.Errorf("can't apply %s to %s", e.op, rightType),
		// 		}
		// 	}
		// 	return TypeBoolean, nil
	}
	return -1, TypeError{
		Expression: e,
		Err:        fmt.Errorf("can't apply unary operator %s to type %s", e.op, rightType),
	}
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
}

func (e BinaryExpression) Type() (Type, error) {
	leftType, err := e.left.Type()
	if err != nil {
		return -1, err
	}
	rightType, err := e.right.Type()
	if err != nil {
		return -1, err
	}
	if leftType != rightType {
		return -1, TypeError{
			Expression: e,
			Err:        fmt.Errorf("can't apply binary operator %s to mismatched types %s and %s", e.op, leftType, rightType),
		}
	}
	switch leftType {
	case TypeNumeric:
		switch e.op.Type {
		case OpPlus, OpMinus, OpMultiply, OpDivide:
			return TypeNumeric, nil
		case OpEqualEquals, OpNotEquals, OpLess, OpLessEquals, OpGreater, OpGreaterEquals:
			return TypeBoolean, nil
		}
	case TypeString:
		switch e.op.Type {
		case OpPlus:
			return TypeString, nil
		case OpEqualEquals, OpNotEquals:
			return TypeBoolean, nil
		}
	case TypeBoolean:
		switch e.op.Type {
		case OpEqualEquals, OpNotEquals:
			return TypeBoolean, nil
		}
	}
	return -1, TypeError{
		Expression: e,
		Err:        fmt.Errorf("can't apply binary operator %s to types %s and %s", e.op, leftType, rightType),
	}
}

type GroupingExpression struct {
	expr Expression
}

func (e GroupingExpression) String() string {
	return fmt.Sprintf("(group %s)", e.expr)
}

func (e GroupingExpression) Type() (Type, error) {
	return e.expr.Type()
}

type LiteralExpression interface {
	literal()
}

type StringExpression string

func (e StringExpression) literal() {}

func (e StringExpression) String() string {
	return string(e)
}

func (e StringExpression) Type() (Type, error) {
	return TypeString, nil
}

type NumberExpression float64

func (e NumberExpression) literal() {}

func (e NumberExpression) String() string {
	return fmt.Sprint(float64(e))
}

func (e NumberExpression) Type() (Type, error) {
	return TypeNumeric, nil
}

type BooleanExpression bool

func (e BooleanExpression) literal() {}

func (e BooleanExpression) String() string {
	if bool(e) {
		return "true"
	} else {
		return "false"
	}
}

func (e BooleanExpression) Type() (Type, error) {
	return TypeBoolean, nil
}

type NilExpression struct{}

func (e NilExpression) literal() {}

func (e NilExpression) String() string {
	return "nil"
}

func (e NilExpression) Type() (Type, error) {
	return TypeNil, nil
}
