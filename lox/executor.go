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

func (x *Executor) TypeCheckProgram(stmts []Statement) error {
	for _, stmt := range stmts {
		if err := x.TypeCheck(stmt); err != nil {
			log.Error().Msgf("(typechecker) error in statement %q: %s", stmt, err)
			return err
		}
	}
	return nil
}

func (x *Executor) TypeCheck(stmt Statement) error {
	return stmt.TypeCheck(x.ctx)
}

func (x *Executor) ExecuteProgram(stmts []Statement) error {
	for _, stmt := range stmts {
		if err := x.Execute(stmt); err != nil {
			log.Error().Msgf("(executor) error in statement %q: %s", stmt, err)
			return err
		}
	}
	return nil
}

func (x *Executor) Execute(stmt Statement) error {
	return stmt.Execute(x)
}

func (x *Executor) ExecuteBlockStatement(s *BlockStatement) error {
	log.Trace().Msg("ExecuteBlockStatement")
	log.Debug().Msgf("(executor) enter scope")
	x.ctx.PushEnvironment()
	defer x.ctx.PopEnvironment()
	defer log.Debug().Msgf("(executor) exit scope")
	for _, stmt := range s.stmts {
		if err := stmt.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (x *Executor) ExecuteExpressionStatement(s *ExpressionStatement) error {
	log.Trace().Msg("ExecuteExpressionStatement")
	_, err := s.expr.Evaluate(x.ctx)
	return err
}

func (x *Executor) ExecutePrintStatement(s *PrintStatement) error {
	log.Trace().Msg("ExecutePrintStatement")
	val, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	err = x.runtime.Print(val.String())
	if err != nil {
		return NewRuntimeError(err, s.Position())
	}
	return nil
}

func (x *Executor) ExecuteDeclarationStatement(s *DeclarationStatement) error {
	log.Trace().Msg("ExecuteDeclarationStatement")
	log.Debug().Msgf("(executor) executing %q", s)
	val, err := s.expr.Evaluate(x.ctx)
	if err == nil {
		err = x.ctx.values.Set(s.name, val)
	}
	return err
}
