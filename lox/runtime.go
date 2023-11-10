package lox

import "fmt"

type Runtime struct{}

func (r *Runtime) Execute(stmt Statement) error {
	return stmt.Execute()
}

type RuntimeError struct {
	Err error
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s", e.Err)
}

func (e RuntimeError) Unwrap() error {
	return e.Err
}

type InvalidOperationErr struct {
	Operator
	Values []Value
}

func (e InvalidOperationErr) Error() string {
	return fmt.Sprintf("can't apply operator %s to values %v", e.Operator, e.Values)
}

func InvalidOperation(op Operator, vs ...Value) InvalidOperationErr {
	return InvalidOperationErr{Operator: op, Values: vs}
}
