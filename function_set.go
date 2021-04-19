package template

import (
	"fmt"
	"github.com/ImSingee/dt"
	"reflect"
)

type FunctionSet map[string][]function

func NewFunctionSet() FunctionSet {
	return make(FunctionSet, 32)
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// 注册函数
// 参数类型与 RegisterFunction 一致
// 如果注册失败会直接 panic
func (fs FunctionSet) MustRegisterFunction(params ...interface{}) {
	err := fs.RegisterFunction(params...)
	if err != nil {
		panic(err)
	}
}

// 注册函数
// 接收的参数类型：
// - string 类型：函数名
// - function 类型：函数
func (fs FunctionSet) RegisterFunction(params ...interface{}) error {
	names := make([]string, 0, len(params))
	functions := make([]interface{}, 0, len(params))

	for _, param := range params {
		if s, ok := param.(string); ok {
			names = append(names, s)
		} else {
			functions = append(functions, param)
		}
	}

	for _, name := range names {
		for _, f := range functions {
			err := fs.registerFunction(name, f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fs FunctionSet) registerFunction(name string, f interface{}) error {
	ff := reflect.TypeOf(f)
	if ff.Kind() != reflect.Func {
		return errIsNotFunction
	}

	function := function{F: f}

	numOut := ff.NumOut()
	switch numOut {
	case 0:
		return fmt.Errorf("function return at least one value")
	case 1:
		function.MayBeError = false
	case 2:
		out := ff.Out(1)
		if out.Kind() != reflect.Interface {
			return fmt.Errorf("function's second return value must be error interface")
		}

		if !out.Implements(errorInterface) {
			return fmt.Errorf("function's second return value must be error")
		}

		function.MayBeError = true
	default:
		return fmt.Errorf("function return too much values")
	}

	argsCount := ff.NumIn()
	args := make([]argument, argsCount)
	for i := range args {
		arg := ff.In(i)

		args[i].Name = arg.Name()
		args[i].OutType = arg.Kind()
		args[i].InType = dt.MapReflectType(arg.Kind())
	}
	function.Args = args

	fs[name] = append(fs[name], function)

	return nil
}

// 调用函数
// funcName: 函数名
// args：函数参数
// greedy：如果找到了函数但调用返回错误是否继续尝试其他函数
func (fs FunctionSet) callFunction(funcName string, args []dt.Value, greedy bool) (interface{}, error) {
	funcs := fs[funcName]
	if len(funcs) == 0 {
		return nil, errNoFunction
	}

	var latestError ErrArgsNotMatch
	for _, f := range funcs {
		inArgs, err := f.Apply(args)
		if err != nil {
			latestError = ErrArgsNotMatch{err}
			continue
		}

		ff := reflect.ValueOf(f.F)
		result := ff.Call(inArgs)

		if f.MayBeError {
			err := result[1].Elem()
			if !err.IsNil() {
				if greedy {
					continue
				}

				return nil, ErrFunctionFail{err.Interface().(error)}
			}
		}

		return result[0].Interface(), nil
	}

	return nil, latestError
}
