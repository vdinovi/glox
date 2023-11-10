package lox

import (
	"fmt"

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

func (syms Symbols) TypeCheckUnaryExpression(e UnaryExpression) (Type, error) {
	var err error
	var typ Type
	if typ, err = e.right.Type(syms); err != nil {
		return ErrType, err
	}
	switch typ {
	case TypeNumeric:
		typ, err = syms.typeCheckUnaryNumeric(e.op.Type, typ)
	default:
		err = TypeError{fmt.Errorf("%s can't be applied to type %s", e.op.Type, typ)}
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return typ, nil
}

func (syms Symbols) TypeCheckBinaryExpression(e BinaryExpression) (Type, error) {
	var err error
	var left, right Type
	if left, err = e.left.Type(syms); err != nil {
		return ErrType, err
	}
	if right, err = e.right.Type(syms); err != nil {
		return ErrType, err
	}
	if left != right {
		return ErrType, TypeError{TypeMismatch(left, right)}
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
		err = TypeError{fmt.Errorf("%s can't be applied to types %s and %s", e.op.Type, left, right)}
	}
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return typ, err
}

func (syms Symbols) TypeCheckGroupingExpression(e GroupingExpression) (Type, error) {
	typ, err := e.expr.Type(syms)
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return typ, err
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

func (syms Symbols) typeCheckUnaryNumeric(op OperatorType, typ Type) (Type, error) {
	switch op {
	case OpMinus:
		return TypeNumeric, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to type %s", op, typ)}
	}
}

func (syms Symbols) typeCheckBinaryNumeric(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpPlus, OpMinus, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualEquals, OpNotEquals, OpLess, OpLessEquals, OpGreater, OpGreaterEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}

func (syms Symbols) typeCheckBinaryString(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpPlus:
		return TypeString, nil
	case OpEqualEquals, OpNotEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}

func (syms Symbols) typeCheckBinaryBoolean(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpEqualEquals, OpNotEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}
