package lox

import (
	"github.com/rs/zerolog/log"
)

type Typecheckable interface {
	Typecheck(*Context) error
}

func Typecheck(ctx *Context, elems []Statement) error {
	restore := ctx.StartPhase(PhaseTypecheck)
	defer restore()
	for _, elem := range elems {
		log.Debug().Msgf("(%s) checking %s", ctx.Phase(), elem)
		if err := elem.Typecheck(ctx); err != nil {
			log.Error().Msgf("(%s) error in %q: %s", ctx.Phase(), elem, err)
			return err
		}
	}
	return nil
}

func (s *BlockStatement) Typecheck(ctx *Context) error {
	exit := debugEnterEnv(ctx, "<block>")
	defer exit()
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
	return debugSetType(ctx.Phase(), ctx.env, s.name, s.expr.Type())
}

func (s *FunctionDefinitionStatement) Typecheck(ctx *Context) error {
	// TODO: Typecheck function statement
	// ctx.PushEnvironment()
	// log.Debug().Msgf("(typecheck) enter %s", ctx.types.String())
	// defer func() {
	// 	ctx.PopEnvironment()
	// 	log.Debug().Msgf("(typecheck) enter %s", ctx.types.String())
	// }()
	// for _, p := range s.params {
	// 	if _, err := ctx.types.Set(p, TypeAny); err != nil {
	// 		return err
	// 	}
	// 	log.Debug().Msgf("(typecheck) (%d) %s := %s", ctx.types.depth, p, TypeAny)
	// }
	// for _, st := range s.body {
	// 	if err := st.Typecheck(ctx); err != nil {
	// 		return err
	// 	}
	// }
	s.rtype = TypeAny
	return nil
}

func (s *ReturnStatement) Typecheck(ctx *Context) error {
	if err := s.expr.Typecheck(ctx); err != nil {
		return err
	}
	s.typ = s.expr.Type()
	return nil
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
	prev, env := ctx.env.ResolveType(e.name)
	if prev == TypeNone {
		return NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	return debugSetType(ctx.Phase(), env, e.name, e.right.Type())
}

func (e *VariableExpression) Typecheck(ctx *Context) error {
	typ, _ := ctx.env.ResolveType(e.name)
	if typ == TypeNone {
		return NewTypeError(NewUndefinedVariableError(e.name), e.Position())
	}
	e.typ = typ
	return nil
}

func (e *CallExpression) Typecheck(ctx *Context) error {
	// TODO: Typecheck call expression
	// if err := e.callee.Typecheck(ctx); err != nil {
	// 	return err
	// }
	// typ := e.callee.Type()
	// if !typ.Contains(TypeCallable) {
	// 	return NewTypeError(NewTypeNotCallableError(typ), e.Position())
	// }
	// fn, _ := ctx.funcs.Lookup(typ.name)
	// if fn == nil {
	// 	return NewTypeError(NewUndefinedVariableError(typ.name), e.Position())
	// }
	// if arity, nargs := fn.Arity(), len(e.args); nargs != arity {
	// 	return NewTypeError(NewArityMismatchError(arity, nargs), e.Position())
	// }
	// for i, arg := range e.args {
	// 	if err := arg.Typecheck(ctx); err != nil {
	// 		return err
	// 	}
	// 	if param := fn.Params[i]; !param.Type.Contains(arg.Type()) {
	// 		return NewTypeError(NewInvalidArgumentTypeForParameter(arg.Type(), param), e.Position())
	// 	}
	// }
	e.typ = TypeAny
	return nil
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
