package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vdinovi/glox/lox"
)

const usage = `Usage: glox [file]`

func main() {
	if len(os.Args) > 2 {
		lox.Exitln(lox.ExitCodeErr, usage)
	} else if len(os.Args) > 1 {
		execFile(os.Args[1])
	} else {
		repl()
	}
}

func execFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		lox.ExitErr(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	err = exec(rd, filepath.Clean(path))
	if err != nil {
		lox.ExitErr(err)
	}
}

func repl() {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				lox.Exit(lox.ExitCodeOK)
			} else {
				lox.ExitErr(err)
			}
		}

		line = strings.TrimRight(line, "\n")
		rd := strings.NewReader(line)
		err = exec(rd, "<stdin>")
		if err != nil {
			lox.ExitErr(err)
		}
	}
}

func exec(r io.Reader, fname string) error {
	lexer := lox.NewLexer(bufio.NewReader(r))
	lexer.SetFilename(fname)
	tokens, err := lexer.ScanTokens()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
	}
	return nil
}
