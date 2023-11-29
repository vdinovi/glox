package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
		lox.ExitErr(err)
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
	ctx := lox.NewContext(w)
	tokens, err := lox.Scan(ctx, bufio.NewReader(r))
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
