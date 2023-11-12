package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vdinovi/glox/lox"
	"github.com/vdinovi/glox/tool"
)

func main() {
	err := tool.Setup()
	if err != nil {
		lox.ExitErr(err)
	}
	if flag.NArg() != 1 {
		lox.Exit(lox.ExitCodeErr)
	}
	path := filepath.Clean(flag.Arg(0))
	file, err := os.Open(path)
	if err != nil {
		lox.ExitErr(err)
	}
	err = process(file, os.Stdout)
	file.Close()
	if err != nil {
		lox.ExitErr(err)
	}
	lox.Exit(lox.ExitCodeOK)
}

func process(r io.Reader, w io.Writer) error {
	lexer, err := lox.NewLexer(bufio.NewReader(r))
	if err != nil {
		return err
	}
	tokens, err := lexer.Scan()
	if err != nil {
		return err
	}
	parser := lox.NewParser(tokens)
	program, err := parser.Parse()
	if err != nil {
		return err
	}
	for _, stmt := range program {
		_, err := fmt.Fprintln(w, stmt.String())
		if err != nil {
			return err
		}
	}
	return nil
}
