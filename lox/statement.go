package lox

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Statement interface {
	Execute() error
	TypeCheck() error
	fmt.Stringer
}

type ExpressionStatement struct {
	expr Expression
}

func (s ExpressionStatement) TypeCheck() error {
	t, err := s.expr.Type()
	log.Debug().Msgf("(typecheck) %s expression yields type %s", s, t)
	return err
}

func (s ExpressionStatement) Execute() error {
	log.Debug().Msgf("(runtime) executing statement: %s", s)
	_, _, err := s.expr.Evaluate()
	if err != nil {
		log.Debug().Msgf("(runtime) error in %s: %s", s, err)
		return err
	}
	return nil
}

func (s ExpressionStatement) String() string {
	return fmt.Sprintf("%s ;", s.expr)
}

type PrintStatement struct {
	expr Expression
}

// TODO: ensure only strings can be printed
func (s PrintStatement) TypeCheck() error {
	t, err := s.expr.Type()
	log.Debug().Msgf("(typecheck) %s expression yields type %s", s, t)
	return err
}

func (s PrintStatement) Execute() error {
	log.Debug().Msgf("(runtime) executing statement: %s", s)
	val, _, err := s.expr.Evaluate()
	if err != nil {
		log.Debug().Msgf("(runtime) error in %s: %s", s, err)
		return err
	}
	return s.print(val)
}

func (s PrintStatement) String() string {
	return fmt.Sprintf("print %s ;", s.expr)
}

func (s PrintStatement) print(val Value) error {
	_, err := fmt.Println(val)
	return err
}
