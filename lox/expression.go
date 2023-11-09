package lox

import (
	"fmt"
)

type Evaluable interface {
	Evaluate() (Value, Type, error)
}

type Expression interface {
	Typed
	Evaluable
	fmt.Stringer
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
	}
	return -1, TypeError{e, fmt.Errorf("can't apply unary operator %s to type %s", e.op, rightType)}
}

func (e UnaryExpression) Evaluate() (Value, Type, error) {
	val, typ, err := e.right.Evaluate()
	if err != nil {
		return nil, -1, err
	}
	switch typ {
	case TypeNumeric:
		return e.evalUnaryNumeric(e.op, val)
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, val)}
}

func (e UnaryExpression) evalUnaryNumeric(op Operator, val Value) (Value, Type, error) {
	n, ok := val.Unwrap().(float64)
	if ok {
		return nil, -1, RuntimeError{Err: InvalidOperation(e.op, val)}
	}
	switch op.Type {
	case OpMinus:
		return NumericValue(-n), TypeNumeric, nil
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, val)}
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

func (e BinaryExpression) Evaluate() (Value, Type, error) {
	leftVal, leftType, err := e.left.Evaluate()
	if err != nil {
		return nil, -1, err
	}
	rightVal, rightType, err := e.right.Evaluate()
	if err != nil {
		return nil, -1, err
	}
	if leftType != rightType {
		return nil, -1, RuntimeError{Err: InvalidOperation(e.op, leftVal, rightVal)}
	}
	switch leftType {
	case TypeString:
		return e.evalBinaryString(e.op, leftVal, rightVal)
	case TypeNumeric:
		return e.evalBinaryNumeric(e.op, leftVal, rightVal)
	case TypeBoolean:
		return e.evalBinaryBoolean(e.op, leftVal, rightVal)
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, leftVal, rightVal)}
}

func (e BinaryExpression) evalBinaryString(op Operator, left, right Value) (Value, Type, error) {
	a, aOk := left.Unwrap().(string)
	b, bOk := right.Unwrap().(string)
	if aOk && bOk {
		switch op.Type {
		case OpPlus:
			return StringValue(a + b), TypeString, nil
		case OpEqualEquals:
			return BooleanValue(a == b), TypeBoolean, nil
		case OpNotEquals:
			return BooleanValue(a != b), TypeBoolean, nil
		}
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, left, right)}
}

func (e BinaryExpression) evalBinaryNumeric(op Operator, left, right Value) (Value, Type, error) {
	a, aAk := left.Unwrap().(float64)
	b, bOk := right.Unwrap().(float64)
	if aAk && bOk {
		switch op.Type {
		case OpPlus:
			return NumericValue(a + b), TypeNumeric, nil
		case OpMinus:
			return NumericValue(a - b), TypeNumeric, nil
		case OpMultiply:
			return NumericValue(a * b), TypeNumeric, nil
		case OpDivide:
			return NumericValue(a / b), TypeNumeric, nil
		case OpEqualEquals:
			return BooleanValue(a == b), TypeBoolean, nil
		case OpNotEquals:
			return BooleanValue(a != b), TypeBoolean, nil
		case OpLess:
			return BooleanValue(a < b), TypeBoolean, nil
		case OpLessEquals:
			return BooleanValue(a <= b), TypeBoolean, nil
		case OpGreater:
			return BooleanValue(a > b), TypeBoolean, nil
		case OpGreaterEquals:
			return BooleanValue(a >= b), TypeBoolean, nil
		}
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, left, right)}
}

func (e BinaryExpression) evalBinaryBoolean(op Operator, left, right Value) (Value, Type, error) {
	a, aok := left.Unwrap().(bool)
	b, bok := right.Unwrap().(bool)
	if !aok || !bok {
		return nil, -1, RuntimeError{Err: InvalidOperation(e.op, left, right)}
	}
	switch op.Type {
	case OpEqualEquals:
		return BooleanValue(a == b), TypeBoolean, nil
	case OpNotEquals:
		return BooleanValue(a != b), TypeBoolean, nil
	}
	return nil, -1, RuntimeError{Err: InvalidOperation(e.op, left, right)}
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

func (e GroupingExpression) Evaluate() (Value, Type, error) {
	return e.expr.Evaluate()
}

type LiteralExpression interface {
	literal()
}

type StringExpression string

func (e StringExpression) literal() {}

func (e StringExpression) Type() (Type, error) {
	return TypeString, nil
}

func (e StringExpression) String() string {
	return string(e)
}

func (e StringExpression) Evaluate() (Value, Type, error) {
	return StringValue(e), TypeString, nil
}

type NumericExpression float64

func (e NumericExpression) literal() {}

func (e NumericExpression) Type() (Type, error) {
	return TypeNumeric, nil
}

func (e NumericExpression) String() string {
	return fmt.Sprint(float64(e))
}

func (e NumericExpression) Evaluate() (Value, Type, error) {
	return NumericValue(e), TypeNumeric, nil
}

type BooleanExpression bool

func (e BooleanExpression) literal() {}

func (e BooleanExpression) String() string {
	if bool(e) {
		return "true"
	}
	return "false"
}

func (e BooleanExpression) Type() (Type, error) {
	return TypeBoolean, nil
}
func (e BooleanExpression) Evaluate() (Value, Type, error) {
	return BooleanValue(e), TypeBoolean, nil
}

type NilExpression struct{}

func (e NilExpression) literal() {}

func (e NilExpression) String() string {
	return "nil"
}

func (e NilExpression) Type() (Type, error) {
	return TypeNil, nil
}
func (e NilExpression) Evaluate() (Value, Type, error) {
	return NilValue(e), TypeNil, nil
}
