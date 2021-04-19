package template

import (
	"fmt"
)

// 调用 apply 时传递非指针参数
var errInvalidType = fmt.Errorf("invalid Type")

// 注册函数时传递非函数参数
var errIsNotFunction = fmt.Errorf("is not function")

// 找不到函数
type ErrNoFunction struct {
}

func (e ErrNoFunction) Error() string {
	return "no valid function signature"
}

var errNoFunction = ErrNoFunction{}

// 函数找到但没有适配的参数
type ErrArgsNotMatch struct {
	Inner error
}

func (e ErrArgsNotMatch) Error() string {
	return fmt.Sprintf("function args not match %s", e.Inner)
}

// 函数调用返回异常
type ErrFunctionFail struct {
	Inner error
}

func (e ErrFunctionFail) Error() string {
	return fmt.Sprintf("function call error %s", e.Inner.Error())
}
