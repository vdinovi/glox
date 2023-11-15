package main

import (
	"bufio"
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
	return execute(executor, reader)
}

func interactive() (err error) {
	executor := lox.NewExecutor(os.Stdout)
	reader := bufio.NewReader(os.Stdin)
	var line string
	for err == nil {
		fmt.Printf("> ")
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		line = strings.TrimRight(line, "\n")
		err = execute(executor, strings.NewReader(line))
	}
	return err
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
