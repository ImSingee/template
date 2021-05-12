package template

import (
	"fmt"
	"github.com/ImSingee/dt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

type function struct {
	F          interface{}
	Args       []argument
	ArgsNumber int
	MayBeError bool // 返回值是 error
}

func (f *function) Apply(args []dt.Value, callerCtx *Context) ([]reflect.Value, error) {
	if len(args) != f.ArgsNumber {
		return nil, fmt.Errorf("args number mismatch")
	}

	mappedArgs := make([]reflect.Value, 0, len(f.Args))
	j := 0
	for i := 0; i < len(f.Args); i++ {
		if f.Args[i].IsContext {
			mappedArgs = append(mappedArgs, reflect.ValueOf(callerCtx))
		} else {
			value, ok := dt.ConvertToReflectType(args[j], f.Args[i].OutType)
			if !ok {
				return nil, fmt.Errorf("arg %d's out-type mismatch", i+1)
			}
			mappedArgs = append(mappedArgs, value)
			j += 1
		}
	}

	return mappedArgs, nil
}

type argument struct {
	Name string

	IsContext bool // 应注入当前 Context

	InType  dt.Type      // 用户字面书写的转换后参数类型，支持 string, bool, *GenericNumber
	OutType reflect.Kind // 函数原型参数类型，支持 string, bool, int[X], uint[X]
	// 暂不支持可变参数
}

func mapArgs(args []ast.Expr) ([]dt.Value, error) {
	mapped := make([]dt.Value, len(args))
	var err error

	for i := range args {
		//fmt.Printf("param %d: %#+v\n", i, args[i])

		switch arg := args[i].(type) {
		case *ast.BasicLit:
			mapped[i], err = mapArg(arg.Value)
		case *ast.Ident:
			mapped[i], err = mapArg(arg.String())
		}

		if err != nil {
			return nil, fmt.Errorf("invalid param (E03)\n")
		}
	}

	return mapped, nil
}

func mapArg(arg string) (interface{}, error) {
	// 关键字（*ast.Ident）   <- bool, nil
	switch arg {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "nil":
		return nil, nil
	}

	// 字面值量 (*ast.BasicLit) <- number, string

	if strings.HasPrefix(arg, `'`) || strings.HasPrefix(arg, `"`) || strings.HasPrefix(arg, "`") {
		// 字符串（或字符）
		return strconv.Unquote(arg)
	}

	if num, ok := dt.NumberFromString(arg); ok {
		return num, nil
	}

	return nil, fmt.Errorf("invalid param (E5) V=%s", arg)
}
