package lox

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Function interface {
	fmt.Stringer
	Execute(*Executor, string, ...Value) (Value, error)
	Context() *Context
}

type UserFunction struct {
	name    string
	params  []string
	body    []Statement
	context *Context
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

func (f *UserFunction) Context() *Context {
	return f.context
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

func (f *UserFunction) Execute(x *Executor, name string, args ...Value) (Value, error) {
	log.Debug().Msgf("(execute) executing fn %s with %v", f.name, args)
	if len(args) != len(f.params) {
		return nil, NewArityMismatchError(f.Arity(), len(args))
	}
	prevCtx := x.ctx
	x = &Executor{
		printer: x.printer,
		runtime: x.runtime,
		ctx:     f.context,
	}
	exit := executeEnterEnv(x.ctx, name)
	defer func() {
		exit()
		log.Debug().Msgf("(execute) ENTER %s", prevCtx.values.String())
	}()
	for i, arg := range args {
		if _, err := x.ctx.values.Set(f.params[i], arg); err != nil {
			return nil, err
		}
	}
	for _, s := range f.body {
		if err := s.Execute(x); err != nil {
			if ret, ok := err.(ReturnErr); ok {
				return ret.val, nil
			}
			return nil, err
		}
	}
	return Nil, nil
}
