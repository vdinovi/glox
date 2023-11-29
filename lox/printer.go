package lox

import (
	"fmt"
	"strconv"
	"strings"
)

var defaultPrinter = CompactPrinter{}

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

type CompactPrinter struct{}

func (p *CompactPrinter) Print(elem Printable) (string, error) {
	return elem.Print(p)
}

// type CorrectPrinter struct {
// 	depth int
// }

// func (p *CorrectPrinter) Print(elem Printable) (string, error) {
// 	return elem.Print(p)
// }

// func (p *CorrectPrinter) pad(indent int) string {
// 	return strings.Repeat("\t", p.depth+indent)
// }

func (s *ConditionalStatement) Print(p Printer) (str string, err error) {
	cond, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	thenBranch, err := s.thenBranch.Print(p)
	if err != nil {
		return "", err
	}
	elseBranch := "<none>"
	if s.elseBranch != nil {
		if elseBranch, err = s.elseBranch.Print(p); err != nil {
			return "", err
		}
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("if (%s) { %s } else { %s }", cond, thenBranch, elseBranch)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *WhileStatement) Print(p Printer) (str string, err error) {
	cond, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	body, err := s.body.Print(p)
	if err != nil {
		return "", err
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("while (%s) { %s }", cond, body)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *ForStatement) Print(p Printer) (str string, err error) {
	init := "<none>"
	if s.init != nil {
		if init, err = s.init.Print(p); err != nil {
			return "", err
		}
	}
	cond := "<none>"
	if s.cond != nil {
		if cond, err = s.cond.Print(p); err != nil {
			return "", err
		}
	}
	incr := "<none>"
	if s.incr != nil {
		if incr, err = s.incr.Print(p); err != nil {
			return "", err
		}
	}
	body, err := s.body.Print(p)
	if err != nil {
		return "", err
	}

	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("for (%s; %s; %s) { %s }", init, cond, incr, body)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *ExpressionStatement) Print(p Printer) (str string, err error) {
	expr, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("%s;", expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *PrintStatement) Print(p Printer) (str string, err error) {
	expr, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("print %s;", expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *DeclarationStatement) Print(p Printer) (str string, err error) {
	expr, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("var %s = %s;", s.name, expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *FunctionDefinitionStatement) Print(p Printer) (str string, err error) {
	body := make([]string, len(s.body))
	for i, stmt := range s.body {
		body[i], err = stmt.Print(p)
		if err != nil {
			return "", err
		}
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("fun %s(%s) { %s }", s.name, strings.Join(s.params, ", "), strings.Join(body, " "))
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *ReturnStatement) Print(p Printer) (str string, err error) {
	expr, err := s.expr.Print(p)
	if err != nil {
		return "", err
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("return %s;", expr)
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (s *BlockStatement) Print(p Printer) (str string, err error) {
	body := make([]string, len(s.stmts))
	for i, stmt := range s.stmts {
		body[i], err = stmt.Print(p)
		if err != nil {
			return "", err
		}
	}
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("{ %s }", strings.Join(body, " "))
	default:
		err = UnprintableError{s}
	}
	return str, err
}

func (e *UnaryExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("(%s %s)", e.op.Lexem, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *BinaryExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("(%s %s %s)", e.op.Lexem, e.left, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *GroupingExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("(group %s)", e.expr)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *AssignmentExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("(%s = %s)", e.name, e.right)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *VariableExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("Var(%s)", e.name)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *CallExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
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
	case *CompactPrinter:
		str = fmt.Sprintf("\"%s\"", string(e.value))
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *NumericExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprint(e.value)
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (e *BooleanExpression) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
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
	case *CompactPrinter:
		str = "nil"
	default:
		err = UnprintableError{e}
	}
	return str, err
}

func (v ValueString) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("\"%s\"", string(v))
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueNumeric) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = strconv.FormatFloat(float64(v), 'f', -1, 64)
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueBoolean) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
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
	case *CompactPrinter:
		str = "nil"
	default:
		err = UnprintableError{v}
	}
	return str, err
}

func (v ValueCallable) Print(p Printer) (str string, err error) {
	switch p.(type) {
	case *CompactPrinter:
		str = fmt.Sprintf("Callable(%s)", v.name)
	default:
		err = UnprintableError{v}
	}
	return str, err
}
