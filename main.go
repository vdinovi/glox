package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vdinovi/glox/lox"
)

const usagef = `Usage: %s [file]
       starts a repl if no file is provided.
`

func main() {
	err := setup()
	if err == nil {
		if flag.NArg() == 0 {
			err = interactive()
		} else {
			err = file(flag.Arg(0))
		}
	}
	if err != nil {
		lox.ExitErr(err)
	}
	lox.Exit(lox.ExitCodeOK)
}

func setup() error {
	logLevel := flag.String("log", "", "enable logging at specified level")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), usagef, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *logLevel == "" {
		lox.DisableLogger()
	} else {
		lox.SetConsoleLogOutput(os.Stderr)
		lox.SetLogLevel(*logLevel)
	}
	return nil
}

func file(fpath string) error {
	log.Debug().Msgf("(main) executing %s", fpath)
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	executor := lox.NewExecutor(os.Stdout)
	reader := bufio.NewReader(f)

	err = execute(executor, reader)
	if err != nil {
		return lox.FatalError{err}
	}
	return nil
}

func interactive() (err error) {
	terminal, err := lox.NewTerminal(os.Stdin, "(lox) ")
	if err != nil {
		return err
	}
	defer terminal.Close()
	defer func() {
		if err := recover(); err != nil {
			terminal.Close()
			panic(err)
		}
	}()

	executor := lox.NewExecutor(terminal)

	var line string
	for {
		line, err = terminal.ReadLine()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		err = execute(executor, strings.NewReader(line))
		if err == nil {
			continue
		} else if errors.Is(err, lox.FatalError{}) {
			return err
		} else if e := terminal.WriteError(err); e != nil {
			return e
		}
	}
}

func execute(executor *lox.Executor, reader io.Reader) error {
	tokens, err := lox.Scan(bufio.NewReader(reader))
	if err != nil {
		return err
	}
	stmts, err := lox.Parse(tokens)
	if err != nil {
		return err
	}
	if err = executor.TypeCheckProgram(stmts); err != nil {
		return err
	}
	return executor.ExecuteProgram(stmts)
}
