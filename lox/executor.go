package lox

import (
	"io"
	"strings"

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
	log.Debug().Msgf("(executor) executing %q", stmt)
	return stmt.Execute(x)
}

func (x *Executor) ExecuteBlockStatement(s *BlockStatement) error {
	x.ctx.PushEnvironment()
	log.Debug().Msgf("(executor) entering scope {%s}", x.ctx.values.String())
	defer func() {
		x.ctx.PopEnvironment()
		log.Debug().Msgf("(executor) entering scope {%s}", x.ctx.values.String())
	}()
	defer log.Debug().Msgf("(executor) exit scope")
	for _, stmt := range s.stmts {
		if err := stmt.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (x *Executor) ExecuteConditionalStatement(s *ConditionalStatement) error {
	cond, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	if cond.Truthy() {
		if err := s.thenBranch.Execute(x); err != nil {
			return err
		}
	} else if s.elseBranch != nil {
		if err := s.elseBranch.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (x *Executor) ExecuteWhileStatement(s *WhileStatement) error {
	for {
		cond, err := s.expr.Evaluate(x.ctx)
		if err != nil {
			return err
		}
		if !cond.Truthy() {
			log.Debug().Msg("(executor) break loop")
			break
		}
		if err := s.body.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (x *Executor) ExecuteForStatement(s *ForStatement) error {
	if s.init != nil {
		s.init.Execute(x)
	}
	for {
		if s.cond != nil {
			cond, err := s.cond.Evaluate(x.ctx)
			if err != nil {
				return err
			}
			if !cond.Truthy() {
				log.Debug().Msg("(executor) break loop")
				break
			}
		}
		if err := s.body.Execute(x); err != nil {
			return err
		}
		if s.incr != nil {
			_, err := s.incr.Evaluate(x.ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *Executor) ExecuteExpressionStatement(s *ExpressionStatement) error {
	_, err := s.expr.Evaluate(x.ctx)
	return err
}

func (x *Executor) ExecutePrintStatement(s *PrintStatement) error {
	val, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	err = x.runtime.Print(deparenthesize(val.String()))
	if err != nil {
		return NewRuntimeError(err, s.Position())
	}
	return nil
}

func (x *Executor) ExecuteDeclarationStatement(s *DeclarationStatement) error {
	val, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	prev := x.ctx.values.Get(s.name, nil)
	err = x.ctx.values.Set(s.name, val)
	if err != nil {
		return err
	}
	if prev == nil {
		log.Debug().Msgf("(executor) initialized %s to %s", s.name, val)
	} else {
		log.Debug().Msgf("(executor) %s <- %s (prev %s)", s.name, val, prev)
	}
	return err
}

func deparenthesize(s string) string {
	return strings.Trim(s, "\"")
}
