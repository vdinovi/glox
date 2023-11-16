package lox

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParserExpressionStatement(t *testing.T) {
	tests := []struct {
		text  string
		stmts []ExpressionStatement
		err   error
	}{
		{text: "1;", stmts: []ExpressionStatement{{expr: oneExpr()}}},
		{text: "3.14;", stmts: []ExpressionStatement{{expr: piExpr()}}},
		{text: "\"str\";", stmts: []ExpressionStatement{{expr: strExpr()}}},
		{text: "true;", stmts: []ExpressionStatement{{expr: trueExpr()}}},
		{text: "false;", stmts: []ExpressionStatement{{expr: falseExpr()}}},
		{text: "nil;", stmts: []ExpressionStatement{{expr: nilExpr()}}},
		{text: "//comment\n1;", stmts: []ExpressionStatement{{expr: oneExpr()}}},
		{text: "1;//comment\n", stmts: []ExpressionStatement{{expr: oneExpr()}}},
		{text: "foo;", stmts: []ExpressionStatement{{expr: fooExpr()}}},
		{text: "-1;", stmts: []ExpressionStatement{{expr: uSubExpr(oneExpr())()}}},
		{text: "--1;", stmts: []ExpressionStatement{{expr: uSubExpr(uSubExpr(oneExpr())())()}}},
		{text: "+1;", stmts: []ExpressionStatement{{expr: uAddExpr(oneExpr())()}}},
		{text: "++1;", stmts: []ExpressionStatement{{expr: uAddExpr(uAddExpr(oneExpr())())()}}},
		{text: "(1);", stmts: []ExpressionStatement{{expr: groupExpr(oneExpr())()}}},
		{text: "(-1);", stmts: []ExpressionStatement{{expr: groupExpr(uSubExpr(oneExpr())())()}}},
		{text: "1 + 3.14;", stmts: []ExpressionStatement{{expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "1 - -3.14;", stmts: []ExpressionStatement{{expr: bSubExpr(oneExpr())(uSubExpr(piExpr())())()}}},
		{text: "-1 * 3.14;", stmts: []ExpressionStatement{{expr: bMulExpr(uSubExpr(oneExpr())())(piExpr())()}}},
		{text: "-1 / -3.14;", stmts: []ExpressionStatement{{expr: bDivExpr(uSubExpr(oneExpr())())(uSubExpr(piExpr())())()}}},
		{text: "(1 + 3.14);", stmts: []ExpressionStatement{{expr: groupExpr(bAddExpr(oneExpr())(piExpr())())()}}},
		{text: "1 + (1 + 3.14);", stmts: []ExpressionStatement{{expr: bAddExpr(oneExpr())(groupExpr(bAddExpr(oneExpr())(piExpr())())())()}}},
		{text: "(1 + 3.14) + 1;", stmts: []ExpressionStatement{{expr: bAddExpr(groupExpr(bAddExpr(oneExpr())(piExpr())())())(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParserPrintStatement(t *testing.T) {
	tests := []struct {
		text  string
		stmts []PrintStatement
		err   error
	}{
		{text: "print 1;", stmts: []PrintStatement{{expr: oneExpr()}}},
		{text: "print 3.14;", stmts: []PrintStatement{{expr: piExpr()}}},
		{text: "print \"str\";", stmts: []PrintStatement{{expr: strExpr()}}},
		{text: "print true;", stmts: []PrintStatement{{expr: trueExpr()}}},
		{text: "print false;", stmts: []PrintStatement{{expr: falseExpr()}}},
		{text: "print nil;", stmts: []PrintStatement{{expr: nilExpr()}}},
		{text: "//comment\nprint 1;", stmts: []PrintStatement{{expr: oneExpr()}}},
		{text: "print 1;//comment\n", stmts: []PrintStatement{{expr: oneExpr()}}},
		{text: "print foo;", stmts: []PrintStatement{{expr: fooExpr()}}},
		{text: "print -1;", stmts: []PrintStatement{{expr: uSubExpr(oneExpr())()}}},
		{text: "print 1 + 3.14;", stmts: []PrintStatement{{expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "print (1);", stmts: []PrintStatement{{expr: groupExpr(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParserDeclarationStatement(t *testing.T) {
	tests := []struct {
		text  string
		stmts []DeclarationStatement
		err   error
	}{
		{text: "var foo;", stmts: []DeclarationStatement{{name: "foo", expr: nilExpr()}}},
		{text: "var foo = 1;", stmts: []DeclarationStatement{{name: "foo", expr: oneExpr()}}},
		{text: "var foo = 3.14;", stmts: []DeclarationStatement{{name: "foo", expr: piExpr()}}},
		{text: "var foo = \"str\";", stmts: []DeclarationStatement{{name: "foo", expr: strExpr()}}},
		{text: "var foo = true;", stmts: []DeclarationStatement{{name: "foo", expr: trueExpr()}}},
		{text: "var foo = false;", stmts: []DeclarationStatement{{name: "foo", expr: falseExpr()}}},
		{text: "var foo = nil;", stmts: []DeclarationStatement{{name: "foo", expr: nilExpr()}}},
		{text: "//comment\nvar foo;", stmts: []DeclarationStatement{{name: "foo", expr: nilExpr()}}},
		{text: "var foo;//comment\n", stmts: []DeclarationStatement{{name: "foo", expr: nilExpr()}}},
		{text: "var foo = 1 + 3.14;", stmts: []DeclarationStatement{{name: "foo", expr: bAddExpr(oneExpr())(piExpr())()}}},
		{text: "var foo = (1);", stmts: []DeclarationStatement{{name: "foo", expr: groupExpr(oneExpr())()}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %q, but got %q", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParserBlockStatement(t *testing.T) {
	tests := []struct {
		text  string
		stmts []BlockStatement
		err   error
	}{
		{text: "{var foo;}", stmts: []BlockStatement{{stmts: []Statement{&DeclarationStatement{name: "foo", expr: nilExpr()}}}}},
		{text: "{1;}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}}}}},
		{text: "{print 1;}", stmts: []BlockStatement{{stmts: []Statement{&PrintStatement{expr: oneExpr()}}}}},
		{text: "{1; 1;}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}, &ExpressionStatement{expr: oneExpr()}}}}},
		{text: "{1; {1;}}", stmts: []BlockStatement{{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}, &BlockStatement{stmts: []Statement{&ExpressionStatement{expr: oneExpr()}}}}}}},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}

		for i, want := range test.stmts {
			if got := program[i]; !got.Equals(&want) {
				t.Errorf("Expected %q to be %s, but got %s", test.text, want.String(), got.String())
				break
			}
		}
	}
}

func TestParserConditionalStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt ConditionalStatement
		err  error
	}{
		// if/else without braces
		{
			text: "if (true) 1; else 3.14;",
			stmt: ConditionalStatement{
				expr:       trueExpr(),
				thenBranch: &ExpressionStatement{expr: oneExpr()},
				elseBranch: &ExpressionStatement{expr: piExpr()},
			},
		},
		// if without else
		{
			text: "if (true) 1;",
			stmt: ConditionalStatement{
				expr:       trueExpr(),
				thenBranch: &ExpressionStatement{expr: oneExpr()},
				elseBranch: nil,
			},
		},
		// if with braces without else
		{
			text: "if (true) {1;}",
			stmt: ConditionalStatement{
				expr: trueExpr(),
				thenBranch: &BlockStatement{
					stmts: []Statement{
						&ExpressionStatement{expr: oneExpr()},
					},
				},
				elseBranch: nil,
			},
		},
		// nested if without braces (else binds to closest if)
		{
			text: "if (true) if (false) 1; else 3.14; else 1;",
			stmt: ConditionalStatement{
				expr: trueExpr(),
				thenBranch: &ConditionalStatement{
					expr:       falseExpr(),
					thenBranch: &ExpressionStatement{expr: oneExpr()},
					elseBranch: &ExpressionStatement{expr: piExpr()},
				},
				elseBranch: &ExpressionStatement{
					expr: oneExpr(),
				},
			},
		},
		// nested if with braces (changes else binding)
		{
			text: "if (true) { if (false) 1; else 3.14;} else {1;}",
			stmt: ConditionalStatement{
				expr: trueExpr(),
				thenBranch: &BlockStatement{
					stmts: []Statement{
						&ConditionalStatement{
							expr:       falseExpr(),
							thenBranch: &ExpressionStatement{expr: oneExpr()},
							elseBranch: &ExpressionStatement{expr: piExpr()},
						},
					},
				},
				elseBranch: &BlockStatement{
					stmts: []Statement{
						&ExpressionStatement{expr: oneExpr()},
					},
				},
			},
		},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		if len(program) != 1 {
			t.Errorf("Expected %q to produce 1 statement but got %d", test.text, len(program))
		}
		stmt := program[0]
		if !stmt.Equals(&test.stmt) {
			t.Errorf("Expected %q to be %q, but got %q", test.text, test.stmt.String(), stmt.String())
		}
	}
}

func TestParserWhileStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt WhileStatement
		err  error
	}{
		{
			text: "while (true) 1;",
			stmt: WhileStatement{
				expr: trueExpr(),
				body: &ExpressionStatement{expr: oneExpr()},
			},
		},
		{
			text: "while (true) { 1; }",
			stmt: WhileStatement{
				expr: trueExpr(),
				body: &BlockStatement{
					stmts: []Statement{
						&ExpressionStatement{expr: oneExpr()},
					},
				},
			},
		},
		{
			text: "while (true) { while(false) 1; }",
			stmt: WhileStatement{
				expr: trueExpr(),
				body: &BlockStatement{
					stmts: []Statement{
						&WhileStatement{
							expr: falseExpr(),
							body: &ExpressionStatement{
								expr: oneExpr(),
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		if len(program) != 1 {
			t.Errorf("Expected %q to produce 1 statement but got %d", test.text, len(program))
		}
		stmt := program[0]
		if !stmt.Equals(&test.stmt) {
			t.Errorf("Expected %q to be %q, but got %q", test.text, test.stmt.String(), stmt.String())
		}
	}
}

func TestParserForStatement(t *testing.T) {
	tests := []struct {
		text string
		stmt ForStatement
		err  error
	}{
		{
			text: "for (;;) 1;",
			stmt: ForStatement{body: &ExpressionStatement{expr: oneExpr()}},
		},
		{
			text: "for (;;) { 1; }",
			stmt: ForStatement{
				body: &BlockStatement{
					stmts: []Statement{
						&ExpressionStatement{expr: oneExpr()},
					},
				},
			},
		},
		{
			text: "for (var x = 1;;) 1;",
			stmt: ForStatement{
				init: &DeclarationStatement{name: "x", expr: oneExpr()},
				body: &ExpressionStatement{expr: oneExpr()},
			},
		},
		{
			text: "for (var x = 1; x;) 1;",
			stmt: ForStatement{
				init: &DeclarationStatement{name: "x", expr: oneExpr()},
				cond: &VariableExpression{name: "x"},
				body: &ExpressionStatement{expr: oneExpr()},
			},
		},
		{
			text: "for (var x = 1; x; x = x + 1) 1;",
			stmt: ForStatement{
				init: &DeclarationStatement{name: "x", expr: oneExpr()},
				cond: &VariableExpression{name: "x"},
				incr: &AssignmentExpression{
					name:  "x",
					right: &BinaryExpression{op: addOp, left: makeVarExpr("x")(), right: oneExpr()},
				},
				body: &ExpressionStatement{expr: oneExpr()},
			},
		},
	}
	for _, test := range tests {
		tokens, err := Scan(strings.NewReader(test.text))
		if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		program, err := Parse(tokens)
		if test.err != nil {
			if err != test.err {
				t.Errorf("Expected %q to produce error %q, but got %q", test.text, test.err, err)
				continue
			}
		} else if err != nil {
			t.Errorf("Unexpected error in %q: %s", test.text, err)
			continue
		}
		if len(program) != 1 {
			t.Errorf("Expected %q to produce 1 statement but got %d", test.text, len(program))
		}
		stmt := program[0]
		if !stmt.Equals(&test.stmt) {
			t.Errorf("Expected %q to be %q, but got %q", test.text, test.stmt.String(), stmt.String())
		}
	}
}

func TestParserProgram(t *testing.T) {
	// TODO: Needs to serialize AST to golden file for this test to work
	t.Skip()
	var tokens []Token
	err := json.Unmarshal([]byte(fixtureProgramTokens), &tokens)
	if err != nil {
		t.Fatalf("Failed to deserialize tokens")
	}

	parser := NewParser(tokens)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	_, err = parser.Parse()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func BenchmarkParserFixturePrograms(b *testing.B) {
	tokens, err := Scan(strings.NewReader(fixtureProgram))
	if err != nil {
		b.Errorf("Unexpected error lexing fixture 'program': %s", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(tokens)
		if err != nil {
			b.Errorf("Unexpected error parsing fixture 'program': %s", err)
		}
	}
}
