package lox

//go:generate stringer -type Type  -trimprefix=Type
type Type int

const ErrType Type = -1

const (
	TypeAny Type = iota
	TypeNil
	TypeNumeric
	TypeString
	TypeBoolean
)
