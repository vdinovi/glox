package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) EvaluateUnaryExpression(e *UnaryExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateUnaryExpression")
	var right Value
	right, err = e.right.Evaluate(ctx)
	if err != nil {
		return nil, ErrType, err
	}
	switch typ := right.Type(); typ {
	case TypeNumeric:
		val, err = ctx.evalUnaryNumeric(e, right)
	default:
		err = NewTypeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right.Type()), e.Position())
	}
	if err != nil {
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil

}

func (ctx *Context) EvaluateBinaryExpression(e *BinaryExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateBinaryExpression")
	var left, right Value
	left, err = e.left.Evaluate(ctx)
	if err == nil {
		right, err = e.right.Evaluate(ctx)
	}
	if err != nil {
		return nil, ErrType, err
	}
	if left.Type() != right.Type() {
		return nil, ErrType, NewTypeError(NewTypeMismatchError(left.Type(), right.Type()), e.Position())
	}
	switch typ := left.Type(); typ {
	case TypeString:
		val, err = ctx.evalBinaryString(e, left, right)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e, left, right)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e, left, right)
	default:
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	if err != nil {
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateGroupingExpression(e *GroupingExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateGroupingExpression")
	typ = e.Type()
	if val, err = e.expr.Evaluate(ctx); err != nil {
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateAssignmentExpression(e *AssignmentExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateAssignmentExpression")
	if val, err = e.right.Evaluate(ctx); err == nil {
		err = ctx.values.Set(e.name, val)
	}
	if err != nil {
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateVariableExpression(e *VariableExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateVariableExpression")
	typ = e.Type()
	v := ctx.values.Lookup(e.name)
	if val != nil {
		err := NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
		return nil, ErrType, err
	}
	val = *v
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) evalUnaryNumeric(e *UnaryExpression, right Value) (val Value, err error) {
	num, ok := right.Unwrap().(float64)
	if !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "float64"), Position{})
	}
	switch e.op.Type {
	case OpSubtract:
		val = ValueNumeric(-num)
	default:
		err = NewTypeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right.Type()), e.Position())
	}
	return val, nil
}

func (ctx *Context) evalBinaryString(e *BinaryExpression, left, right Value) (val Value, err error) {
	var ok bool
	var l, r string
	if l, ok = left.Unwrap().(string); !ok {
		err = NewRuntimeError(NewDowncastError(left, "string"), e.Position())
	}
	if ok {
		if r, ok = right.Unwrap().(string); !ok {
			err = NewRuntimeError(NewDowncastError(right, "string"), e.Position())
		}
	}
	if err != nil {
		return nil, err
	}
	switch e.op.Type {
	case OpAdd:
		val = ValueString(l + r)
	case OpEqualTo:
		val = ValueBoolean(l == r)
	case OpNotEqualTo:
		val = ValueBoolean(l != r)
	default:
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	return val, nil
}

func (ctx *Context) evalBinaryNumeric(e *BinaryExpression, left, right Value) (val Value, err error) {
	var ok bool
	var l, r float64
	if l, ok = left.Unwrap().(float64); !ok {
		err = NewRuntimeError(NewDowncastError(left, "float64"), e.Position())
	}
	if ok {
		if r, ok = right.Unwrap().(float64); !ok {
			err = NewRuntimeError(NewDowncastError(right, "float64"), e.Position())
		}
	}
	if err != nil {
		return nil, err
	}
	switch e.op.Type {
	case OpAdd:
		val = ValueNumeric(l + r)
	case OpSubtract:
		val = ValueNumeric(l - r)
	case OpMultiply:
		val = ValueNumeric(l * r)
	case OpDivide:
		val = ValueNumeric(l / r)
	case OpEqualTo:
		val = ValueBoolean(l == r)
	case OpNotEqualTo:
		val = ValueBoolean(l != r)
	case OpLessThan:
		val = ValueBoolean(l < r)
	case OpLessThanOrEqualTo:
		val = ValueBoolean(l <= r)
	case OpGreaterThan:
		val = ValueBoolean(l > r)
	case OpGreaterThanOrEqualTo:
		val = ValueBoolean(l >= r)
	default:
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	return val, nil
}

func (ctx *Context) evalBinaryBoolean(e *BinaryExpression, left, right Value) (val Value, err error) {
	var l, r bool
	var ok bool
	if l, ok = left.Unwrap().(bool); !ok {
		err = NewRuntimeError(NewDowncastError(left, "bool"), e.Position())
	}
	if ok {
		if r, ok = right.Unwrap().(bool); !ok {
			err = NewRuntimeError(NewDowncastError(right, "bool"), e.Position())
		}
	}
	if err != nil {
		return nil, err
	}
	switch e.op.Type {
	case OpEqualTo:
		val = ValueBoolean(l == r)
	case OpNotEqualTo:
		val = ValueBoolean(l != r)
	default:
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	return val, nil
}
