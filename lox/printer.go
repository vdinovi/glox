package lox

import (
	"fmt"
	"strconv"
	"strings"
)

var defaultPrinter = DefaultPrinter{}

type Printable interface {
	Print(Printer) (string, error)
}

type UnprintableError struct {
	p Printable
}

func (e UnprintableError) Error() string {
	return fmt.Sprintf("unprintable: %v", e.p)
}

type Printer interface {
	Print(Printable) (string, error)
}

type DefaultPrinter struct{}

func (p *DefaultPrinter) Print(thing Printable) (str string, err error) {
	switch x := thing.(type) {
	case Statement:
		str, err = p.printStatement(x)
	case Expression:
		str, err = p.printExpression(x)
	case Value:
		str, err = p.printValue(x)
	default:
		err = UnprintableError{x}
	}
	if err != nil {
		return "", err
	}
	return str, nil
}

func (p *DefaultPrinter) printStatement(stmt Statement) (str string, err error) {
	switch s := stmt.(type) {
	case *BlockStatement:
		str, err = p.printBlockStatement(s)
	case *ConditionalStatement:
		str, err = p.printConditionalStatement(s)
	case *WhileStatement:
		str = fmt.Sprintf("while ( %s ) %s", s.expr, s.body)
	case *ForStatement:
		str, err = p.printForStatement(s)
	case *ExpressionStatement:
		str = fmt.Sprintf("%s ;", s.expr)
	case *PrintStatement:
		str = fmt.Sprintf("print %s ;", s.expr)
	case *DeclarationStatement:
		str = fmt.Sprintf("var %s = %s ;", s.name, s.expr)
	default:
		err = UnprintableError{s}
	}
	if err != nil {
		return "", err
	}
	return str, nil
}

func (p *DefaultPrinter) printExpression(expr Expression) (str string, err error) {
	switch e := expr.(type) {
	case *UnaryExpression:
		str = fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
	case *BinaryExpression:
		str = fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
	case *GroupingExpression:
		str = fmt.Sprintf("(group %s)", e.expr)
	case *AssignmentExpression:
		str = fmt.Sprintf("(%s = %s)", e.name, e.right)
	case *VariableExpression:
		str = fmt.Sprintf("Var(%s)", e.name)
	case *CallExpression:
		str, err = p.printCallExpression(e)
	case *StringExpression:
		str = fmt.Sprintf("\"%s\"", string(e.value))
	case *NumericExpression:
		str = fmt.Sprint(e.value)
	case *BooleanExpression:
		if e.value {
			str = "true"
		} else {
			str = "false"
		}
	case *NilExpression:
		str = "nil"
	default:
		err = UnprintableError{e}
	}
	if err != nil {
		return "", err
	}
	return str, nil

}

func (p *DefaultPrinter) printValue(val Value) (str string, err error) {
	switch v := val.(type) {
	case ValueString:
		str = fmt.Sprintf("\"%s\"", string(v))
	case ValueNumeric:
		str = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case ValueBoolean:
		if bool(v) {
			str = "true"
		} else {
			str = "false"
		}
	case ValueNil:
		str = "nil"
	default:
		err = UnprintableError{v}
	}
	if err != nil {
		return "", err
	}
	return str, nil
}

func (p *DefaultPrinter) printBlockStatement(s *BlockStatement) (str string, err error) {
	var sb strings.Builder
	sb.WriteString("{ ")
	for _, stmt := range s.stmts {
		fmt.Fprint(&sb, " ", stmt.String())
	}
	sb.WriteString(" } ;")
	return sb.String(), err
}

func (p *DefaultPrinter) printConditionalStatement(s *ConditionalStatement) (str string, err error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "if %s %s", s.expr, s.thenBranch)
	if s.elseBranch != nil {
		fmt.Fprintf(&sb, "else %s", s.elseBranch)
	}
	return sb.String(), err
}

func (p *DefaultPrinter) printForStatement(s *ForStatement) (str string, err error) {
	var sb strings.Builder
	fmt.Fprint(&sb, "for ( ")
	if s.init != nil {
		sb.WriteString(s.init.String())
	} else {
		sb.WriteString(" ;")
	}
	if s.cond != nil {
		sb.WriteString(" ")
		sb.WriteString(s.cond.String())
	}
	sb.WriteString(" ; ")
	if s.incr != nil {
		sb.WriteString(s.incr.String())
	}
	sb.WriteString(" ) ")
	sb.WriteString(s.body.String())
	return sb.String(), err
}

func (p *DefaultPrinter) printCallExpression(e *CallExpression) (str string, err error) {
	var args []string
	for _, arg := range e.args {
		args = append(args, arg.String())
	}
	str = fmt.Sprintf("%s(%s)", e.callee, strings.Join(args, ", "))
	return str, err
}
