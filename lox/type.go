package lox

//go:generate stringer -type Type  -trimprefix=Type
type Type int

const (
	ErrType Type = iota
	TypeNil
	TypeNumeric
	TypeString
	TypeBoolean
)
