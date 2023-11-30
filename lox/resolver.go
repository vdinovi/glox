package lox

import (
	"github.com/rs/zerolog/log"
)

func Resolve(ctx *Context, stmts []Statement) error {
	restore := ctx.StartPhase(PhaseResolve)
	defer restore()
	for _, stmt := range stmts {
		log.Debug().Msgf("(%s) resolving %s", ctx.Phase(), stmt)
		if err := stmt.Resolve(ctx); err != nil {
			log.Error().Msgf("(%s) error in %q: %s", ctx.Phase(), stmt, err)
			return err
		}
	}
	return nil
}

func (s *BlockStatement) Resolve(ctx *Context) error {
	// TODO: block introduces new scope
	return nil
}

func (s *FunctionDefinitionStatement) Resolve(ctx *Context) error {
	// TODO: fundef introduces new scope for body and binds params within scope
	return nil
}

func (s *DeclarationStatement) Resolve(ctx *Context) error {
	// TODO: adds var to scope
	return nil
}

func (e *AssignmentExpression) Resolve(ctx *Context) error {
	// TODO: resolve var
	return nil
}

func (e *VariableExpression) Resolve(ctx *Context) error {
	// TODO: resolve var
	return nil
}

func (s *ConditionalStatement) Resolve(ctx *Context) error {
	return nil
}

func (s *WhileStatement) Resolve(ctx *Context) error {
	return nil
}

func (s *ForStatement) Resolve(ctx *Context) error {
	return nil
}

func (s *ExpressionStatement) Resolve(ctx *Context) error {
	return nil
}

func (s *PrintStatement) Resolve(ctx *Context) error {
	return nil
}

func (s *ReturnStatement) Resolve(ctx *Context) error {
	return nil
}

func (e *UnaryExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *BinaryExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *GroupingExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *CallExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *StringExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *NumericExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *BooleanExpression) Resolve(ctx *Context) error {
	return nil
}

func (e *NilExpression) Resolve(ctx *Context) error {
	return nil
}
