package lox

import (
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Sets the global log level
func SetLogLevel(level string) error {
	l, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(l)
	return nil
}

// Sets the global logger to output to console
// Warning: slow
func SetConsoleLogOutput(output io.Writer) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: output})
}
