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

type Typed interface {
	Type(Symbols) (Type, error)
}

// func typeFor(v any) Type {
// 	if v == nil {
// 		return TypeNil
// 	}
// 	if t, ok := v.(Type); ok {
// 		return t
// 	}
// 	switch v.(type) {
// 	case float64:
// 		return TypeNumeric
// 	case string:
// 		return TypeString
// 	case bool:
// 		return TypeBoolean
// 	}
// 	return ErrType
// }
