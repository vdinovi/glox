package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type Typed interface {
	Type() Type
}

const (
	typeNilBit = 1 << iota
	typeByteBit
	typeBooleanBit
	typeUint64Bit
	typeInt64Bit
	typeFloat64Bit
	typeUtf8Bit
	typeStringBit
	typeCallable
)

var TypeAny = Type{bits: ^uint(0)}
var TypeNone = Type{bits: uint(0)}

var TypeNil = Type{bits: uint(typeNilBit)}
var TypeBoolean = Type{bits: uint(typeBooleanBit)}
var TypeNumeric = Type{bits: uint(typeUint64Bit | typeInt64Bit | typeFloat64Bit)}
var TypeString = Type{bits: uint(typeStringBit)}
var TypeCallable = Type{bits: uint(typeCallable)}

var allTypes = [...]Type{TypeNil, TypeBoolean, TypeNumeric, TypeString, TypeCallable}
var typeStrings = [...]string{"TypeNil", "TypeBoolean", "TypeNumeric", "TypeString", "TypeCallable"}

type Type struct {
	bits uint
}

func (t Type) String() string {
	switch t {
	case TypeAny:
		return "Any"
	case TypeNone:
		return "None"
	}
	rem := t
	ts := []string{}
	for i, v := range allTypes {
		if t.Contains(v) {
			rem.Clear(v)
			ts = append(ts, typeStrings[i])
		}
	}
	if rem.bits != 0 {
		ts = append(ts, fmt.Sprintf("rem=%s", strconv.FormatUint(uint64(rem.bits), 2)))
	}
	return fmt.Sprintf("Type{%s}", strings.Join(ts, ", "))
}

func (t Type) Union(u Type) Type {
	return Type{bits: t.bits | u.bits}
}

func (t Type) Subtract(u Type) Type {
	return Type{bits: t.bits &^ u.bits}
}

func (t Type) Contains(us ...Type) bool {
	x := TypeNone
	for _, u := range us {
		x.Set(u)
	}
	return t.bits&x.bits == x.bits
}

func (t Type) Within(us ...Type) bool {
	x := TypeNone
	for _, u := range us {
		x.Set(u)
	}
	return t.bits&x.bits == t.bits
}

func (t Type) Test(u Type) bool {
	return t.bits&u.bits != 0
}

func (t *Type) Set(u Type) {
	t.bits |= u.bits
}

func (t *Type) Clear(u Type) {
	t.bits &^= u.bits
}

func (t *Type) Zero() {
	t.bits = 0
}
