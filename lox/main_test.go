package lox

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestMain(m *testing.M) {
	//zerolog.SetGlobalLevel(zerolog.Disabled)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	os.Exit(m.Run())
}
