package lox

import (
	"errors"

	"github.com/rs/zerolog/log"
)

type Evaluable interface {
	Evaluate(*Context) (Value, error)
}

func (e *GroupingExpression) Evaluate(ctx *Context) (Value, error) {
	return e.expr.Evaluate(ctx)
}

func (e *AssignmentExpression) Evaluate(ctx *Context) (Value, error) {
	val, err := e.right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	prev, env := ctx.values.Lookup(e.name)
	if prev == nil {
		return nil, NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
	}
	if _, err = env.Set(e.name, val); err != nil {
		return nil, err
	}
	log.Debug().Msgf("(evaluate) (%d) %s = %s (prev %s)", env.depth, e.name, val, *prev)
	return val, nil
}

func (e *VariableExpression) Evaluate(ctx *Context) (Value, error) {
	if val, _ := ctx.values.Lookup(e.name); val != nil {
		return *val, nil
	}
	return nil, NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
}

func (e *CallExpression) Evaluate(ctx *Context) (Value, error) {
	return nil, ErrNotYetImplemented
}

func (e *UnaryExpression) Evaluate(ctx *Context) (Value, error) {
	right, err := e.right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return evaluateUnary(ctx, e, right)
}

func (e *BinaryExpression) Evaluate(ctx *Context) (Value, error) {
	switch e.op.Type {
	case OpAnd, OpOr:
		return evaluateBinaryWithShortCircuit(ctx, e)
	}
	left, err := e.left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	right, err := e.right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return evaluateBinary(ctx, e, left, right)

}

func (e *StringExpression) Evaluate(*Context) (Value, error) {
	return ValueString(e.value), nil
}

func (e *NumericExpression) Evaluate(*Context) (Value, error) {
	return ValueNumeric(e.value), nil
}

func (e *BooleanExpression) Evaluate(*Context) (Value, error) {
	return ValueBoolean(e.value), nil
}

func (e *NilExpression) Evaluate(*Context) (Value, error) {
	return Nil, nil
}

func evaluateUnary(ctx *Context, e *UnaryExpression, right Value) (val Value, err error) {
	if e.op.Type == OpNegate {
		return ValueBoolean(!right.Truthy()), nil
	}

	var invalid bool
	switch right.Type() {
	case TypeNumeric:
		var n ValueNumeric
		n, ok := right.(ValueNumeric)
		if invalid = !ok; invalid {
			break
		}
		switch e.op.Type {
		case OpAdd:
			val = right
		case OpSubtract:
			val, err = n.Negative()
		default:
			invalid = true
		}
	case TypeAny:
		switch e.op.Type {
		case OpNegate:
			val = ValueBoolean(right.Truthy())
		default:
			invalid = true
		}
	default:
		invalid = true
	}
	if err == ErrInvalidType || err == nil && invalid {
		err = NewInvalidUnaryOperatorForTypeError(e.op.Type, right.Type())
	}
	if err != nil {
		if !errors.Is(err, RuntimeError{}) {
			err = NewRuntimeError(err, e.Position())
		}
		return nil, err
	}
	return val, err
}

func evaluateBinaryWithShortCircuit(ctx *Context, e *BinaryExpression) (val Value, err error) {
	left, err := e.left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	var invalid bool
	switch e.op.Type {
	case OpAnd:
		val, err = left, nil
		if val.Truthy() {
			val, err = e.right.Evaluate(ctx)
		}
	case OpOr:
		val, err = left, nil
		if !left.Truthy() {
			val, err = e.right.Evaluate(ctx)
		}
	default:
		invalid = true
	}
	if err == ErrInvalidType || err == nil && invalid {
		err = NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), e.right.Type())
	}
	if err != nil {
		if !errors.Is(err, RuntimeError{}) {
			err = NewRuntimeError(err, e.Position())
		}
		return nil, err
	}
	return val, err
}

func evaluateBinary(ctx *Context, e *BinaryExpression, left, right Value) (val Value, err error) {
	switch e.op.Type {
	case OpEqualTo:
		return ValueBoolean(left.Equals(right)), nil
	case OpNotEqualTo:
		return ValueBoolean(!left.Equals(right)), nil
	}

	var invalid bool
	var cmp int
	switch left.Type() {
	case TypeNumeric:
		var n ValueNumeric
		n, ok := left.(ValueNumeric)
		if invalid = !ok; invalid {
			break
		}
		switch e.op.Type {
		case OpAdd:
			val, err = n.Add(right)
		case OpSubtract:
			val, err = n.Subtract(right)
		case OpMultiply:
			val, err = n.Multiply(right)
		case OpDivide:
			val, err = n.Divide(right)
		case OpLessThan:
			cmp, err = n.Compare(right)
			if err == nil {
				val = ValueBoolean(cmp < 0)
			}
		case OpLessThanOrEqualTo:
			cmp, err = n.Compare(right)
			if err == nil {
				val = ValueBoolean(cmp <= 0)
			}
		case OpGreaterThan:
			cmp, err = n.Compare(right)
			if err == nil {
				val = ValueBoolean(cmp > 0)
			}
		case OpGreaterThanOrEqualTo:
			cmp, err = n.Compare(right)
			if err == nil {
				val = ValueBoolean(cmp >= 0)
			}
		default:
			invalid = true
		}
	case TypeString:
		var s ValueString
		s, ok := left.(ValueString)
		if invalid = !ok; invalid {
			break
		}
		switch e.op.Type {
		case OpAdd:
			val, err = s.Concat(right)
		default:
			invalid = true
		}
	default:
		invalid = true
	}
	if err == ErrInvalidType || err == nil && invalid {
		err = NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type())
	}
	if err != nil {
		if !errors.Is(err, RuntimeError{}) {
			err = NewRuntimeError(err, e.Position())
		}
		return nil, err
	}
	return val, err
}
