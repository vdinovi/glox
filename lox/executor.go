package lox

import (
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
