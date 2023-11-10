package lox

// var one = NumericExpression(1)
// var str = StringExpression("str")
// var yes = BooleanExpression(true)

// func TestTypeCheckUnaryOps(t *testing.T) {
// 	tests := [][]Statement{
// 		{
// 			// Numeric
// 			ExpressionStatement{expr: UnaryExpression{op: operators[OpMinus], right: one}},
// 		},
// 	}

// 	for _, test := range tests {
// 		err := TypeCheck(test)
// 		if err != nil {
// 			t.Errorf("Unexpected error: %s", err)
// 			continue
// 		}
// 	}

//}

// func TestTypeCheckBinaryOps(t *testing.T) {
// 	tests := [][]Statement{
// 		{
// 			// Numeric
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpPlus], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpMinus], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpMultiply], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpDivide], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualEquals], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEquals], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpLess], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpLessEquals], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpGreater], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpGreaterEquals], left: one, right: one}},
// 			// String
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpPlus], left: str, right: str}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualEquals], left: str, right: str}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEquals], left: str, right: str}},
// 			// Bool
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualEquals], left: yes, right: yes}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEquals], left: yes, right: yes}},
// 		},
// 	}

// 	for _, test := range tests {
// 		err := TypeCheck(test)
// 		if err != nil {
// 			t.Errorf("Unexpected error: %s", err)
// 			continue
// 		}
// 	}
// }

// func TestTypeCheckFail(t *testing.T) {
// 	one := NumericExpression(1)

// 	tests := [][]Statement{
// 		{
// 			ExpressionStatement{expr: UnaryExpression{op: Operator{OpMinus, "-"}, right: one}},
// 		},
// 	}

// 	var tErr *ErrTypeor
// 	for _, test := range tests {
// 		err := TypeCheck(test)
// 		if err == nil {
// 			t.Errorf("Unexpected error: %s", err)
// 			continue
// 		}
// 	}
// }

// var operators = map[OperatorType]Operator{
// 	OpMinus:         Operator{OpMinus, "-"},
// 	OpPlus:          Operator{OpPlus, "+"},
// 	OpMultiply:      Operator{OpMultiply, "*"},
// 	OpDivide:        Operator{OpDivide, "/"},
// 	OpEqualEquals:   Operator{OpEqualEquals, "=="},
// 	OpNotEquals:     Operator{OpNotEquals, "!="},
// 	OpLess:          Operator{OpLess, "<"},
// 	OpLessEquals:    Operator{OpLessEquals, "<="},
// 	OpGreater:       Operator{OpGreater, ">"},
// 	OpGreaterEquals: Operator{OpGreater, ">="},
// }
