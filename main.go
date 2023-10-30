package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

	err = exec(f)
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
		err = exec(rd)
		if err != nil {
			lox.ExitErr(err)
		}
	}
}

func exec(r io.Reader) error {
	// data, err := io.ReadAll(r)
	// if err != nil {
	// 	return err
	// }
	tokens, errs := lox.Lex(r)
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	for _, tok := range tokens {
		fmt.Fprintf(os.Stdout, "%v\n", tok)
	}
	return nil
}
