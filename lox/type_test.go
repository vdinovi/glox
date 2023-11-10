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
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualTo], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEqualTo], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpLessThan], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpLessThanOrEqualTo], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpGreaterThan], left: one, right: one}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpGreaterThanOrEqualTo], left: one, right: one}},
// 			// String
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpPlus], left: str, right: str}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualTo], left: str, right: str}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEqualTo], left: str, right: str}},
// 			// Bool
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpEqualTo], left: yes, right: yes}},
// 			ExpressionStatement{expr: BinaryExpression{op: operators[OpNotEqualTo], left: yes, right: yes}},
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
// 	OpEqualTo:   Operator{OpEqualTo, "=="},
// 	OpNotEqualTo:     Operator{OpNotEqualTo, "!="},
// 	OpLessThan:          Operator{OpLessThan, "<"},
// 	OpLessThanOrEqualTo:    Operator{OpLessThanOrEqualTo, "<="},
// 	OpGreaterThan:       Operator{OpGreaterThan, ">"},
// 	OpGreaterThanOrEqualTo: Operator{OpGreaterThan, ">="},
// }
