package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vdinovi/glox/lox"
)

const usagef = `Usage: %s [file...]
       starts a repl if no files are provided.
`

func main() {
	debug := flag.Bool("debug", false, "debug logging")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usagef, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if flag.NArg() == 0 {
		repl()
	} else {
		for _, path := range flag.Args() {
			execFile(filepath.Clean(path))
		}
	}
}

func execFile(path string) {
	log.Debug().Msgf("(main) executing %s", path)
	f, err := os.Open(path)
	if err != nil {
		lox.ExitErr(err)
	}
	defer f.Close()

	x := lox.NewExecutor(os.Stdout)
	rd := bufio.NewReader(f)
	err = exec(x, rd)
	if err != nil {
		lox.ExitErr(err)
	}
}

func repl() {
	x := lox.NewExecutor(os.Stdout)
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
		err = exec(x, rd)
		if err != nil {
			lox.ExitErr(err)
		}
	}
}

func exec(x *lox.Executor, r io.Reader) error {
	lexer, err := lox.NewLexer(bufio.NewReader(r))
	if err != nil {
		return err
	}
	tokens, err := lexer.Scan()
	if err != nil {
		return err
	}
	parser := lox.NewParser(tokens)
	stmts, err := parser.Parse()
	if err != nil {
		return err
	}
	for _, stmt := range stmts {
		err := x.Execute(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}
