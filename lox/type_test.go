package lox

import (
	"testing"
)

func TestTypeSet(t *testing.T) {
	ts := &TypeSet{}
	assertNotSet(t, ts, TypeAny, TypeNil, TypeNumeric, TypeString, TypeBoolean)
	ts.Set(TypeNil)
	assertSet(t, ts, TypeNil)
	assertNotSet(t, ts, TypeAny, TypeNumeric, TypeString, TypeBoolean)
	ts.Set(TypeString)
	assertSet(t, ts, TypeNil, TypeString)
	assertNotSet(t, ts, TypeAny, TypeNumeric, TypeBoolean)
	ts.Set(TypeBoolean)
	assertSet(t, ts, TypeNil, TypeString, TypeBoolean)
	assertNotSet(t, ts, TypeAny, TypeNumeric)
	ts.Clear(TypeString)
	assertSet(t, ts, TypeNil, TypeBoolean)
	assertNotSet(t, ts, TypeAny, TypeString, TypeNumeric)
	ts.Zero()
	assertNotSet(t, ts, TypeAny, TypeNil, TypeNumeric, TypeString, TypeBoolean)
}

func assertSet(t *testing.T, ts *TypeSet, types ...Type) {
	t.Helper()
	for _, typ := range types {
		if !ts.Test(typ) {
			t.Errorf("Expected %s to be set but was not", typ)
			return
		}
	}
}

func assertNotSet(t *testing.T, ts *TypeSet, types ...Type) {
	t.Helper()
	for _, typ := range types {
		if ts.Test(typ) {
			t.Errorf("Expected %s to not be set but was", typ)
			return
		}
	}
}
