package lox

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog/log"
)

type Runtime struct {
	writer io.Writer
	funcs  map[string]Function
}

func NewRuntime(w io.Writer) *Runtime {
	r := &Runtime{
		writer: w,
		funcs:  make(map[string]Function, 1),
	}
	r.defun("clock", clock)
	r.defun("sleep", sleep)
	return r
}

func (r *Runtime) defun(name string, fn func(...Value) (Value, error)) {
	r.funcs[name] = &BuiltinFunction{name: name, exec: fn}
}

func (r *Runtime) Function(name string) Function {
	if fn, ok := r.funcs[name]; ok {
		return fn
	}
	return nil
}

func (r *Runtime) Print(s string) error {
	_, err := fmt.Fprintln(r.writer, s)
	return err
}

func clock(args ...Value) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("clock expects no arguments, but got %d", len(args))
	}
	return ValueNumeric(time.Now().Unix()), nil
}

func sleep(args ...Value) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sleep expects one argument, but got %d", len(args))
	}
	n, ok := args[0].(ValueNumeric)
	if !ok {
		return nil, fmt.Errorf("sleep expects one numeric argument, but got %s", args[0].Type())
	}
	secs := time.Duration(n) * time.Second
	log.Debug().Msgf("(runtime) sleeping for %v", secs)
	time.Sleep(secs)
	return Nil, nil
}
