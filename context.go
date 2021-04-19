package template

import (
	"errors"
	"fmt"
	"github.com/ImSingee/dt"
	"go/ast"
	"go/parser"
	"reflect"
	"regexp"
)

var re = regexp.MustCompile(`@(@|(?P<call>(?P<function>[a-zA-Z0-9\-_]+)(?P<args>\(.*?\))))`)

type Context struct {
	Option Option

	values       map[string]interface{}
	parent       *Context
	functionSets []FunctionSet
}

func NewContext() Context {
	return Context{
		values: make(map[string]interface{}, 8),
	}
}

func (ctx *Context) With(functionSets ...FunctionSet) Context {
	c := NewContext()
	c.parent = ctx
	c.functionSets = functionSets

	return c
}

func (ctx *Context) AddFunctionSet(functionSet FunctionSet) {
	ctx.functionSets = append(ctx.functionSets, functionSet)
}

func (ctx *Context) Execute(source string) (s string, err error) {
	defer func() {
		e := recover()
		if ee, ok := e.(error); ok {
			err = ee
			return
		} else if e != nil {
			err = fmt.Errorf("unknown Error (%v)", e)
		}
	}()

	return re.ReplaceAllStringFunc(source, func(s string) string {
		// 转义
		if s == "@@" {
			return "@"
		}

		// 函数调用
		f := s[1:]

		tree, err := parser.ParseExpr(f)

		if err != nil {
			panic(fmt.Errorf("parse error (%s) (%w)", f, err))
			return ""
		}
		//fmt.Printf("%s: %#+v\n", f, tree)

		call, ok := tree.(*ast.CallExpr)
		if !ok {
			panic(fmt.Errorf("parse as Call error (%s) (%T)", f, tree))
			return ""
		}

		funcName := call.Fun.(*ast.Ident).String()
		args, err := mapArgs(call.Args)

		if err != nil {
			panic(err)
		}

		result, err := ctx.callFunction(funcName, args, ctx.Option)
		if err != nil {
			panic(err)
		}

		return dt.ToString(result)
	}), nil
}

func (ctx *Context) callFunction(funcName string, args []dt.Value, option Option) (interface{}, error) {
	functionArgsNotMatchMatchErr := (error)(nil)
	functionRunError := (error)(nil)

	if len(ctx.functionSets) != 0 {
		for _, fs := range ctx.functionSets {
			result, err := fs.callFunction(funcName, args, option.Greedy)

			if err == nil {
				return result, nil
			}

			if errors.Is(err, errNoFunction) {
				continue
			} else if errors.Is(err, &ErrArgsNotMatch{}) {
				functionArgsNotMatchMatchErr = err
			} else {
				functionRunError = err
			}
		}
	}

	if ctx.parent != nil {
		result, err := ctx.parent.callFunction(funcName, args, option)

		if err == nil {
			return result, nil
		}

		if errors.Is(err, errNoFunction) {
			// do nothing
		} else if errors.Is(err, &ErrArgsNotMatch{}) {
			functionArgsNotMatchMatchErr = err
		} else {
			functionRunError = err
		}
	}

	if functionRunError != nil {
		return nil, functionRunError
	} else if functionArgsNotMatchMatchErr != nil {
		return nil, functionArgsNotMatchMatchErr
	} else {
		return nil, errNoFunction
	}
}

func (ctx *Context) Apply(source interface{}) (err error) {
	return ctx.apply(reflect.ValueOf(source))
}

func (ctx *Context) apply(source reflect.Value) (err error) {
	if source.Kind() != reflect.Ptr || source.IsNil() {
		return errInvalidType
	}

	switch source.Elem().Kind() {
	case reflect.String:
		return ctx.applyOnString(source.Interface().(*string))
	case reflect.Struct:
		return ctx.applyOnStruct(source)
	default:
		return errInvalidType
	}
}

func (ctx *Context) applyOnString(source *string) (err error) {
	result, err := ctx.Execute(*source)
	if err != nil {
		return err
	}
	*source = result
	return nil
}

func (ctx *Context) applyOnStruct(source reflect.Value) (err error) {
	e := source.Elem()
	numField := e.NumField()
	for i := 0; i < numField; i++ {
		ev := e.Field(i)
		err = nil

		if ev.Kind() == reflect.Ptr {
			if !ev.IsNil() {
				err = ctx.apply(ev)
			}
		} else {
			if ev.CanAddr() {
				err = ctx.apply(ev.Addr())
			}
		}

		if err != nil && err != errInvalidType {
			return err
		}
	}

	return nil
}
