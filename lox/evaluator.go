package lox

import (
	"fmt"

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

func (ctx *EvaluationContext) EvaluateUnaryExpression(e UnaryExpression) (Value, error) {
	var err error
	var right Value
	var typ Type
	if typ, err = ctx.symbols.TypeCheckUnaryExpression(e); err != nil {
		return nil, err
	}
	if right, err = e.right.Evaluate(ctx); err != nil {
		return nil, err
	}
	var val Value
	switch typ {
	case TypeNumeric:
		val, err = ctx.evalUnaryNumeric(e.op.Type, right)
	default:
		err = ExecutorError{fmt.Errorf("unary expression does not support type %s", typ)}
	}
	if err == nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, err
	}
	log.Debug().Msgf("(evaluator) binary expr %s evaluates to %s", e, val)
	return val, err

}

func (ctx *EvaluationContext) EvaluateBinaryExpression(e BinaryExpression) (Value, error) {
	var err error
	var left, right Value
	var typ Type
	if typ, err = ctx.symbols.TypeCheckBinaryExpression(e); err != nil {
		return nil, err
	}
	if left, right, err = ctx.evalBinaryOperands(e.left, e.right); err != nil {
		return nil, err
	}
	var val Value
	switch typ {
	case TypeString:
		val, err = ctx.evalBinaryString(e.op.Type, left, right)
	case TypeNumeric:
		val, err = ctx.evalBinaryNumeric(e.op.Type, left, right)
	case TypeBoolean:
		val, err = ctx.evalBinaryBoolean(e.op.Type, left, right)
	default:
		err = ExecutorError{fmt.Errorf("binary expression does not support type %s", typ)}
	}
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, err
	}
	log.Debug().Msgf("(evaluator) binary expr %s evaluates to %s", e, val)
	return val, err
}

func (ctx *EvaluationContext) EvaluateGroupingExpression(e GroupingExpression) (Value, error) {
	val, err := e.expr.Evaluate(ctx)
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, err
	}
	log.Debug().Msgf("(evaluator) grouping expr %s evaluates to %s", e, val)
	return val, nil
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
		return nil, ExecutorError{DowncastError{right, "float64"}}
	}
	switch op {
	case OpMinus:
		return ValueNumeric(-n), nil
	default:
		return nil, ExecutorError{TypeError{fmt.Errorf("%s can't be applied to value %s", op, right)}}
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
		return nil, ExecutorError{DowncastError{left, "string"}}
	}
	if r, ok = right.Unwrap().(string); !ok {
		return nil, ExecutorError{DowncastError{right, "string"}}
	}
	switch op {
	case OpPlus:
		return ValueString(l + r), nil
	case OpEqualEquals:
		return ValueBoolean(l == r), nil
	case OpNotEquals:
		return ValueBoolean(l != r), nil
	default:
		return nil, ExecutorError{TypeError{fmt.Errorf("%s can't be applied to values %s and %s", op, left, right)}}
	}
}

func (ctx *EvaluationContext) evalBinaryNumeric(op OperatorType, left, right Value) (Value, error) {
	var ok bool
	var l, r float64
	if l, ok = left.Unwrap().(float64); !ok {
		return nil, ExecutorError{DowncastError{left, "float64"}}
	}
	if r, ok = right.Unwrap().(float64); !ok {
		return nil, ExecutorError{DowncastError{right, "float64"}}
	}
	switch op {
	case OpPlus:
		return ValueNumeric(l + r), nil
	case OpMinus:
		return ValueNumeric(l - r), nil
	case OpMultiply:
		return ValueNumeric(l * r), nil
	case OpDivide:
		return ValueNumeric(l / r), nil
	case OpEqualEquals:
		return ValueBoolean(l == r), nil
	case OpNotEquals:
		return ValueBoolean(l != r), nil
	case OpLess:
		return ValueBoolean(l < r), nil
	case OpLessEquals:
		return ValueBoolean(l <= r), nil
	case OpGreater:
		return ValueBoolean(l > r), nil
	case OpGreaterEquals:
		return ValueBoolean(l >= r), nil
	default:
		return nil, ExecutorError{TypeError{fmt.Errorf("%s can't be applied to values %s and %s", op, left, right)}}
	}
}

func (ctx *EvaluationContext) evalBinaryBoolean(op OperatorType, left, right Value) (Value, error) {
	var ok bool
	var l, r bool
	if l, ok = left.Unwrap().(bool); !ok {
		return nil, ExecutorError{DowncastError{left, "bool"}}
	}
	if r, ok = right.Unwrap().(bool); !ok {
		return nil, ExecutorError{DowncastError{right, "bool"}}
	}
	switch op {
	case OpEqualEquals:
		return ValueBoolean(l == r), nil
	case OpNotEquals:
		return ValueBoolean(l != r), nil
	default:
		return nil, ExecutorError{TypeError{fmt.Errorf("%s can't be applied to values %s and %s", op, left, right)}}
	}
}
