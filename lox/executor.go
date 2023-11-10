package lox

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	runtime Runtime
}

func NewExecutor(printer io.Writer) *Executor {
	return &Executor{
		runtime: Runtime{printer},
	}
}

func (e *Executor) Execute(stmt Statement) error {
	return stmt.Execute(e)
}

func (e *Executor) Print(val Value) error {
	return e.runtime.Print(val.String())
}

type ExecutorError struct {
	Err error
}

func (e ExecutorError) Error() string {
	return fmt.Sprintf("Executor Error: %s", e.Err)
}

func (e ExecutorError) Unwrap() error {
	return e.Err
}

type DowncastError struct {
	Value
	To string
}

func (e DowncastError) Error() string {
	return fmt.Sprintf("failed to cast %s to %s", e.Value, e.To)
}

func (*Executor) ExecuteExpressionStatement(s ExpressionStatement) error {
	log.Debug().Msgf("(executor) executing %q", s)
	_, err := s.expr.Evaluate()
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msg("(executor) success")
	return nil
}

func (e *Executor) ExecutePrintStatement(s PrintStatement) error {
	log.Debug().Msgf("(executor) executing %q", s)
	val, err := s.expr.Evaluate()
	if err == nil {
		err = e.Print(val)
	}
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msg("(executor) success")
	return nil
}

func EvaluateUnaryExpression(e UnaryExpression) (Value, error) {
	var err error
	var right Value
	var typ Type
	if typ, err = TypeCheckUnaryExpression(e); err != nil {
		return nil, err
	}
	if right, err = e.right.Evaluate(); err != nil {
		return nil, err
	}
	var val Value
	switch typ {
	case TypeNumeric:
		val, err = evalUnaryNumeric(e.op.Type, right)
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

func EvaluateBinaryExpression(e BinaryExpression) (Value, error) {
	var err error
	var left, right Value
	var typ Type
	if typ, err = TypeCheckBinaryExpression(e); err != nil {
		return nil, err
	}
	if left, right, err = evalBinaryOperands(e.left, e.right); err != nil {
		return nil, err
	}
	var val Value
	switch typ {
	case TypeString:
		val, err = evalBinaryString(e.op.Type, left, right)
	case TypeNumeric:
		val, err = evalBinaryNumeric(e.op.Type, left, right)
	case TypeBoolean:
		val, err = evalBinaryBoolean(e.op.Type, left, right)
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

func EvaluateGroupingExpression(e GroupingExpression) (Value, error) {
	val, err := e.expr.Evaluate()
	if err != nil {
		log.Error().Msgf("(evaluator) error in %q: %s", e, err)
		return nil, err
	}
	log.Debug().Msgf("(evaluator) grouping expr %s evaluates to %s", e, val)
	return val, nil
}

func evalUnaryNumeric(op OperatorType, right Value) (Value, error) {
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

func evalBinaryOperands(left, right Expression) (Value, Value, error) {
	lv, err := left.Evaluate()
	if err != nil {
		return nil, nil, err
	}
	rv, err := right.Evaluate()
	if err != nil {
		return nil, nil, err
	}
	return lv, rv, nil
}

func evalBinaryString(op OperatorType, left, right Value) (Value, error) {
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

func evalBinaryNumeric(op OperatorType, left, right Value) (Value, error) {
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

func evalBinaryBoolean(op OperatorType, left, right Value) (Value, error) {
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
