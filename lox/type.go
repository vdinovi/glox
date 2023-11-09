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

func TypeCheck(expr Expression) (Type, error) {
	log.Debug().Msgf("(typecheck) type checking expression %s", expr)
	t, err := expr.Type()
	if err != nil {
		log.Error().Msgf("(typecheck) error: %s", err)
		return -1, err
	}
	log.Debug().Msgf("(typecheck) yielded type %s", t)
	return t, err
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
