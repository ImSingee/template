package template

import (
	"github.com/ImSingee/tt"
	"testing"
)

func TestApplyString(t *testing.T) {
	s := "@number(10)"
	tt.AssertIsNil(t, Apply(&s))
	tt.AssertEqual(t, s, "10")
}

type TT struct {
	A string
	B string
	C *string
}

type T struct {
	A string
	B string
	C *string
	D *string
	X TT
	Y *T
	Z *T
}

func TestApplyStruct(t *testing.T) {

	k := T{
		A: "77",
		B: "@number(99)",
	}

	tt.AssertIsNil(t, Apply(&k))
	tt.AssertEqual(t, k.A, "77")
	tt.AssertEqual(t, k.B, "99")
}

func TestApplyNestedStruct(t *testing.T) {
	var p = "88"
	var q = "@number(88)"

	k := T{
		A: "77",
		B: "@number(99)",
		C: &p,
		D: nil,
		X: TT{
			A: "@number(77)",
			B: "88",
			C: &q,
		},
		Y: &T{
			A: "@number(66)",
			B: "@number(33)",
			X: TT{
				A: "@number(66)",
			},
		},
	}

	tt.AssertIsNil(t, Apply(&k))
	tt.AssertEqual(t, k.A, "77")
	tt.AssertEqual(t, k.B, "99")
	tt.AssertEqual(t, *k.C, "88")
	tt.AssertIsNil(t, k.D)
	tt.AssertEqual(t, k.X.A, "77")
	tt.AssertEqual(t, k.X.B, "88")
	tt.AssertEqual(t, *k.X.C, "88")
	tt.AssertEqual(t, k.Y.A, "66")
	tt.AssertEqual(t, k.Y.B, "33")
	tt.AssertEqual(t, k.Y.X.A, "66")
}

func TestApplyAliasString(t *testing.T) {
	type A = string
	var s A
	s = "@number(10)"

	tt.AssertIsNil(t, Apply(&s))
	tt.AssertEqual(t, s, "10")
}

func TestApplyAliasStruct(t *testing.T) {
	type A = string
	type T struct {
		A A
	}

	v := T{
		A: "@number(10)",
	}

	tt.AssertIsNil(t, Apply(&v))
	tt.AssertEqual(t, v.A, "10")
}
