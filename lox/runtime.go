package lox

import (
	"fmt"
	"io"
)

type Runtime struct {
	writer io.Writer
}

func NewRuntime(w io.Writer) *Runtime {
	return &Runtime{
		writer: w,
	}
}

func (r *Runtime) Print(s string) error {
	_, err := fmt.Fprintln(r.writer, s)
	return err
}
