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
	ctx.types = NewEnvironment(ctx.types)
	defer func() { ctx.types = ctx.types.parent }()
	for _, stmt := range s.stmts {
		if err := stmt.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckPrintStatement(s *PrintStatement) error {
	_, err := s.expr.TypeCheck(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	return nil
}

func (ctx *Context) TypeCheckExpressionStatement(s *ExpressionStatement) error {
	_, err := s.expr.TypeCheck(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	return nil
}

func (ctx *Context) TypeCheckDeclarationStatement(s *DeclarationStatement) error {
	typ, err := s.expr.TypeCheck(ctx)
	if err == nil {
		err = ctx.types.Set(s.name, typ)
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	return nil
}

func (ctx *Context) TypeCheckUnaryExpression(e *UnaryExpression) (right, result Type, err error) {
	right, err = e.right.TypeCheck(ctx)
	if err != nil {
		return ErrType, ErrType, err
	}
	switch right {
	case TypeNumeric:
		result, err = ctx.typeCheckUnaryNumeric(e.op, right)
	default:
		err = NewTypeError(NewInvalidOperatorForTypeError(e.op.Type, right), e.Position())
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return right, result, nil
}

func (ctx *Context) TypeCheckBinaryExpression(e *BinaryExpression) (left, right, result Type, err error) {
	if left, err = e.left.TypeCheck(ctx); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if right, err = e.right.TypeCheck(ctx); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if left != right {
		return ErrType, ErrType, ErrType, NewTypeError(NewTypeMismatchError(left, right), e.Position())
	}
	switch left {
	case TypeNumeric:
		result, err = ctx.typeCheckBinaryNumeric(e.op, left, right)
	case TypeString:
		result, err = ctx.typeCheckBinaryString(e.op, left, right)
	case TypeBoolean:
		result, err = ctx.typeCheckBinaryBoolean(e.op, left, right)
	default:
		err = NewTypeError(NewInvalidOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return left, right, result, nil
}

func (ctx *Context) TypeCheckGroupingExpression(e *GroupingExpression) (inner, result Type, err error) {
	inner, err = e.expr.TypeCheck(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return inner, result, nil
}

func (ctx *Context) TypeCheckAssignmentExpression(e *AssignmentExpression) (right, result Type, err error) {
	right, err = e.right.TypeCheck(ctx)
	if err == nil {
		err = ctx.types.Set(e.name, right)
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	result = right
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return right, result, nil
}

func (ctx *Context) TypeCheckVariableExpression(e *VariableExpression) (result Type, err error) {
	typ := ctx.types.Lookup(e.name)
	if typ == nil {
		err := NewTypeError(NewUndefinedVariableError(e.name), e.Position())
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, err
	}
	result = *typ
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return result, nil
}

func (ctx *Context) typeCheckUnaryNumeric(op Operator, typ Type) (Type, error) {
	switch op.Type {
	case OpSubtract, OpAdd:
		return TypeNumeric, nil
	default:
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, typ), ErrPosition)
	}
}

func (ctx *Context) typeCheckBinaryNumeric(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpAdd, OpSubtract, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualTo, OpNotEqualTo, OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		return TypeBoolean, nil
	default:
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), ErrPosition)
	}
}

func (ctx *Context) typeCheckBinaryString(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpAdd:
		return TypeString, nil
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), ErrPosition)
	}
}

func (ctx *Context) typeCheckBinaryBoolean(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), ErrPosition)
	}
}
