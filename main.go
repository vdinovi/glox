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
		line = strings.TrimRight(line, "\n")
		if err != nil {
			lox.ExitErr(err)
		}
		if line == "" {
			lox.Exit(lox.ExitCodeOK)
		}
		rd := strings.NewReader(line)
		err = exec(rd)
		if err != nil {
			lox.ExitErr(err)
		}
	}
}

func exec(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
