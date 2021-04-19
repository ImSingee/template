package template

import (
	"github.com/ImSingee/tt"
	"testing"
)

func TestContextValue(t *testing.T) {
	ctx := NewContext()
	ctx.SetValue("a", 1)
	ctx.SetValue("b", 1)
	newCtx := ctx.With()
	newCtx.SetValue("b", 2)
	newCtx.SetValue("c", 2)

	tt.AssertEqual(t, 1, ctx.GetValue("a"))
	tt.AssertEqual(t, 1, ctx.GetValue("b"))
	tt.AssertEqual(t, nil, ctx.GetValue("c"))
	tt.AssertEqual(t, nil, ctx.GetValue("d"))

	tt.AssertEqual(t, 1, newCtx.GetValue("a"))
	tt.AssertEqual(t, 2, newCtx.GetValue("b"))
	tt.AssertEqual(t, 2, newCtx.GetValue("c"))
	tt.AssertEqual(t, nil, newCtx.GetValue("d"))
}
