package lox

import (
	"fmt"
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
			if ret, ok := err.(ReturnErr); ok {
				err = NewRuntimeError(
					fmt.Errorf("out of place return statement"),
					ret.Position(),
				)
			}
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
	exit := executeEnterEnv(x.ctx, "<block>")
	defer exit()
	for _, stmt := range s.stmts {
		if err := stmt.Execute(x); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConditionalStatement) Execute(x *Executor) error {
	cond, err := s.expr.Evaluate(x)
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
		cond, err := s.expr.Evaluate(x)
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
			cond, err := s.cond.Evaluate(x)
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
			_, err := s.incr.Evaluate(x)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ExpressionStatement) Execute(x *Executor) error {
	_, err := s.expr.Evaluate(x)
	return err
}

func (s *PrintStatement) Execute(x *Executor) error {
	val, err := s.expr.Evaluate(x)
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
	val, err := s.expr.Evaluate(x)
	if err != nil {
		return err
	}
	prev, err := x.ctx.values.Set(s.name, val)
	if prev == nil {
		log.Debug().Msgf("(execute) Env(%s) %s := %s", x.ctx.values.name, s.name, val)
	} else {
		log.Debug().Msgf("(execute) Env(%s) %s = %s (was %s)", x.ctx.values.name, s.name, val, *prev)
	}
	return err
}

func (s *FunctionStatement) Execute(x *Executor) error {
	fn := &UserFunction{
		name:    s.name,
		params:  s.params,
		body:    s.body,
		context: x.ctx.Copy(),
	}
	val := &ValueCallable{
		name: s.name,
		fn:   fn,
	}
	prev, err := x.ctx.values.Set(fn.name, val)
	if err != nil {
		return err
	}
	if prev == nil {
		log.Debug().Msgf("(execute) Env(%s) %s := %s", x.ctx.values.name, s.name, val)
	} else {
		log.Debug().Msgf("(execute) Env(%s) %s = %s (was %s)", x.ctx.values.name, s.name, val, *prev)
	}
	return nil
}

func (s *ReturnStatement) Execute(x *Executor) error {
	val, err := s.expr.Evaluate(x)
	if err != nil {
		return err
	}
	return ReturnErr{val: val, pos: s.Position()}
}

func deparenthesize(s string) string {
	return strings.Trim(s, "\"")
}

func executeEnterEnv(ctx *Context, name string) (exit func()) {
	ctx.PushEnvironment(name)
	log.Debug().Msgf("(execute) ENTER %s", ctx.values.String())
	return func() {
		log.Debug().Msgf("(execute) EXIT %s", ctx.values.String())
		ctx.PopEnvironment()
		log.Debug().Msgf("(execute) ENTER %s", ctx.values.String())
	}
}
