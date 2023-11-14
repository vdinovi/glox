package lox

import (
	"io"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	runtime *Runtime
	ctx     *Context
}

func NewExecutor(printer io.Writer) *Executor {
	return &Executor{
		runtime: NewRuntime(printer),
		ctx:     NewContext(),
	}
}

func (e *Executor) Execute(stmt Statement) error {
	return stmt.Execute(e)
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
		err = x.runtime.Print(val.String())
	}
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(executor) success: printed %s", val)
	return nil
}

func (x *Executor) ExecuteDeclarationStatement(s DeclarationStatement) error {
	log.Debug().Msgf("(executor) executing %q", s)
	val, err := s.expr.Evaluate(x.ctx)
	if err == nil {
		err = x.ctx.values.Set(s.name, val)
		if err == nil {
			err = x.ctx.types.Set(s.name, val.Type())
		}
	}
	if err != nil {
		log.Error().Msgf("(executor) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(executor) success: %q = %s", s.name, val)
	return nil
}
