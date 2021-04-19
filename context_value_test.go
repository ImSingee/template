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

func TestValueInFunc(t *testing.T) {
	ctx := NewContext()
	ctx.SetValue("a", 1)

	set := NewFunctionSet()
	set.MustRegisterFunction("test1", func(ctx *Context, i int) int {
		return i + ctx.GetValue("a").(int)
	})
	set.MustRegisterFunction("test2", func(i int, ctx *Context) int {
		return i + ctx.GetValue("a").(int)
	})

	ctx.AddFunctionSet(set)

	value, err := ctx.Execute("@test1(2)")
	tt.AssertIsNil(t, err)
	tt.AssertEqual(t, "3", value)

	value, err = ctx.Execute("@test2(5)")
	tt.AssertIsNil(t, err)
	tt.AssertEqual(t, "6", value)
}
