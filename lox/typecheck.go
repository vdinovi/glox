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
	log.Debug().Msgf("(typechecker) enter scope")
	ctx.PushEnvironment()
	defer ctx.PopEnvironment()
	defer log.Debug().Msgf("(typechecker) exit scope")
	for _, stmt := range s.stmts {
		if err := stmt.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *Context) TypeCheckPrintStatement(s *PrintStatement) error {
	log.Trace().Msg("TypeCheckPrintStatement")
	_, err := ctx.typeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckExpressionStatement(s *ExpressionStatement) error {
	log.Trace().Msg("TypeCheckExpressionStatement")
	_, err := ctx.typeCheckExpression(s.expr)
	return err
}

func (ctx *Context) TypeCheckDeclarationStatement(s *DeclarationStatement) error {
	log.Trace().Msg("TypeCheckDeclarationStatement")
	typ, err := ctx.typeCheckExpression(s.expr)
	if err == nil {
		err = ctx.types.Set(s.name, typ)
	}
	if err == nil {
		log.Debug().Msgf("(typechecker) type(%s) <- %s", s.name, typ)
	}
	return err
}

func (ctx *Context) typeCheckExpression(e Expression) (typ Type, err error) {
	typ, err = e.TypeCheck(ctx)
	if err != nil {
		return ErrType, err
	}
	log.Debug().Msgf("(typechecker) type(%s) => %s", e, typ)
	return typ, nil
}

func (ctx *Context) TypeCheckUnaryExpression(e *UnaryExpression) (right, result Type, err error) {
	log.Trace().Msg("TypeCheckUnaryExpression")
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
		return ErrType, ErrType, err
	}
	return right, result, err
}

func (ctx *Context) TypeCheckBinaryExpression(e *BinaryExpression) (left, right, result Type, err error) {
	log.Trace().Msg("TypeCheckBinaryExpression")
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
		return ErrType, ErrType, ErrType, err
	}
	return left, right, result, nil
}

func (ctx *Context) TypeCheckGroupingExpression(e *GroupingExpression) (inner, result Type, err error) {
	log.Trace().Msg("TypeCheckGroupingExpression")
	inner, err = e.expr.TypeCheck(ctx)
	if err != nil {
		return ErrType, ErrType, err
	}
	return inner, inner, nil
}

func (ctx *Context) TypeCheckAssignmentExpression(e *AssignmentExpression) (right, result Type, err error) {
	log.Trace().Msg("TypeCheckAssignmentExpression")
	right, err = e.right.TypeCheck(ctx)
	if err == nil {
		err = ctx.types.Set(e.name, right)
	}
	if err != nil {
		return ErrType, ErrType, err
	}
	return right, right, nil
}

func (ctx *Context) TypeCheckVariableExpression(e *VariableExpression) (result Type, err error) {
	log.Trace().Msg("TypeCheckVariableExpression")
	typ := ctx.types.Lookup(e.name)
	if typ == nil {
		return ErrType, NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	return *typ, nil
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
