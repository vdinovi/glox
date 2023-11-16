package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) TypeCheckProgram(stmts []Statement) error {
	for _, stmt := range stmts {
		log.Debug().Msgf("(typechecker) checking statement %s", stmt)
		err := stmt.TypeCheck(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckBlockStatement(s *BlockStatement) error {
	ctx.PushEnvironment()
	log.Debug().Msgf("(typechecker) enter %s", ctx.types.String())
	defer func() {
		ctx.PopEnvironment()
		log.Debug().Msgf("(typechecker) enter %s", ctx.types.String())
	}()
	for _, stmt := range s.stmts {
		if err := stmt.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckConditionalStatement(s *ConditionalStatement) error {
	if _, err := s.expr.TypeCheck(ctx); err != nil {
		return err
	}
	if err := s.thenBranch.TypeCheck(ctx); err != nil {
		return err
	}
	if s.elseBranch != nil {
		if err := s.elseBranch.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckWhileStatement(s *WhileStatement) error {
	if _, err := s.expr.TypeCheck(ctx); err != nil {
		return err
	}
	if err := s.body.TypeCheck(ctx); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) TypeCheckForStatement(s *ForStatement) error {
	if s.init != nil {
		if err := s.init.TypeCheck(ctx); err != nil {
			return err
		}
	}
	if s.cond != nil {
		if _, err := s.cond.TypeCheck(ctx); err != nil {
			return err
		}
	}
	if s.incr != nil {
		if _, err := s.incr.TypeCheck(ctx); err != nil {
			return err
		}
	}
	if s.body != nil {
		if err := s.body.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckPrintStatement(s *PrintStatement) error {
	_, err := ctx.TypeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckExpressionStatement(s *ExpressionStatement) error {
	_, err := ctx.TypeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckDeclarationStatement(s *DeclarationStatement) error {
	typ, err := ctx.TypeCheckExpression(s.expr)
	if err != nil {
		return err
	}
	prev, err := ctx.types.Set(s.name, typ)
	if err != nil {
		return err
	}
	if prev == nil {
		log.Debug().Msgf("(typechecker) (%d) %s := %s", ctx.types.depth, s.name, typ)
	} else {
		log.Debug().Msgf("(typechecker) (%d) %s = %s (was %s)", ctx.types.depth, s.name, typ, prev)
	}
	return err
}

func (ctx *Context) TypeCheckExpression(e Expression) (Type, error) {
	return e.TypeCheck(ctx)
}

func (ctx *Context) TypeCheckUnaryExpression(e *UnaryExpression) (right, result Type, err error) {
	right, err = e.right.TypeCheck(ctx)
	if err == nil {
		result, err = resolveUnary(e, right)
		return right, result, err
	}
	return right, result, err
}

func (ctx *Context) TypeCheckBinaryExpression(e *BinaryExpression) (left, right, result Type, err error) {
	left, err = e.left.TypeCheck(ctx)
	if err == nil {
		right, err = e.right.TypeCheck(ctx)
		if err == nil {
			result, err = resolveBinary(e, left, right)
		}
	}
	return left, right, result, err
}

func (ctx *Context) TypeCheckGroupingExpression(e *GroupingExpression) (inner, result Type, err error) {
	inner, err = e.expr.TypeCheck(ctx)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	return inner, inner, nil
}

func (ctx *Context) TypeCheckAssignmentExpression(e *AssignmentExpression) (right, result Type, err error) {
	right, err = e.right.TypeCheck(ctx)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	prev, env := ctx.types.Lookup(e.name)
	if prev == nil {
		return TypeAny, TypeAny, NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	_, err = env.Set(e.name, right)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	log.Debug().Msgf("(typechecker) (%d) %s = %s (was %s)", ctx.types.depth, e.name, right, *prev)
	return right, right, nil
}

func (ctx *Context) TypeCheckVariableExpression(e *VariableExpression) (result Type, err error) {
	typ, _ := ctx.types.Lookup(e.name)
	if typ == nil {
		return TypeAny, NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	return *typ, nil
}

func resolveUnary(e *UnaryExpression, right Type) (result Type, err error) {
	var invalid bool
	switch e.op.Type {
	case OpAdd, OpSubtract:
		switch right {
		case TypeNumeric:
			result = TypeNumeric
		default:
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

func resolveBinary(e *BinaryExpression, left Type, right Type) (result Type, err error) {
	var invalid bool
	switch e.op.Type {
	case OpAnd, OpOr:
		if left == right {
			result = left
		} else {
			result = TypeAny
		}
	case OpAdd:
		if invalid = left != right; invalid {
			break
		}
		switch left {
		case TypeNumeric:
			result = TypeNumeric
		case TypeString:
			result = TypeString
		default:
			invalid = true
		}
	case OpSubtract, OpMultiply, OpDivide:
		if invalid = left != right; invalid {
			break
		}
		switch left {
		case TypeNumeric:
			result = TypeNumeric
		default:
			invalid = true
		}
	case OpEqualTo, OpNotEqualTo:
		if invalid = left != right; invalid {
			break
		}
		result = TypeBoolean
	case OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		if invalid = left != right; invalid {
			break
		}
		switch left {
		case TypeNumeric:
			result = TypeBoolean
		default:
			invalid = true
		}
	default:
		invalid = true
	}
	if invalid {
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	return result, err
}
