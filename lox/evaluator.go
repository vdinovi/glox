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
		val, err = ctx.evalUnaryNumeric(e.op.Type, rightVal)
	default:
		// Handled by prior type check
		panic("unreachable")
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
		val, err = ctx.evalBinaryString(e.op.Type, leftVal, rightVal)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e.op.Type, leftVal, rightVal)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e.op.Type, leftVal, rightVal)
	default:
		// Handled by prior type check
		panic("unreachable")
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

func (ctx *EvaluationContext) evalUnaryNumeric(op OperatorType, right Value) (Value, error) {
	n, ok := right.Unwrap().(float64)
	if ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(right, "float64"), 0, 0)
	}
	switch op {
	case OpSubtract:
		return ValueNumeric(-n), nil
	default:
		// Handled by prior type check
		panic("unreachable")
	}
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

func (ctx *EvaluationContext) evalBinaryString(op OperatorType, left, right Value) (Value, error) {
	var ok bool
	var l, r string
	if l, ok = left.Unwrap().(string); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(left, "string"), 0, 0)
	}
	if r, ok = right.Unwrap().(string); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(right, "string"), 0, 0)
	}
	switch op {
	case OpAdd:
		return ValueString(l + r), nil
	case OpEqualTo:
		return ValueBoolean(l == r), nil
	case OpNotEqualTo:
		return ValueBoolean(l != r), nil
	default:
		// Handled by prior type check
		panic("unreachable")
	}
}

func (ctx *EvaluationContext) evalBinaryNumeric(op OperatorType, left, right Value) (Value, error) {
	var ok bool
	var l, r float64
	if l, ok = left.Unwrap().(float64); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(left, "string"), 0, 0)
	}
	if r, ok = right.Unwrap().(float64); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(right, "string"), 0, 0)
	}
	switch op {
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
		// Handed by prior type check
		panic("unreachable")
	}
}

func (ctx *EvaluationContext) evalBinaryBoolean(op OperatorType, left, right Value) (Value, error) {
	var ok bool
	var l, r bool
	if l, ok = left.Unwrap().(bool); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(left, "bool"), 0, 0)
	}
	if r, ok = right.Unwrap().(bool); !ok {
		// TODO: feed line and column
		return nil, NewRuntimeError(NewDowncastError(right, "bool"), 0, 0)
	}
	switch op {
	case OpEqualTo:
		return ValueBoolean(l == r), nil
	case OpNotEqualTo:
		return ValueBoolean(l != r), nil
	default:
		// Handed by prior type check
		panic("unreachable")
	}
}
