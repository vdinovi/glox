package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) EvaluateUnaryExpression(e *UnaryExpression) (val Value, typ Type, err error) {
	switch typ := e.Type(); typ {
	case TypeNumeric:
		val, err = e.right.Evaluate(ctx)
		if err != nil {
			return nil, ErrType, err
		}
		val, err = ctx.evalUnaryNumeric(e.op, val)
	default:
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	if err != nil {
		return nil, ErrType, err
	}
	return val, typ, nil

}

func (ctx *Context) EvaluateBinaryExpression(e *BinaryExpression) (val Value, typ Type, err error) {
	switch typ := e.Type(); typ {
	case TypeString:
		val, err = ctx.evalBinaryString(e)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e)
	default:
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) binary expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateGroupingExpression(e *GroupingExpression) (val Value, typ Type, err error) {
	typ = e.Type()
	if val, err = e.expr.Evaluate(ctx); err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) grouping expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateAssignmentExpression(e *AssignmentExpression) (val Value, typ Type, err error) {
	if val, err = e.right.Evaluate(ctx); err == nil {
		err = ctx.values.Set(e.name, val)
	}
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) assignment expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateVariableExpression(e *VariableExpression) (val Value, typ Type, err error) {
	typ = e.Type()
	v := ctx.values.Lookup(e.name)
	if val == nil {
		err := NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	val = *v
	log.Debug().Msgf("(evaluator) variable expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) evalUnaryNumeric(op Operator, right Value) (val Value, err error) {
	num, ok := right.Unwrap().(float64)
	if ok {
		return nil, NewRuntimeError(NewDowncastError(right, "float64"), Position{})
	}
	switch op.Type {
	case OpSubtract:
		val = ValueNumeric(-num)
	default:
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	return val, nil
}

func (ctx *Context) evalBinaryString(e *BinaryExpression) (val Value, err error) {
	var left, right Value
	var ok bool
	left, err = e.left.Evaluate(ctx)
	if err == nil {
		right, err = e.right.Evaluate(ctx)
	}
	if err != nil {
		return nil, err
	}
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
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	return val, nil
}

func (ctx *Context) evalBinaryNumeric(e *BinaryExpression) (val Value, err error) {
	var left, right Value
	var ok bool
	left, err = e.left.Evaluate(ctx)
	if err == nil {
		right, err = e.right.Evaluate(ctx)
	}
	if err != nil {
		return nil, err
	}
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
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	return val, nil
}

func (ctx *Context) evalBinaryBoolean(e *BinaryExpression) (val Value, err error) {
	var left, right Value
	var ok bool
	left, err = e.left.Evaluate(ctx)
	if err == nil {
		right, err = e.right.Evaluate(ctx)
	}
	if err != nil {
		return nil, err
	}
	var l, r bool
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
		// TODO: replace with runtime type error
		unreachable("prevented by prior type check")
	}
	return val, nil
}
