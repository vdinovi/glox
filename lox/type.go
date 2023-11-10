package lox

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

//go:generate stringer -type Type  -trimprefix=Type
type Type int

const (
	TypeNil Type = iota
	TypeNumeric
	TypeString
	TypeBoolean
)

type Typed interface {
	Type() (Type, error)
}

func TypeCheck(stmts []Statement) error {
	for _, stmt := range stmts {
		log.Debug().Msgf("(typecheck) type checking statement: %s", stmt)
		err := stmt.TypeCheck()
		if err != nil {
			log.Error().Msgf("(typecheck) error: %s", err)
			return err
		}
	}
	return nil
}

type TypeError struct {
	Expression
	Err error
}

func (e TypeError) Error() string {
	return fmt.Sprintf("type error in %q: %s", e.Expression, e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}
