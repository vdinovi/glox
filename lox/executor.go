package lox

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	runtime *Runtime
	ctx     *EvaluationContext
}

func NewExecutor(printer io.Writer) *Executor {
	return &Executor{
		runtime: NewRuntime(printer),
		ctx:     NewEvaluationContext(),
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

func (x *Executor) ExecuteExpressionStatement(s ExpressionStatement) error {
	log.Debug().Msgf("(executor) executing %q", s)
	_, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msg("(executor) success")
	return nil
}

func (x *Executor) ExecutePrintStatement(s PrintStatement) error {
	log.Debug().Msgf("(executor) executing %q", s)
	val, err := s.expr.Evaluate(x.ctx)
	if err == nil {
		err = x.Print(val)
	}
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msg("(executor) success")
	return nil
}
