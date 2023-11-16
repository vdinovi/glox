package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *Context) TypeCheckProgram(stmts []Statement) error {
	log.Trace().Msg("TypeCheckProgram")
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
	log.Trace().Msg("TypeCheckBlockStatement")
	ctx.PushEnvironment()
	log.Debug().Msgf("(typechecker) entering scope {%s}", ctx.types.String())
	defer func() {
		ctx.PopEnvironment()
		log.Debug().Msgf("(typechecker) entering scope {%s}", ctx.types.String())
	}()
	for _, stmt := range s.stmts {
		if err := stmt.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckConditionalStatement(s *ConditionalStatement) error {
	log.Trace().Msg("TypeCheckConditionalStatement")
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

func (ctx *Context) TypeCheckPrintStatement(s *PrintStatement) error {
	log.Trace().Msg("TypeCheckPrintStatement")
	_, err := ctx.TypeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckExpressionStatement(s *ExpressionStatement) error {
	log.Trace().Msg("TypeCheckExpressionStatement")
	_, err := ctx.TypeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckDeclarationStatement(s *DeclarationStatement) error {
	log.Trace().Msg("TypeCheckDeclarationStatement")
	typ, err := ctx.TypeCheckExpression(s.expr)
	if err != nil {
		return err
	}
	prev := ctx.types.Get(s.name, TypeAny)
	err = ctx.types.Set(s.name, typ)
	if err != nil {
		return err
	}
	if prev == TypeAny {
		log.Debug().Msgf("(typechecker) initialized type(%s) to %s", s.name, typ)
	} else {
		log.Debug().Msgf("(typechecker) %s <- %s (prev %s)", s.name, typ, prev)
	}
	return err
}

func (ctx *Context) TypeCheckExpression(e Expression) (typ Type, err error) {
	typ, err = e.TypeCheck(ctx)
	if err != nil {
		return TypeAny, err
	}
	log.Debug().Msgf("(typechecker) type(%s) => %s", e, typ)
	return typ, nil
}

func (ctx *Context) TypeCheckUnaryExpression(e *UnaryExpression) (right, result Type, err error) {
	log.Trace().Msg("TypeCheckUnaryExpression")
	right, err = e.right.TypeCheck(ctx)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	switch right {
	case TypeNumeric:
		result, err = ctx.typeCheckUnaryNumeric(e, right)
	default:
		err = NewTypeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right), e.Position())
	}
	if err != nil {
		return TypeAny, TypeAny, err
	}
	return right, result, err
}

func (ctx *Context) TypeCheckBinaryExpression(e *BinaryExpression) (left, right, result Type, err error) {
	log.Trace().Msg("TypeCheckBinaryExpression")
	if left, err = e.left.TypeCheck(ctx); err != nil {
		return TypeAny, TypeAny, TypeAny, err
	}
	if right, err = e.right.TypeCheck(ctx); err != nil {
		return TypeAny, TypeAny, TypeAny, err
	}
	if left != right {
		return TypeAny, TypeAny, TypeAny, NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	switch left {
	case TypeNumeric:
		result, err = ctx.typeCheckBinaryNumeric(e, left, right)
	case TypeString:
		result, err = ctx.typeCheckBinaryString(e, left, right)
	case TypeBoolean:
		result, err = ctx.typeCheckBinaryBoolean(e, left, right)
	default:
		err = NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	if err != nil {
		return TypeAny, TypeAny, TypeAny, err
	}
	return left, right, result, nil
}

func (ctx *Context) TypeCheckGroupingExpression(e *GroupingExpression) (inner, result Type, err error) {
	log.Trace().Msg("TypeCheckGroupingExpression")
	inner, err = e.expr.TypeCheck(ctx)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	return inner, inner, nil
}

func (ctx *Context) TypeCheckAssignmentExpression(e *AssignmentExpression) (right, result Type, err error) {
	log.Trace().Msg("TypeCheckAssignmentExpression")
	right, err = e.right.TypeCheck(ctx)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	err = ctx.types.Set(e.name, right)
	if err != nil {
		return TypeAny, TypeAny, err
	}
	return right, right, nil
}

func (ctx *Context) TypeCheckVariableExpression(e *VariableExpression) (result Type, err error) {
	log.Trace().Msg("TypeCheckVariableExpression")
	typ, _ := ctx.types.Lookup(e.name)
	if typ == nil {
		return TypeAny, NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	return *typ, nil
}

func (ctx *Context) typeCheckUnaryNumeric(e *UnaryExpression, right Type) (Type, error) {
	switch e.op.Type {
	case OpSubtract, OpAdd:
		return TypeNumeric, nil
	default:
		return TypeAny, NewTypeError(NewInvalidUnaryOperatorForTypeError(e.op.Type, right), e.Position())
	}
}

func (ctx *Context) typeCheckBinaryNumeric(e *BinaryExpression, left, right Type) (Type, error) {
	switch e.op.Type {
	case OpAdd, OpSubtract, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualTo, OpNotEqualTo, OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		return TypeBoolean, nil
	default:
		return TypeAny, NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
}

func (ctx *Context) typeCheckBinaryString(e *BinaryExpression, left, right Type) (Type, error) {
	switch e.op.Type {
	case OpAdd:
		return TypeString, nil
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		return TypeAny, NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
}

func (ctx *Context) typeCheckBinaryBoolean(e *BinaryExpression, left, right Type) (Type, error) {
	switch e.op.Type {
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		return TypeAny, NewTypeError(NewInvalidBinaryOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
}
