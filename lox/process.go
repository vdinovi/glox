package lox

import (
	"fmt"
	"os"
)

type ExitCode int

const (
	ExitCodeOK  ExitCode = 0
	ExitCodeErr ExitCode = 1
)

func Exit(code ExitCode) {
	os.Exit(int(code))
}

func Exitf(code ExitCode, format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	Exit(code)
}

func Exitln(code ExitCode, msg string) {
	fmt.Fprintln(os.Stderr, msg)
	Exit(code)
}

func ExitErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
	Exit(ExitCodeErr)
}

func unreachable(note string) {
	panic("unreachable: " + note)
}
