package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vdinovi/glox/lox"
)

const usagef = `Usage: %s [file...]
       starts a repl if no files are provided.
`

func main() {
	verbose := flag.Bool("verbose", false, "emit logging to stderr")
	displayTokens := flag.Bool("tokens", true, "display tokens")
	displayAST := flag.Bool("ast", true, "display ast")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usagef, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	fmt.Println(flag.NArg())
	if flag.NArg() == 0 {
		repl(*verbose, *displayTokens, *displayAST)
	} else {
		for _, path := range os.Args[1:] {
			execFile(filepath.Clean(path), *verbose, *displayTokens, *displayAST)
		}
	}
}

func execFile(path string, verbose, displayTokens, displayAST bool) {
	f, err := os.Open(path)
	if err != nil {
		lox.ExitErr(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	err = exec(rd, verbose, displayTokens, displayAST)
	if err != nil {
		lox.ExitErr(err)
	}
}

func repl(verbose, displayTokens, displayAST bool) {
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
		err = exec(rd, verbose, displayTokens, displayAST)
		if err != nil {
			lox.ExitErr(err)
		}
	}
}

func exec(r io.Reader, _, displayTokens, displayAST bool) error {
	lexer, err := lox.NewLexer(bufio.NewReader(r))
	if err != nil {
		return err
	}
	tokens, err := lexer.Scan()
	if err != nil {
		return err
	}
	if displayTokens {
		fmt.Println("=== Tokens ===")
		for _, token := range tokens {
			fmt.Printf("%+v\n", token)
		}
	}
	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()
	if err != nil {
		return err
	}
	if displayAST {
		fmt.Println("=== AST ===")
		fmt.Println(expr.String())
	}

	t, err := expr.Type()
	if err != nil {
		return err
	}
	fmt.Println("=== Type ===")
	fmt.Println(t)

	runtime := lox.Runtime{}
	result, _, err := runtime.Evaluate(expr)
	if err != nil {
		return err
	}
	fmt.Println("=== Eval ===")
	fmt.Println(result)

	return nil
}
