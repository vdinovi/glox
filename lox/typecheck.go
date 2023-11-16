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

func (ctx *Context) TypeCheckExpression(e Expression) (typ Type, err error) {
	typ, err = e.TypeCheck(ctx)
	if err != nil {
		return TypeAny, err
	}
	return typ, nil
}

func (ctx *Context) TypeCheckUnaryExpression(e *UnaryExpression) (right, result Type, err error) {
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
	if left, err = e.left.TypeCheck(ctx); err != nil {
		return TypeAny, TypeAny, TypeAny, err
	}
	if right, err = e.right.TypeCheck(ctx); err != nil {
		return TypeAny, TypeAny, TypeAny, err
	}
	switch e.op.Type {
	case OpAnd, OpOr:
		if left == right {
			return left, right, left, nil
		}
		return left, right, TypeAny, nil
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
