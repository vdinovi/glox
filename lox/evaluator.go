package lox

import (
	"github.com/rs/zerolog/log"
)

type Symbols map[string]Type
type Bindings map[string]Value

type EvaluationContext struct {
	symbols  Symbols
	bindings Bindings
}

func NewEvaluationContext() *EvaluationContext {
	return &EvaluationContext{
		symbols:  make(Symbols),
		bindings: make(Bindings),
	}
}

func (ctx *EvaluationContext) EvaluateUnaryExpression(e UnaryExpression) (Value, Type, error) {
	rightType, typ, err := ctx.symbols.TypeCheckUnaryExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	rightVal, err := e.right.Evaluate(ctx)
	if err != nil {
		return nil, ErrType, err
	}
	var val Value
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
	log.Debug().Msgf("(evaluator) binary expr %s evaluates to %s", e, val)
	return val, typ, err

}

func (ctx *EvaluationContext) EvaluateBinaryExpression(e BinaryExpression) (Value, Type, error) {
	leftType, _, typ, err := ctx.symbols.TypeCheckBinaryExpression(e)
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
	log.Debug().Msgf("(evaluator) binary expr %s evaluates to %s", e, val)
	return val, typ, err
}

func (ctx *EvaluationContext) EvaluateGroupingExpression(e GroupingExpression) (Value, Type, error) {
	_, typ, err := ctx.symbols.TypeCheckGroupingExpression(e)
	if err != nil {
		return nil, ErrType, err
	}
	val, err := e.expr.Evaluate(ctx)
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, ErrType, err
	}
	log.Debug().Msgf("(evaluator) grouping expr %s evaluates to %s", e, val)
	return val, typ, nil
}

func (ctx *EvaluationContext) EvaluateStringExpression(e StringExpression) (Value, error) {
	return ValueString(e), nil
}

func (ctx *EvaluationContext) EvaluateNumericExpression(e NumericExpression) (Value, error) {
	return ValueNumeric(e), nil
}

func (ctx *EvaluationContext) EvaluateBooleanExpression(e BooleanExpression) (Value, error) {
	return ValueBoolean(e), nil
}

func (ctx *EvaluationContext) EvaluateNilExpression(e NilExpression) (Value, error) {
	return Nil, nil
}

func (ctx *EvaluationContext) evalUnaryNumeric(op Operator, right Value) (Value, error) {
	n, ok := right.Unwrap().(float64)
	if ok {
		return nil, NewRuntimeError(NewDowncastError(right, "float64"), op.Line, op.Column)
	}
	switch op.Type {
	case OpSubtract:
		return ValueNumeric(-n), nil
	}
	unreachable("prevented by prior type check")
	return nil, nil
}

func (ctx *EvaluationContext) evalBinaryOperands(left, right Expression) (Value, Value, error) {
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

func (ctx *EvaluationContext) evalBinaryString(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r string
	if l, ok = left.Unwrap().(string); !ok {
		return nil, NewRuntimeError(NewDowncastError(left, "string"), op.Line, op.Column)
	}
	if r, ok = right.Unwrap().(string); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "string"), op.Line, op.Column)
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

func (ctx *EvaluationContext) evalBinaryNumeric(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r float64
	if l, ok = left.Unwrap().(float64); !ok {
		return nil, NewRuntimeError(NewDowncastError(left, "string"), op.Line, op.Column)
	}
	if r, ok = right.Unwrap().(float64); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "string"), op.Line, op.Column)
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

func (ctx *EvaluationContext) evalBinaryBoolean(op Operator, left, right Value) (Value, error) {
	var ok bool
	var l, r bool
	if l, ok = left.Unwrap().(bool); !ok {
		return nil, NewRuntimeError(NewDowncastError(left, "bool"), op.Line, op.Column)
	}
	if r, ok = right.Unwrap().(bool); !ok {
		return nil, NewRuntimeError(NewDowncastError(right, "bool"), op.Line, op.Column)
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
