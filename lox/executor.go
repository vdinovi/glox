package lox

import (
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

type Executor struct {
	ctx *Context
}

func NewExecutor(w io.Writer) *Executor {
	return &Executor{
		ctx: NewContext(w),
	}
}

type Executable interface {
	Execute(*Context) error
}

func Execute(ctx *Context, elems []Statement) error {
	restore := ctx.StartPhase(PhaseExecute)
	defer restore()
	for _, elem := range elems {
		log.Debug().Msgf("(%s) executing %s", ctx.Phase(), elem)
		if err := elem.Execute(ctx); err != nil {
			if ret, ok := err.(ReturnErr); ok {
				err = NewRuntimeError(
					fmt.Errorf("out of place return statement"),
					ret.Position(),
				)
			}
			log.Error().Msgf("(%s) error in %q: %s", ctx.Phase(), elem, err)
			return err
		}
	}
	return nil
}

func (s *BlockStatement) Execute(ctx *Context) error {
	exit := debugEnterEnv(ctx, "<block>")
	defer exit()
	for _, stmt := range s.stmts {
		if err := stmt.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConditionalStatement) Execute(ctx *Context) error {
	cond, err := s.expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	if cond.Truthy() {
		if err := s.thenBranch.Execute(ctx); err != nil {
			return err
		}
	} else if s.elseBranch != nil {
		if err := s.elseBranch.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *WhileStatement) Execute(ctx *Context) error {
	for {
		cond, err := s.expr.Evaluate(ctx)
		if err != nil {
			return err
		}
		if !cond.Truthy() {
			log.Debug().Msgf("(%s) break loop", ctx.Phase())
			break
		}
		if err := s.body.Execute(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *ForStatement) Execute(ctx *Context) error {
	if s.init != nil {
		s.init.Execute(ctx)
	}
	for {
		if s.cond != nil {
			cond, err := s.cond.Evaluate(ctx)
			if err != nil {
				return err
			}
			if !cond.Truthy() {
				log.Debug().Msgf("(%s) break loop", ctx.Phase())
				break
			}
		}
		if err := s.body.Execute(ctx); err != nil {
			return err
		}
		if s.incr != nil {
			_, err := s.incr.Evaluate(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ExpressionStatement) Execute(ctx *Context) error {
	_, err := s.expr.Evaluate(ctx)
	return err
}

func (s *PrintStatement) Execute(ctx *Context) error {
	val, err := s.expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	str, err := val.Print(ctx.printer)
	if err != nil {
		return err
	}
	err = ctx.runtime.Print(deparenthesize(str))
	if err != nil {
		return NewRuntimeError(err, s.Position())
	}
	return nil
}

func (s *DeclarationStatement) Execute(ctx *Context) error {
	val, err := s.expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	return debugSetValue(ctx.Phase(), ctx.env, s.name, val)
}

func (s *FunctionDefinitionStatement) Execute(ctx *Context) error {
	fn := &UserFunction{
		name:   s.name,
		params: s.params,
		body:   s.body,
		env:    ctx.env,
	}
	ctx.funcs = append(ctx.funcs, fn)
	log.Debug().Msgf("(%s) created user func %s(...) { ... }", ctx.Phase(), s.name)
	val := &ValueCallable{
		name: s.name,
		fn:   fn,
	}
	return debugSetValue(ctx.Phase(), ctx.env, fn.name, val)
}

func (s *ReturnStatement) Execute(ctx *Context) error {
	val, err := s.expr.Evaluate(ctx)
	if err != nil {
		return err
	}
	return ReturnErr{val: val, pos: s.Position()}
}

func deparenthesize(s string) string {
	return strings.Trim(s, "\"")
}
