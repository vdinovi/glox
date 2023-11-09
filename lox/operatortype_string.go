// Code generated by "stringer -type OperatorType -trimprefix=Op"; DO NOT EDIT.

package lox

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpPlus-0]
	_ = x[OpMinus-1]
	_ = x[OpMultiply-2]
	_ = x[OpDivide-3]
	_ = x[OpEquals-4]
	_ = x[OpEqualEquals-5]
	_ = x[OpNotEquals-6]
	_ = x[OpLess-7]
	_ = x[OpLessEquals-8]
	_ = x[OpGreater-9]
	_ = x[OpGreaterEquals-10]
}

const _OperatorType_name = "PlusMinusMultiplyDivideEqualsEqualEqualsNotEqualsLessLessEqualsGreaterGreaterEquals"

var _OperatorType_index = [...]uint8{0, 4, 9, 17, 23, 29, 40, 49, 53, 63, 70, 83}

func (i OperatorType) String() string {
	if i < 0 || i >= OperatorType(len(_OperatorType_index)-1) {
		return "OperatorType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OperatorType_name[_OperatorType_index[i]:_OperatorType_index[i+1]]
}
