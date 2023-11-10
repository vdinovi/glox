package lox

import (
	"github.com/rs/zerolog/log"
)

func (syms Symbols) TypeCheckProgram(prog Program) error {
	for _, stmt := range prog {
		log.Debug().Msgf("(typechecker) checking statement %s", stmt)
		if err := stmt.TypeCheck(syms); err != nil {
			return err
		}
	}
	return nil
}

func (syms Symbols) TypeCheckPrintStatement(s PrintStatement) error {
	typ, err := s.expr.Type(syms)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func (syms Symbols) TypeCheckExpressionStatement(s ExpressionStatement) error {
	typ, err := s.expr.Type(syms)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func (syms Symbols) TypeCheckUnaryExpression(e UnaryExpression) (right, result Type, err error) {
	right, err = e.right.Type(syms)
	if err != nil {
		return ErrType, ErrType, err
	}
	switch right {
	case TypeNumeric:
		result, err = syms.typeCheckUnaryNumeric(e.op.Type, right)
	default:
		// TODO: feed line and column
		err = NewTypeError(NewInvalidOperatorForTypeError(e.op.Type, right), 0, 0)
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return right, result, nil
}

func (syms Symbols) TypeCheckBinaryExpression(e BinaryExpression) (left, right, result Type, err error) {
	if left, err = e.left.Type(syms); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if right, err = e.right.Type(syms); err != nil {
		return ErrType, ErrType, ErrType, err
	}
	if left != right {
		// TODO: feed line and column
		return ErrType, ErrType, ErrType, NewTypeError(NewTypeMismatchError(left, right), 0, 0)
	}
	var typ Type
	switch left {
	case TypeNumeric:
		typ, err = syms.typeCheckBinaryNumeric(e.op.Type, left, right)
	case TypeString:
		typ, err = syms.typeCheckBinaryString(e.op.Type, left, right)
	case TypeBoolean:
		typ, err = syms.typeCheckBinaryBoolean(e.op.Type, left, right)
	default:
		// TODO: feed line and column
		err = NewTypeError(NewInvalidOperatorForTypeError(e.op.Type, left, right), 0, 0)
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return left, right, typ, err
}

func (syms Symbols) TypeCheckGroupingExpression(e GroupingExpression) (inner, result Type, err error) {
	inner, err = e.expr.Type(syms)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, result)
	return inner, result, err
}

func (syms Symbols) TypeCheckStringExpression(e StringExpression) (Type, error) {
	return TypeString, nil
}

func (syms Symbols) TypeCheckNumericExpression(NumericExpression) (Type, error) {
	return TypeNumeric, nil
}

func (syms Symbols) TypeCheckBooleanExpression(BooleanExpression) (Type, error) {
	return TypeBoolean, nil
}

func (syms Symbols) TypeCheckNilExpression(e NilExpression) (Type, error) {
	return TypeNil, nil
}

func (syms Symbols) typeCheckUnaryNumeric(opType OperatorType, typ Type) (Type, error) {
	switch opType {
	case OpAdd:
		return TypeNumeric, nil
	default:
		// TODO: feed line and column
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(opType, typ), 0, 0)
	}
}

func (syms Symbols) typeCheckBinaryNumeric(opType OperatorType, left, right Type) (Type, error) {
	switch opType {
	case OpAdd, OpSubtract, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualTo, OpNotEqualTo, OpLessThan, OpLessThanOrEqualTo, OpGreaterThan, OpGreaterThanOrEqualTo:
		return TypeBoolean, nil
	default:
		// TODO: feed line and column
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(opType, left, right), 0, 0)
	}
}

func (syms Symbols) typeCheckBinaryString(opType OperatorType, left, right Type) (Type, error) {
	switch opType {
	case OpAdd:
		return TypeString, nil
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		// TODO: feed line and column
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(opType, left, right), 0, 0)
	}
}

func (syms Symbols) typeCheckBinaryBoolean(opType OperatorType, left, right Type) (Type, error) {
	switch opType {
	case OpEqualTo, OpNotEqualTo:
		return TypeBoolean, nil
	default:
		// TODO: feed line and column
		return ErrType, NewTypeError(NewInvalidOperatorForTypeError(opType, left, right), 0, 0)
	}
}
