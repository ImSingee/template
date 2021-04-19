package template

import (
	"github.com/ImSingee/tt"
	"testing"
)

func TestWith(t *testing.T) {
	ctx := NewContext()

	_, err := ctx.Execute("@test()")
	tt.AssertIsNotNil(t, err)

	set := NewFunctionSet()
	set.MustRegisterFunction("test", func() string {
		return "test"
	})
	newCtx := ctx.With(set)

	content, err := newCtx.Execute("@test()")
	tt.AssertIsNil(t, err)
	tt.AssertEqual(t, "test", content)

	_, err = ctx.Execute("@test()")
	tt.AssertIsNotNil(t, err)
}
