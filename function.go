package template

import (
	"fmt"
	"github.com/ImSingee/dt"
	"go/ast"
	"reflect"
	"strings"
)

type function struct {
	F          interface{}
	Args       []argument
	MayBeError bool // 返回值是 error
}

func (f *function) Apply(args []dt.Value) ([]reflect.Value, error) {
	if len(args) != len(f.Args) {
		return nil, fmt.Errorf("args number mismatch")
	}

	mappedArgs := make([]reflect.Value, 0, len(args))
	for i, dtValue := range args {
		value, ok := dt.ConvertToReflectType(dtValue, f.Args[i].OutType)
		if !ok {
			return nil, fmt.Errorf("arg %d's out-type mismatch", i+1)
		}
		mappedArgs = append(mappedArgs, value)
	}

	return mappedArgs, nil
}

type argument struct {
	Name string

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

	if strings.HasPrefix(arg, `'`) || strings.HasPrefix(arg, `"`) {
		if !strings.HasPrefix(arg, `'`) && !strings.HasPrefix(arg, `"`) {
			return nil, fmt.Errorf("invalid literal")
		}
		return arg[1 : len(arg)-1], nil
	}

	if num, ok := dt.NumberFromString(arg); ok {
		return num, nil
	}

	return nil, fmt.Errorf("invalid param (E5) V=%s", arg)
}
