package lox

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

type Terminal struct {
	terminal *term.Terminal
	restore  func()
}

func NewTerminal(in *os.File, prompt string) (*Terminal, error) {
	fd := int(in.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	var restoreOnce sync.Once
	restore := func() {
		restoreOnce.Do(func() {
			term.Restore(fd, state)
		})
	}
	return &Terminal{
		terminal: term.NewTerminal(in, prompt),
		restore:  restore,
	}, nil
}

func (t *Terminal) Close() {
	t.restore()
	t.terminal = nil
}

func (t *Terminal) ReadLine() (string, error) {
	if t.terminal == nil {
		return "", TerminalClosedError{"ReadLine"}
	}
	return t.terminal.ReadLine()
}

func (t *Terminal) Write(buf []byte) (int, error) {
	if t.terminal == nil {
		return -1, TerminalClosedError{"Write"}
	}
	return t.terminal.Write(buf)
}

func (t *Terminal) WriteError(err error) error {
	if t.terminal == nil {
		return TerminalClosedError{"WriteError"}
	}
	data := t.terminal.Escape.Yellow
	data = append(data, []byte(err.Error())...)
	data = append(data, t.terminal.Escape.Reset...)
	data = append(data, '\n')
	_, e := t.terminal.Write(data)
	return e
}

func (t *Terminal) WriteString(s string) error {
	if t.terminal == nil {
		return TerminalClosedError{"WriteString"}
	}
	_, e := t.terminal.Write([]byte(s + "\n"))
	return e
}

type TerminalClosedError struct {
	operation string
}

func (e TerminalClosedError) Error() string {
	return fmt.Sprintf("attempted %s on a closed terminal", e.operation)
}
