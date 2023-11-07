// Code generated by "stringer -type TokenType -trimprefix=Token"; DO NOT EDIT.

package lox

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenNone-0]
	_ = x[TokenLeftParen-1]
	_ = x[TokenRightParen-2]
	_ = x[TokenLeftBrace-3]
	_ = x[TokenRightBrace-4]
	_ = x[TokenComma-5]
	_ = x[TokenDot-6]
	_ = x[TokenMinus-7]
	_ = x[TokenPlus-8]
	_ = x[TokenSemicolon-9]
	_ = x[TokenSlash-10]
	_ = x[TokenStar-11]
	_ = x[TokenBang-12]
	_ = x[TokenBangEqual-13]
	_ = x[TokenEqual-14]
	_ = x[TokenEqualEqual-15]
	_ = x[TokenGreater-16]
	_ = x[TokenGreaterEqual-17]
	_ = x[TokenLess-18]
	_ = x[TokenLessEqual-19]
	_ = x[TokenIdentifier-20]
	_ = x[TokenString-21]
	_ = x[TokenNumber-22]
	_ = x[TokenAnd-23]
	_ = x[TokenClass-24]
	_ = x[TokenElse-25]
	_ = x[TokenFalse-26]
	_ = x[TokenFun-27]
	_ = x[TokenFor-28]
	_ = x[TokenIf-29]
	_ = x[TokenNil-30]
	_ = x[TokenOr-31]
	_ = x[TokenPrint-32]
	_ = x[TokenReturn-33]
	_ = x[TokenSuper-34]
	_ = x[TokenThis-35]
	_ = x[TokenTrue-36]
	_ = x[TokenVar-37]
	_ = x[TokenWhile-38]
	_ = x[TokenComment-39]
	_ = x[TokenEOF-40]
}

const _TokenType_name = "NoneLeftParenRightParenLeftBraceRightBraceCommaDotMinusPlusSemicolonSlashStarBangBangEqualEqualEqualEqualGreaterGreaterEqualLessLessEqualIdentifierStringNumberAndClassElseFalseFunForIfNilOrPrintReturnSuperThisTrueVarWhileCommentEOF"

var _TokenType_index = [...]uint8{0, 4, 13, 23, 32, 42, 47, 50, 55, 59, 68, 73, 77, 81, 90, 95, 105, 112, 124, 128, 137, 147, 153, 159, 162, 167, 171, 176, 179, 182, 184, 187, 189, 194, 200, 205, 209, 213, 216, 221, 228, 231}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
