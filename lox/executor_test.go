package lox

import (
	"reflect"
	"testing"
)

func TestExecutor(t *testing.T) {
	tests := []struct {
		text   string
		prints []string
		err    error
	}{
		{text: "print 1;", prints: []string{"1"}},
		{text: "print 3.14;", prints: []string{"3.14"}},
		{text: "print \"str\";", prints: []string{"str"}},
		{text: "print true;", prints: []string{"true"}},
		{text: "print false;", prints: []string{"false"}},
		{text: "print nil;", prints: []string{"nil"}},
		{text: "print -3.14;", prints: []string{"-3.14"}},
		{text: "print 1 + 2;", prints: []string{"3"}},
		{text: "print 1; print 2;", prints: []string{"1", "2"}},
		{text: "var x = 1; print x;", prints: []string{"1"}},
		{text: "var x = 1; { var x = 2; print x; } print x;", prints: []string{"2", "1"}},
		{text: "if (true) print 1; else print 2;", prints: []string{"1"}},
		{text: "if (false) print 1; else print 2;", prints: []string{"2"}},
		{text: "var x = 1; if (true) { x = 2; } print x;", prints: []string{"2"}},
		{text: "var x = 1; if (false) { x = 2; } print x;", prints: []string{"1"}},
		{text: "while (false) print 1;", prints: nil},
		{text: "var x = true; while (x) {print 1; x = false;}", prints: []string{"1"}},
		{text: "for (;false;) print 1;", prints: nil},
		{text: "for (var x = true; x; x = false) print 1;", prints: []string{"1"}},
	}
	for _, test := range tests {
		td := NewTestDriver(t, test.text)
		td.Lex()
		td.Parse()
		td.TypeCheck()
		td.Fatal()

		td.Execute()
		err := td.Err
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected execution of %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		if !reflect.DeepEqual(td.Printer.Prints, test.prints) {
			t.Errorf("Expected execution of %q to print %v, but printed %v", test.text, test.prints, td.Printer.Prints)
		}
	}
}
