package lox

//go:generate stringer -type Type  -trimprefix=Type
type Type uint64

const (
	TypeAny Type = iota
	TypeNil
	TypeNumeric
	TypeString
	TypeBoolean
)

type TypeSet struct {
	bits uint64
}

func (ts *TypeSet) Test(t Type) bool {
	return ts.bits>>uint64(t)&uint64(1) != 0
}

func (ts *TypeSet) Set(t Type) {
	ts.bits |= 1 << uint(t)
}

func (ts *TypeSet) Clear(t Type) {
	ts.bits &= ^(1 << uint(t))
}

func (ts *TypeSet) Zero() {
	ts.bits = 0
}
