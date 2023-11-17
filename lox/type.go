package lox

import (
	"fmt"
	"strconv"
)

//go:generate stringer -type Type  -trimprefix=Type

const (
	typeNilBit = 1 << iota
	typeByteBit
	typeBooleanBit
	typeUint64Bit
	typeInt64Bit
	typeFloat64Bit
	typeUtf8Bit
	typeStringBit
)

var TypeAny = Type{bits: ^uint(0)}
var TypeNone = Type{bits: uint(0)}

var TypeNil = Type{bits: uint(typeNilBit)}
var TypeBoolean = Type{bits: uint(typeBooleanBit)}
var TypeNumeric = Type{bits: uint(typeUint64Bit | typeInt64Bit | typeFloat64Bit)}
var TypeString = Type{bits: uint(typeStringBit)}

type Type struct {
	bits uint
}

func (t Type) String() string {
	switch t {
	case TypeAny:
		return "Any"
	case TypeNone:
		return "None"
	case TypeNil:
		return "Nil"
	case TypeBoolean:
		return "Boolean"
	case TypeNumeric:
		return "Numeric"
	case TypeString:
		return "String"
	default:
		return fmt.Sprintf("Type{%s}", strconv.FormatUint(uint64(t.bits), 2))
	}
}

func (t *Type) Accepts(u Type) bool {
	return t.bits&u.bits == u.bits
}

func (t *Type) Test(n uint) bool {
	return t.bits>>n&1 != 0
}

func (t *Type) Set(u uint) {
	t.bits |= 1 << u
}

func (t *Type) Clear(u uint) {
	t.bits &= ^(1 << u)
}

func (t *Type) Zero() {
	t.bits = 0
}
