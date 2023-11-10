package lox

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func TypeCheck(stmts []Statement) error {
	for _, s := range stmts {
		log.Debug().Msgf("(typechecker) checking statement %s", s)
		if err := s.TypeCheck(); err != nil {
			return err
		}
	}
	return nil
}

func TypeCheckPrintStatement(s PrintStatement) error {
	typ, err := s.expr.Type()
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func TypeCheckExpressionStatement(s ExpressionStatement) error {
	typ, err := s.expr.Type()
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", s, err)
		return err
	}
	log.Debug().Msgf("(typechecker) %q => %s", s, typ)
	return nil
}

func TypeCheckUnaryExpression(e UnaryExpression) (Type, error) {
	var err error
	var typ Type
	if typ, err = e.right.Type(); err != nil {
		return ErrType, err
	}
	switch typ {
	case TypeNumeric:
		typ, err = typeCheckUnaryNumeric(e.op.Type, typ)
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

func TypeCheckBinaryExpression(e BinaryExpression) (Type, error) {
	var err error
	var left, right Type
	if left, err = e.left.Type(); err != nil {
		return ErrType, err
	}
	if right, err = e.right.Type(); err != nil {
		return ErrType, err
	}
	if left != right {
		return ErrType, TypeError{TypeMismatch(left, right)}
	}
	var typ Type
	switch left {
	case TypeNumeric:
		typ, err = typeCheckBinaryNumeric(e.op.Type, left, right)
	case TypeString:
		typ, err = typeCheckBinaryString(e.op.Type, left, right)
	case TypeBoolean:
		typ, err = typeCheckBinaryBoolean(e.op.Type, left, right)
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

func TypeCheckGroupingExpression(e GroupingExpression) (Type, error) {
	typ, err := e.expr.Type()
	if err != nil {
		log.Error().Msgf("(typechecker) error in %q: %s", e, err)
		return ErrType, err
	}
	log.Debug().Msgf("(typechecker) %q => %s", e, typ)
	return typ, err
}

func typeCheckUnaryNumeric(op OperatorType, typ Type) (Type, error) {
	switch op {
	case OpMinus:
		return TypeNumeric, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to type %s", op, typ)}
	}
}

func typeCheckBinaryNumeric(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpPlus, OpMinus, OpMultiply, OpDivide:
		return TypeNumeric, nil
	case OpEqualEquals, OpNotEquals, OpLess, OpLessEquals, OpGreater, OpGreaterEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}

func typeCheckBinaryString(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpPlus:
		return TypeString, nil
	case OpEqualEquals, OpNotEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}

func typeCheckBinaryBoolean(op OperatorType, left, right Type) (Type, error) {
	switch op {
	case OpEqualEquals, OpNotEquals:
		return TypeBoolean, nil
	default:
		return ErrType, TypeError{fmt.Errorf("%s can't be applied to types %s and %s", op, left, right)}
	}
}
