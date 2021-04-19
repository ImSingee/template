package template

var DefaultFunctionSet FunctionSet
var DefaultContext Context

func init() {
	DefaultFunctionSet = NewFunctionSet()
	DefaultContext = NewContext()
	DefaultContext.AddFunctionSet(DefaultFunctionSet)
}

func With(functionSets ...FunctionSet) Context {
	return DefaultContext.With(functionSets...)
}

func AddFunctionSet(functionSet FunctionSet) {
	DefaultContext.AddFunctionSet(functionSet)
}

func Execute(source string) (s string, err error) {
	return DefaultContext.Execute(source)
}

func Apply(source interface{}) (err error) {
	return DefaultContext.Apply(source)
}
