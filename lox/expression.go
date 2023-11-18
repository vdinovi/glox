package lox

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type OperatorType int

//go:generate stringer -type OperatorType -trimprefix=Op
const (
	ErrOp OperatorType = iota
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpAnd
	OpOr
	OpEqualTo
	OpNotEqualTo
	OpLessThan
	OpLessThanOrEqualTo
	OpGreaterThan
	OpGreaterThanOrEqualTo
)

type Operator struct {
	Type  OperatorType // type of the operator
	Lexem string       // associated string
}

type Expression interface {
	fmt.Stringer
	Printable
	Located
	Typecheckable
	Evaluable
	Typed
	Equals(Expression) bool
}

type UnaryExpression struct {
	op    Operator
	right Expression
	pos   Position
	typ   Type
}

func (e *UnaryExpression) Position() Position {
	return e.pos
}

func (e *UnaryExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *UnaryExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *UnaryExpression) Type() Type {
	return e.typ
}

func (e *UnaryExpression) TypeCheck(ctx *Context) error {
	_, typ, err := ctx.TypeCheckUnaryExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *UnaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateUnaryExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *UnaryExpression) Equals(other Expression) bool {
	unary, ok := other.(*UnaryExpression)
	if !ok || e.op != unary.op {
		return false
	}
	return e.right.Equals(unary.right)
}

type BinaryExpression struct {
	op    Operator
	left  Expression
	right Expression
	pos   Position
	typ   Type
}

func (e *BinaryExpression) Position() Position {
	return e.pos
}

func (e *BinaryExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *BinaryExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *BinaryExpression) Type() Type {
	return e.typ
}

func (e *BinaryExpression) TypeCheck(ctx *Context) error {
	_, _, typ, err := ctx.TypeCheckBinaryExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *BinaryExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateBinaryExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *BinaryExpression) Equals(other Expression) bool {
	binary, ok := other.(*BinaryExpression)
	if !ok || e.op != binary.op {
		return false
	}
	return e.left.Equals(binary.left) && e.right.Equals(binary.right)
}

type GroupingExpression struct {
	expr Expression
	pos  Position
	typ  Type
}

func (e *GroupingExpression) Position() Position {
	return e.pos
}

func (e *GroupingExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *GroupingExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *GroupingExpression) Type() Type {
	return e.typ
}

func (e *GroupingExpression) TypeCheck(ctx *Context) error {
	_, typ, err := ctx.TypeCheckGroupingExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *GroupingExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateGroupingExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *GroupingExpression) Equals(other Expression) bool {
	group, ok := other.(*GroupingExpression)
	if !ok {
		return false
	}
	return e.expr.Equals(group.expr)
}

type AssignmentExpression struct {
	name  string
	right Expression
	pos   Position
	typ   Type
}

func (e *AssignmentExpression) Position() Position {
	return e.pos
}

func (e *AssignmentExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *AssignmentExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *AssignmentExpression) Type() Type {
	return e.typ
}

func (e *AssignmentExpression) TypeCheck(ctx *Context) error {
	_, typ, err := ctx.TypeCheckAssignmentExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *AssignmentExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateAssignmentExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *AssignmentExpression) Equals(other Expression) bool {
	assign, ok := other.(*AssignmentExpression)
	if !ok || e.name != assign.name {
		return false
	}
	return e.right.Equals(assign.right)
}

type VariableExpression struct {
	name string
	pos  Position
	typ  Type
}

func (e *VariableExpression) Position() Position {
	return e.pos
}

func (e *VariableExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *VariableExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *VariableExpression) Type() Type {
	return e.typ
}

func (e *VariableExpression) TypeCheck(ctx *Context) error {
	typ, err := ctx.TypeCheckVariableExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *VariableExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateVariableExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *VariableExpression) Equals(other Expression) bool {
	variable, ok := other.(*VariableExpression)
	return ok && e.name == variable.name
}

type CallExpression struct {
	callee Expression
	args   []Expression
	pos    Position
	typ    Type
}

func (e *CallExpression) Position() Position {
	return e.pos
}

func (e *CallExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *CallExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *CallExpression) Type() Type {
	return e.typ
}

func (e *CallExpression) TypeCheck(ctx *Context) error {
	typ, err := ctx.TypeCheckCallExpression(e)
	if err != nil {
		return err
	}
	e.typ = typ
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, typ)
	return nil
}

func (e *CallExpression) Evaluate(ctx *Context) (Value, error) {
	val, _, err := ctx.EvaluateCallExpression(e)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, err
}

func (e *CallExpression) Equals(other Expression) bool {
	call, ok := other.(*CallExpression)
	if !ok {
		return false
	}
	if len(e.args) != len(call.args) {
		return false
	}
	for i, arg := range e.args {
		if arg != call.args[i] {
			return false
		}
	}
	return e.callee == call.callee
}

type StringExpression struct {
	value string
	pos   Position
}

func (e *StringExpression) Position() Position {
	return e.pos
}

func (e *StringExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *StringExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *StringExpression) Type() Type {
	return TypeString
}

func (e *StringExpression) TypeCheck(*Context) error {
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, e.Type())
	return nil
}

func (e *StringExpression) Evaluate(ctx *Context) (Value, error) {
	val := ValueString(e.value)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, nil
}

func (e *StringExpression) Equals(other Expression) bool {
	str, ok := other.(*StringExpression)
	return ok && e.value == str.value
}

type NumericExpression struct {
	value float64
	pos   Position
}

func (e *NumericExpression) Position() Position {
	return e.pos
}

func (e *NumericExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *NumericExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *NumericExpression) Type() Type {
	return TypeNumeric
}

func (e *NumericExpression) TypeCheck(*Context) error {
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, e.Type())
	return nil
}

func (e *NumericExpression) Evaluate(ctx *Context) (Value, error) {
	val := ValueNumeric(e.value)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, nil
}

func (e *NumericExpression) Equals(other Expression) bool {
	num, ok := other.(*NumericExpression)
	return ok && e.value == num.value
}

type BooleanExpression struct {
	value bool
	pos   Position
}

func (e *BooleanExpression) Position() Position {
	return e.pos
}

func (e *BooleanExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *BooleanExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *BooleanExpression) Type() Type {
	return TypeBoolean
}

func (e *BooleanExpression) TypeCheck(*Context) error {
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, e.Type())
	return nil
}

func (e *BooleanExpression) Evaluate(ctx *Context) (Value, error) {
	val := ValueBoolean(e.value)
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, nil
}

func (e *BooleanExpression) Equals(other Expression) bool {
	boolean, ok := other.(*BooleanExpression)
	return ok && e.value == boolean.value
}

type NilExpression struct {
	pos Position
}

func (e *NilExpression) Position() Position {
	return e.pos
}

func (e *NilExpression) String() string {
	str, err := e.Print(&defaultPrinter)
	if err != nil {
		panic(err)
	}
	return str
}

func (e *NilExpression) Print(p Printer) (string, error) {
	return p.Print(e)
}

func (e *NilExpression) Type() Type {
	return TypeNil
}

func (e *NilExpression) TypeCheck(*Context) error {
	log.Debug().Msgf("(typechecker) typecheck(%s) => %s", e, e.Type())
	return nil
}

func (e *NilExpression) Evaluate(ctx *Context) (Value, error) {
	val := Nil
	log.Debug().Msgf("(executor) eval(%s) => %s", e, val)
	return val, nil
}

func (e *NilExpression) Equals(other Expression) bool {
	_, ok := other.(*NilExpression)
	return ok
}
