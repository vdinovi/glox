package lox

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func (ctx *Context) EvaluateGroupingExpression(e *GroupingExpression) (val Value, typ Type, err error) {
	typ = e.Type()
	if val, err = e.expr.Evaluate(ctx); err != nil {
		return nil, TypeAny, err
	}
	return val, typ, nil
}

func (ctx *Context) EvaluateAssignmentExpression(e *AssignmentExpression) (val Value, typ Type, err error) {
	if val, err = e.right.Evaluate(ctx); err != nil {
		return nil, TypeAny, err
	}
	prev, env := ctx.values.Lookup(e.name)
	if prev == nil {
		return nil, TypeAny, NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
	}
	_, err = env.Set(e.name, val)
	if err != nil {
		return nil, TypeAny, err
	}
	log.Debug().Msgf("(executor) (%d) %s = %s (prev %s)", env.depth, e.name, val, *prev)
	return val, typ, nil
}

func (ctx *Context) EvaluateVariableExpression(e *VariableExpression) (val Value, typ Type, err error) {
	typ = e.Type()
	v, _ := ctx.values.Lookup(e.name)
	if val != nil {
		err := NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
		return nil, TypeAny, err
	}
	val = *v
	return val, typ, nil
}

func (ctx *Context) EvaluateUnaryExpression(e *UnaryExpression) (val Value, typ Type, err error) {
	var right Value
	right, err = e.right.Evaluate(ctx)
	if err != nil {
		return nil, TypeAny, err
	}

	switch e.op.Type {
	case OpNegate:
		return ValueBoolean(!right.Truthy()), TypeBoolean, nil
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
		return nil, TypeAny, err
	}
	return val, val.Type(), err
}

func (ctx *Context) EvaluateBinaryExpression(e *BinaryExpression) (val Value, typ Type, err error) {
	var left, right Value
	left, err = e.left.Evaluate(ctx)
	if err != nil {
		return nil, TypeAny, err
	}
	switch e.op.Type {
	// important that don't evaluate right until we know the value of left so that we can short-circuit evaluation
	case OpAnd:
		if !left.Truthy() {
			return left, left.Type(), nil
		}
		right, err = e.right.Evaluate(ctx)
		if err != nil {
			return nil, TypeAny, err
		}
		return right, right.Type(), nil
	case OpOr:
		if left.Truthy() {
			return left, left.Type(), nil
		}
		right, err = e.right.Evaluate(ctx)
		if err != nil {
			return nil, TypeAny, err
		}
		return right, right.Type(), nil
	}
	right, err = e.right.Evaluate(ctx)
	if err != nil {
		return nil, TypeAny, err
	}

	switch e.op.Type {
	case OpEqualTo:
		return ValueBoolean(left.Equals(right)), TypeBoolean, nil
	case OpNotEqualTo:
		return ValueBoolean(!left.Equals(right)), TypeBoolean, nil
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
		return nil, TypeAny, err
	}
	return val, val.Type(), err
}
