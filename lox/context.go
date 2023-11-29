package lox

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Phase string

const (
	PhaseInit      Phase = "init"
	PhaseLex       Phase = "lex"
	PhaseParse     Phase = "parse"
	PhaseTypecheck Phase = "typecheck"
	PhaseExecute   Phase = "execute"
)

func (p Phase) String() string {
	return string(p)
}

type Context struct {
	phase   Phase
	env     *Env
	runtime *Runtime
	printer Printer
	funcs   []Function
}

func NewContext(w io.Writer) *Context {
	return &Context{
		phase:   PhaseInit,
		env:     NewEnv("root", nil),
		runtime: NewRuntime(w),
		printer: &defaultPrinter,
		funcs:   make([]Function, 0),
	}
}

func (ctx *Context) Phase() Phase {
	return ctx.phase
}

func (ctx *Context) Copy() Context {
	return *ctx
}

func (ctx *Context) PushEnv(name string) (pop func()) {
	prev := ctx.env
	pop = func() {
		ctx.env = prev
	}
	ctx.env = NewEnv(name, prev)
	return pop
}

func (ctx *Context) StartPhase(phase Phase) (restore func()) {
	p := ctx.phase
	restore = func() {
		ctx.phase = p
	}
	ctx.phase = phase
	return restore
}

func (ctx *Context) debug() string {
	var sb strings.Builder
	sb.WriteString("Context:\n")
	fmt.Fprintf(&sb, "\tphase: %s\n", ctx.Phase())
	sb.WriteString("\tfuncs:\n")
	for _, f := range ctx.funcs {
		fmt.Fprintf(&sb, "\t\t%s\n", f.String())
	}
	envs := []*Env{}
	for env := ctx.env; env != nil; env = env.parent {
		envs = append(envs, env)
	}
	slices.Reverse(envs)
	sb.WriteString("\tenv:\n")
	for i, env := range envs {
		fmt.Fprintf(&sb, "\t\t%s%s\n", strings.Repeat("\t", i), env.String())
	}
	return sb.String()
}
