package lox

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type Function interface {
	fmt.Stringer
	Execute(*Executor, ...Value) (Value, error)
}

type UserFunction struct {
	name   string
	params []string
	body   []Statement
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

func (f *UserFunction) Execute(x *Executor, args ...Value) (Value, error) {
	log.Debug().Msgf("(execute) executing fn %s with %v", f.name, args)
	if len(args) != len(f.params) {
		return nil, NewArityMismatchError(f.Arity(), len(args))
	}
	exit := x.ctx.Enter("execute")
	defer exit()
	for i, arg := range args {
		if _, err := x.ctx.values.Set(f.params[i], arg); err != nil {
			return nil, err
		}
	}
	for _, s := range f.body {
		// TODO: return value from return statement
		if err := s.Execute(x); err != nil {
			return nil, err
		}
	}
	return Nil, nil
}
