package lox

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Function interface {
	fmt.Stringer
	Execute(*Context, string, ...Value) (Value, error)
}

type UserFunction struct {
	name   string
	params []string
	body   []Statement
	env    *Env
}

func (f *UserFunction) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "fun %s(", f.name)
	for i, param := range f.params {
		if i+1 == len(f.params) {
			fmt.Fprintf(&sb, "%s) { ", param)
		} else {
			fmt.Fprintf(&sb, "%s, ", param)
		}
	}
	for _, stmt := range f.body {
		s, err := stmt.Print(&defaultPrinter)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(&sb, "%s ", s)
	}
	fmt.Fprint(&sb, "}")
	return sb.String()
}

func (f *UserFunction) Arity() int {
	return len(f.params)
}

// Hijacking the err return for return handling
type ReturnErr struct {
	val Value
	pos Position
}

func (e ReturnErr) Value() Value {
	return e.val
}

func (e ReturnErr) Position() Position {
	return e.pos
}

func (e ReturnErr) Error() string {
	return "return"
}

func (f *UserFunction) Execute(ctx *Context, name string, args ...Value) (Value, error) {
	stringArgs := make([]string, len(args))
	for i, arg := range args {
		stringArgs[i] = arg.String()
	}
	log.Debug().Msgf("(execute) executing fn %s with (%s)", f.name, strings.Join(stringArgs, ", "))
	if len(args) != len(f.params) {
		return nil, NewArityMismatchError(f.Arity(), len(args))
	}

	if ctx.env != f.env {
		prevEnv := ctx.env
		defer func() {
			ctx.env = prevEnv
		}()
		ctx.env = f.env
		log.Debug().Msgf("(%s) CHANGE %s -> %s", ctx.Phase(), prevEnv, ctx.env)
	}
	exit := debugEnterEnv(ctx, name)
	defer exit()
	for i, arg := range args {
		if err := debugSetValue(ctx.Phase(), ctx.env, f.params[i], arg); err != nil {
			return nil, err
		}
	}
	for _, s := range f.body {
		if err := s.Execute(ctx); err != nil {
			if ret, ok := err.(ReturnErr); ok {
				return ret.val, nil
			}
			return nil, err
		}
	}
	return Nil, nil
}
