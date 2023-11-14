package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) EvaluateUnaryExpression(e *UnaryExpression) (val Value, typ Type, err error) {
	rightType, typ, err := ctx.TypeCheckUnaryExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	rightVal, err := e.right.Evaluate(ctx)
	if err != nil {
		return nil, ErrType, err
	}
	switch rightType {
	case TypeNumeric:
		val, err = ctx.evalUnaryNumeric(e.op, rightVal)
	default:
		unreachable("prevented by prior type check")
	}
	if err == nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) binary expr %q evaluates to %s", e, val)
	return val, typ, err

}

func (ctx *Context) EvaluateBinaryExpression(e *BinaryExpression) (Value, Type, error) {
	leftType, _, typ, err := ctx.TypeCheckBinaryExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	leftVal, rightVal, err := ctx.evalBinaryOperands(e.left, e.right)
	if err != nil {
		return nil, ErrType, err
	}
	var val Value
	switch leftType {
	case TypeString:
		val, err = ctx.evalBinaryString(e.op, leftVal, rightVal)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e.op, leftVal, rightVal)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e.op, leftVal, rightVal)
	default:
		unreachable("prevented by prior type check")
	}
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) binary expr %q evaluates to %s", e, val)
	return val, typ, err
}

func (ctx *Context) EvaluateGroupingExpression(e *GroupingExpression) (Value, Type, error) {
	_, typ, err := ctx.TypeCheckGroupingExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	val, err := e.expr.Evaluate(ctx)
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) grouping expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateAssignmentExpression(e *AssignmentExpression) (val Value, typ Type, err error) {
	_, typ, err = ctx.TypeCheckAssignmentExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	val, err = e.right.Evaluate(ctx)
	if err == nil {
		err = ctx.values.Set(e.name, val)
	}
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) assignment expr %q evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateVariableExpression(e *VariableExpression) (Value, Type, error) {
	val := ctx.values.Lookup(e.name)
	if val == nil {
		err := NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) variable expr %q evaluates to %s", e, *val)
	return *val, (*val).Type(), nil
}

func (ctx *Context) evalUnaryNumeric(op Operator, right Value) (Value, error) {
	n, ok := right.Unwrap().(float64)
	if ok {
		return nil, NewRuntimeError(NewDowncastError(right, "float64"), Position{})
	}
	switch op.Type {
	case OpSubtract:
		return ValueNumeric(-n), nil
	}
	unreachable("prevented by prior type check")
	return nil, nil
}

func (ctx *Context) evalBinaryOperands(left, right Expression) (Value, Value, error) {
	lv, err := left.Evaluate(ctx)
	if err != nil {
		return nil, nil, err
	}
	rv, err := right.Evaluate(ctx)
	if err != nil {
		return nil, nil, err
	}
	return lv, rv, nil
}

func (ctx *Context) evalBinaryString(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r string
	if l, ok = left.Unwrap().(string); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(left, "string"), Position{})
	}
	if r, ok = right.Unwrap().(string); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(right, "string"), Position{})
	}
	switch op.Type {
	case OpAdd:
		return ValueString(l + r), nil
	case OpEqualTo:
		return ValueBoolean(l == r), nil
	case OpNotEqualTo:
		return ValueBoolean(l != r), nil
	}
	unreachable("prevented by prior type check")
	return nil, nil
}

func (ctx *Context) evalBinaryNumeric(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r float64
	if l, ok = left.Unwrap().(float64); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(left, "string"), Position{})
	}
	if r, ok = right.Unwrap().(float64); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(right, "string"), Position{})
	}
	switch op.Type {
	case OpAdd:
		return ValueNumeric(l + r), nil
	case OpSubtract:
		return ValueNumeric(l - r), nil
	case OpMultiply:
		return ValueNumeric(l * r), nil
	case OpDivide:
		return ValueNumeric(l / r), nil
	case OpEqualTo:
		return ValueBoolean(l == r), nil
	case OpNotEqualTo:
		return ValueBoolean(l != r), nil
	case OpLessThan:
		return ValueBoolean(l < r), nil
	case OpLessThanOrEqualTo:
		return ValueBoolean(l <= r), nil
	case OpGreaterThan:
		return ValueBoolean(l > r), nil
	case OpGreaterThanOrEqualTo:
		return ValueBoolean(l >= r), nil
	default:
	}
	unreachable("prevented by prior type check")
	return nil, nil
}

func (ctx *Context) evalBinaryBoolean(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r bool
	if l, ok = left.Unwrap().(bool); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(left, "bool"), Position{})
	}
	if r, ok = right.Unwrap().(bool); !ok {
		// TODO
		return nil, NewRuntimeError(NewDowncastError(right, "bool"), Position{})
	}
	switch op.Type {
	case OpEqualTo:
		return ValueBoolean(l == r), nil
	case OpNotEqualTo:
		return ValueBoolean(l != r), nil
	default:
	}
	unreachable("prevented by prior type check")
	return nil, nil
}
