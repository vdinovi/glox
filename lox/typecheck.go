package lox

import (
	"github.com/rs/zerolog/log"
)

type Typecheckable interface {
	Typecheck(*Context) error
}

func (ctx *Context) Typecheck(elems []Statement) error {
	for _, elem := range elems {
		log.Debug().Msgf("(typecheck) checking %s", elem)
		if err := elem.Typecheck(ctx); err != nil {
			log.Error().Msgf("(typecheck) error in %q: %s", elem, err)
			return err
		}
	}
	return nil
}

func (s *BlockStatement) Typecheck(ctx *Context) error {
	ctx.PushEnvironment()
	log.Debug().Msgf("(typecheck) enter %s", ctx.types.String())
	defer func() {
		ctx.PopEnvironment()
		log.Debug().Msgf("(typecheck) enter %s", ctx.types.String())
	}()
	for _, stmt := range s.stmts {
		if err := stmt.Typecheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *ConditionalStatement) Typecheck(ctx *Context) error {
	if err := s.expr.Typecheck(ctx); err != nil {
		return err
	}
	if err := s.thenBranch.Typecheck(ctx); err != nil {
		return err
	}
	if s.elseBranch != nil {
		if err := s.elseBranch.Typecheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *WhileStatement) Typecheck(ctx *Context) error {
	if err := s.expr.Typecheck(ctx); err != nil {
		return err
	}
	if err := s.body.Typecheck(ctx); err != nil {
		return err
	}
	return nil
}

func (s *ForStatement) Typecheck(ctx *Context) error {
	if s.init != nil {
		if err := s.init.Typecheck(ctx); err != nil {
			return err
		}
	}
	if s.cond != nil {
		if err := s.cond.Typecheck(ctx); err != nil {
			return err
		}
	}
	if s.incr != nil {
		if err := s.incr.Typecheck(ctx); err != nil {
			return err
		}
	}
	if s.body != nil {
		if err := s.body.Typecheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *PrintStatement) Typecheck(ctx *Context) error {
	return s.expr.Typecheck(ctx)
}

func (s *ExpressionStatement) Typecheck(ctx *Context) error {
	return s.expr.Typecheck(ctx)
}

func (s *DeclarationStatement) Typecheck(ctx *Context) error {
	err := s.expr.Typecheck(ctx)
	if err != nil {
		return err
	}
	typ := s.expr.Type()
	prev, err := ctx.types.Set(s.name, typ)
	if err != nil {
		return err
	}
	if prev == nil {
		log.Debug().Msgf("(typecheck) (%d) %s := %s", ctx.types.depth, s.name, typ)
	} else {
		log.Debug().Msgf("(typecheck) (%d) %s = %s (was %s)", ctx.types.depth, s.name, typ, prev)
	}
	return err
}

func (e *UnaryExpression) Typecheck(ctx *Context) error {
	if err := e.right.Typecheck(ctx); err != nil {
		return err
	}
	typ, err := typecheckUnary(e, e.right.Type())
	if err != nil {
		return err
	}
	e.typ = typ
	return nil
}

func (e *BinaryExpression) Typecheck(ctx *Context) error {
	if err := e.left.Typecheck(ctx); err != nil {
		return err
	}
	if err := e.right.Typecheck(ctx); err != nil {
		return err
	}
	typ, err := typecheckBinary(e, e.left.Type(), e.right.Type())
	if err != nil {
		return err
	}
	e.typ = typ
	return nil
}

func (e *GroupingExpression) Typecheck(ctx *Context) error {
	if err := e.expr.Typecheck(ctx); err != nil {
		return err
	}
	e.typ = e.expr.Type()
	return nil
}

func (e *AssignmentExpression) Typecheck(ctx *Context) error {
	if err := e.right.Typecheck(ctx); err != nil {
		return err
	}
	prev, env := ctx.types.Lookup(e.name)
	if prev == nil {
		return NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	if _, err := env.Set(e.name, e.right.Type()); err != nil {
		return err
	}
	log.Debug().Msgf("(typecheck) (%d) %s = %s (was %s)", ctx.types.depth, e.name, e.right.Type(), *prev)
	return nil
}

func (e *VariableExpression) Typecheck(ctx *Context) error {
	typ, _ := ctx.types.Lookup(e.name)
	if typ == nil {
		return NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	e.typ = *typ
	return nil
}

func (e *CallExpression) Typecheck(ctx *Context) error {
	return ErrNotYetImplemented
}

func (e *StringExpression) Typecheck(*Context) error {
	return nil
}

func (e *NumericExpression) Typecheck(*Context) error {
	return nil
}

func (e *BooleanExpression) Typecheck(*Context) error {
	return nil
}

func (e *NilExpression) Typecheck(*Context) error {
	return nil
}

func typecheckUnary(e *UnaryExpression, right Type) (result Type, err error) {
	var invalid bool
	switch e.op.Type {
	case OpAdd, OpSubtract:
		if right.Within(TypeNumeric) {
			result = right
		} else {
			invalid = true
		}
	case OpNegate:
		result = TypeBoolean
	default:
		invalid = true
	}
	if invalid {
		err = NewTypeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right), e.Position())
	}
	return result, err
}

func typecheckBinary(e *BinaryExpression, left Type, right Type) (result Type, err error) {
	var invalid bool
	switch e.op.Type {
	case OpAnd, OpOr:
		result = left.Union(right)
	case OpAdd:
		if left.Within(TypeNumeric, TypeString) && right.Within(TypeNumeric, TypeString) && left.Contains(right) {
			result = left.Union(right)
		} else {
			invalid = true
		}
	case OpSubtract, OpMultiply, OpDivide:
		if left.Within(TypeNumeric) && right.Within(TypeNumeric) && left.Contains(right) {
			result = left.Union(right)
		} else {
			invalid = true
		}
	case OpEqualTo, OpNotEqualTo:
		if left.Contains(right) {
			result = TypeBoolean
		} else {
			invalid = true
		}
	case OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		if left.Within(TypeNumeric) && right.Within(TypeNumeric) && left.Contains(right) {
			result = TypeBoolean
		} else {
			invalid = true
		}
	default:
		invalid = true
	}
	if invalid {
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	s := result.String()
	_ = s
	return result, err
}
