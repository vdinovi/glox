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

func (p *DefaultPrinter) Print(elem Printable) (string, error) {
	return elem.Print(p)
}

func (s *ConditionalStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		var sb strings.Builder
		fmt.Fprintf(&sb, "if %s %s", s.expr, s.thenBranch)
		if s.elseBranch != nil {
			fmt.Fprintf(&sb, "else %s", s.elseBranch)
		}
		str = sb.String()
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *WhileStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("while ( %s ) %s", s.expr, s.body)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *ForStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
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
		str = sb.String()
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *ExpressionStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("%s ;", s.expr)

	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *PrintStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("print %s ;", s.expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *DeclarationStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("var %s = %s ;", s.name, s.expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *BlockStatement) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		var sb strings.Builder
		sb.WriteString("{ ")
		for _, stmt := range s.stmts {
			fmt.Fprint(&sb, " ", stmt.String())
		}
		sb.WriteString(" } ;")
		str = sb.String()
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (e *UnaryExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *BinaryExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *GroupingExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("(group %s)", e.expr)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *AssignmentExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("(%s = %s)", e.name, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *VariableExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("Var(%s)", e.name)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *CallExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		var args []string
		for _, arg := range e.args {
			args = append(args, arg.String())
		}
		str = fmt.Sprintf("%s(%s)", e.callee, strings.Join(args, ", "))
	default:
		err = UnprintableError{e}
	}
	return str, err
}
func (e *StringExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("\"%s\"", string(e.value))
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *NumericExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprint(e.value)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *BooleanExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		if e.value {
			str = "true"
		} else {
			str = "false"
		}
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *NilExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = "nil"
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (v ValueString) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = fmt.Sprintf("\"%s\"", string(v))
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueNumeric) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueBoolean) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		if bool(v) {
			str = "true"
		} else {
			str = "false"
		}
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueNil) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *DefaultPrinter:
		str = "nil"
	default:
		err = UnprintableError{v}
	}
	return str, err
}
