package lox

// func TestExecutor(t *testing.T) {
// 	rt := Executor{}
// 	expr := BinaryExpression{
// 		op: Operator{Type: OpMultiply, Lexem: "*"},
// 		left: UnaryExpression{
// 			op:    Operator{Type: OpSubtract, Lexem: "-"},
// 			right: NumericExpression(123),
// 		},
// 		right: GroupingExpression{
// 			expr: NumericExpression(45.67),
// 		},
// 	}
// 	stmt := ExpressionStatement{expr: expr}
// 	err := rt.Execute(stmt)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %s", err)
// 	}
// }
