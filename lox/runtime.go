package lox

import (
	"fmt"
	"io"
)

type Runtime struct {
	printer io.Writer
}

func NewRuntime(printer io.Writer) *Runtime {
	return &Runtime{
		printer: printer,
	}
}

func (r *Runtime) Print(s string) error {
	_, err := fmt.Fprintln(r.printer, s)
	return err
}
