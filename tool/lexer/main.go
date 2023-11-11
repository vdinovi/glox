package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vdinovi/glox/lox"
)

func main() {
	logLevel := flag.String("log", "", "enable logging at specified level")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [file]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		lox.Exit(lox.ExitCodeErr)
	}

	lox.DisableLogger()
	if *logLevel != "" {
		lox.SetConsoleLogOutput(os.Stderr)
		lox.SetLogLevel(*logLevel)
	}

	path := filepath.Clean(os.Args[1])
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
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	written := 0
	for written < len(data) {
		n, err := w.Write(data[written:])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}
