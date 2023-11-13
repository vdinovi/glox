package lox

import (
	"github.com/rs/zerolog/log"
)

func (ctx *EvaluationContext) TypeCheckProgram(prog Program) error {
	for _, stmt := range prog {
		log.Debug().Msgf("(typechecker) checking statement %s", stmt)
		if err := stmt.TypeCheck(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *EvaluationContext) TypeCheckPrintStatement(s PrintStatement) error {
	typ, err := s.expr.Type(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func (ctx *EvaluationContext) TypeCheckExpressionStatement(s ExpressionStatement) error {
	typ, err := s.expr.Type(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func (ctx *EvaluationContext) TypeCheckUnaryExpression(e UnaryExpression) (right, result Type, err error) {
	right, err = e.right.Type(ctx)
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

func (ctx *EvaluationContext) TypeCheckBinaryExpression(e BinaryExpression) (left, right, result Type, err error) {
	if left, err = e.left.Type(ctx); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if right, err = e.right.Type(ctx); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if left != right {
		return ErrType, ErrType, ErrType, NewTypeError(NewTypeMismatchError(left, right), e.Position())
	}
	var typ Type
	switch left {
	case TypeNumeric:
		typ, err = ctx.typeCheckBinaryNumeric(e.op, left, right)
	case TypeString:
		typ, err = ctx.typeCheckBinaryString(e.op, left, right)
	case TypeBoolean:
		typ, err = ctx.typeCheckBinaryBoolean(e.op, left, right)
	default:
		err = NewTypeError(NewInvalidOperatorForTypeError(e.op.Type, left, right), e.Position())
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return left, right, typ, err
}

func (ctx *EvaluationContext) TypeCheckGroupingExpression(e GroupingExpression) (inner, result Type, err error) {
	inner, err = e.expr.Type(ctx)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return inner, result, err
}

func (ctx *EvaluationContext) TypeCheckStringExpression(e StringExpression) (Type, error) {
	return TypeString, nil
}

func (ctx *EvaluationContext) TypeCheckNumericExpression(NumericExpression) (Type, error) {
	return TypeNumeric, nil
}

func (ctx *EvaluationContext) TypeCheckBooleanExpression(BooleanExpression) (Type, error) {
	return TypeBoolean, nil
}

func (ctx *EvaluationContext) TypeCheckNilExpression(e NilExpression) (Type, error) {
	return TypeNil, nil
}

func (ctx *EvaluationContext) typeCheckUnaryNumeric(op Operator, typ Type) (Type, error) {
	switch op.Type {
	case OpAdd:
		return TypeNumeric, nil
	default:
		// TODO
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, typ), Position{})
	}
}

func (ctx *EvaluationContext) typeCheckBinaryNumeric(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpAdd, OpSubtract, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualTo, OpNotEqualTo, OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		return TypeBoolean, nil
	default:
		// TODO
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), Position{})
	}
}

func (ctx *EvaluationContext) typeCheckBinaryString(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpAdd:
		return TypeString, nil
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		// TODO
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), Position{})
	}
}

func (ctx *EvaluationContext) typeCheckBinaryBoolean(op Operator, left, right Type) (Type, error) {
	switch op.Type {
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		// TODO
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(op.Type, left, right), Position{})
	}
}
