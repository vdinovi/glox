package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) EvaluateUnaryExpression(e *UnaryExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateUnaryExpression")
	var right Value
	right, err = e.right.Evaluate(ctx)
	if err != nil {
		return nil, TypeAny, err
	}
	switch typ := right.Type(); typ {
	case TypeNumeric:
		val, err = ctx.evalUnaryNumeric(e, right)
	default:
		err = NewRuntimeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right.Type()), e.Position())
	}
	if err != nil {
		return nil, TypeAny, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, err

}

func (ctx *Context) EvaluateBinaryExpression(e *BinaryExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateBinaryExpression")
	var left, right Value
	left, err = e.left.Evaluate(ctx)
	if err == nil {
		right, err = e.right.Evaluate(ctx)
	}
	if err != nil {
		return nil, TypeAny, err
	}
	if left.Type() != right.Type() {
		return nil, TypeAny, NewRuntimeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	switch typ := left.Type(); typ {
	case TypeString:
		val, err = ctx.evalBinaryString(e, left, right)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e, left, right)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e, left, right)
	default:
		err = NewRuntimeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type()), e.Position())
	}
	if err != nil {
		return nil, TypeAny, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateGroupingExpression(e *GroupingExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateGroupingExpression")
	typ = e.Type()
	if val, err = e.expr.Evaluate(ctx); err != nil {
		return nil, TypeAny, err
	}
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) EvaluateAssignmentExpression(e *AssignmentExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateAssignmentExpression")
	if val, err = e.right.Evaluate(ctx); err != nil {
		return nil, TypeAny, err
	}
	prev, env := ctx.values.Lookup(e.name)
	if prev == nil {
		return nil, TypeAny, NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
	}
	err = env.Set(e.name, val)
	if err != nil {
		return nil, TypeAny, err
	}
	log.Debug().Msgf("(evaluator) %s <- %s (prev %s)", e.name, val, *prev)
	return val, typ, nil
}

func (ctx *Context) EvaluateVariableExpression(e *VariableExpression) (val Value, typ Type, err error) {
	log.Trace().Msg("EvaluateVariableExpression")
	typ = e.Type()
	v, _ := ctx.values.Lookup(e.name)
	if val != nil {
		err := NewRuntimeError(NewUndefinedVariableError(e.name), e.Position())
		return nil, TypeAny, err
	}
	val = *v
	log.Debug().Msgf("(evaluator) eval(%s) => %s", e, val)
	return val, typ, nil
}

func (ctx *Context) evalUnaryNumeric(e *UnaryExpression, right Value) (val Value, err error) {
	var num ValueNumeric
	var ok bool
	if num, ok = right.(ValueNumeric); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "Numeric"), Position{})
	}
	switch e.op.Type {
	case OpSubtract:
		val, err = num.Negative()
	case OpAdd:
		val = num
	default:
		err = NewInvalidUnaryOperatorForTypeError(e.op.Type, right.Type())
	}
	if err != nil {
		err = NewRuntimeError(err, e.Position())
	}
	return val, err
}

func (ctx *Context) evalBinaryString(e *BinaryExpression, left, right Value) (val Value, err error) {
	var l, r ValueString
	var ok bool
	if l, ok = left.(ValueString); !ok {
		return nil, NewRuntimeError(NewDowncastError(left, "String"), Position{})
	}
	if r, ok = right.(ValueString); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "String"), Position{})
	}
	switch e.op.Type {
	case OpAdd:
		val, err = l.Concat(r)
	default:
		err = NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type())
	}
	if err != nil {
		err = NewRuntimeError(err, e.Position())
	}
	return val, err
}

func (ctx *Context) evalBinaryNumeric(e *BinaryExpression, left, right Value) (val Value, err error) {
	var l, r ValueNumeric
	var ok bool
	if l, ok = left.(ValueNumeric); !ok {
		return nil, NewRuntimeError(NewDowncastError(left, "Numeric"), Position{})
	}
	if r, ok = right.(ValueNumeric); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "Numeric"), Position{})
	}
	switch e.op.Type {
	case OpAdd:
		val, err = l.Add(r)
	case OpSubtract:
		val, err = l.Subtract(r)
	case OpMultiply:
		val, err = l.Multiply(r)
	case OpDivide:
		val, err = l.Divide(r)
	default:
		err = NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type())
	}
	if err != nil {
		err = NewRuntimeError(err, e.Position())
	}
	return val, err
}

func (ctx *Context) evalBinaryBoolean(e *BinaryExpression, left, right Value) (val Value, err error) {
	switch e.op.Type {
	default:
		err = NewInvalidBinaryOperatorForTypeError(e.op.Type, left.Type(), right.Type())
	}
	//if err != nil {}
	err = NewRuntimeError(err, e.Position())
	return val, err
}
