package lox

import (
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	runtime *Runtime
	ctx     *Context
	printer Printer
}

func NewExecutor(printer io.Writer) *Executor {
	return &Executor{
		runtime: NewRuntime(printer),
		ctx:     NewContext(),
		printer: &defaultPrinter,
	}
}

type Executable interface {
	Execute(*Executor) error
}

func (x *Executor) Execute(elems []Statement) error {
	for _, elem := range elems {
		log.Debug().Msgf("(execute) executing %s", elem)
		if err := elem.Execute(x); err != nil {
			log.Error().Msgf("(execute) error in %q: %s", elem, err)
			return err
		}
	}
	return nil
}

func (x *Executor) Typecheck(elems []Statement) error {
	for _, elem := range elems {
		log.Debug().Msgf("(typecheck) checking %s", elem)
		if err := elem.Typecheck(x.ctx); err != nil {
			log.Error().Msgf("(typecheck) error in %q: %s", elem, err)
			return err
		}
	}
	return nil
}

func (s *BlockStatement) Execute(x *Executor) error {
	x.ctx.PushEnvironment()
	log.Debug().Msgf("(execute) enter %s", x.ctx.values.String())
	defer func() {
		x.ctx.PopEnvironment()
		log.Debug().Msgf("(execute) enter %s", x.ctx.values.String())
	}()
	for _, stmt := range s.stmts {
		if err := stmt.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConditionalStatement) Execute(x *Executor) error {
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

func (s *WhileStatement) Execute(x *Executor) error {
	for {
		cond, err := s.expr.Evaluate(x.ctx)
		if err != nil {
			return err
		}
		if !cond.Truthy() {
			log.Debug().Msg("(execute) break loop")
			break
		}
		if err := s.body.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (s *ForStatement) Execute(x *Executor) error {
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
				log.Debug().Msg("(execute) break loop")
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

func (s *ExpressionStatement) Execute(x *Executor) error {
	_, err := s.expr.Evaluate(x.ctx)
	return err
}

func (s *PrintStatement) Execute(x *Executor) error {
	val, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	str, err := val.Print(x.printer)
	if err != nil {
		return err
	}
	err = x.runtime.Print(deparenthesize(str))
	if err != nil {
		return NewRuntimeError(err, s.Position())
	}
	return nil
}

func (s *DeclarationStatement) Execute(x *Executor) error {
	val, err := s.expr.Evaluate(x.ctx)
	if err != nil {
		return err
	}
	prev, err := x.ctx.values.Set(s.name, val)
	if prev == nil {
		log.Debug().Msgf("(execute) (%d) %s := %s", x.ctx.values.depth, s.name, val)
	} else {
		log.Debug().Msgf("(execute) (%d) %s = %s (was %s)", x.ctx.values.depth, s.name, val, *prev)
	}
	return err
}

func deparenthesize(s string) string {
	return strings.Trim(s, "\"")
}
