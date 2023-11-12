package tool

import (
	"flag"
	"fmt"
	"os"

	"github.com/vdinovi/glox/lox"
)

func Setup() error {
	logLevel := flag.String("log", "", "enable logging at specified level")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [file]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	lox.DisableLogger()
	if *logLevel != "" {
		lox.SetConsoleLogOutput(os.Stderr)
		lox.SetLogLevel(*logLevel)
	}
	return nil
}
