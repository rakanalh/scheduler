package task

import (
	"fmt"
	"reflect"
	"runtime"
)

type Function interface{}

type Param interface{}

type FunctionMeta struct {
	Name     string
	function Function
	params   map[string]reflect.Type
}

type FuncRegistry struct {
	funcs map[string]FunctionMeta
}

func NewFuncRegistry() *FuncRegistry {
	return &FuncRegistry{
		funcs: make(map[string]FunctionMeta),
	}
}

func (reg *FuncRegistry) Add(function Function) (FunctionMeta, error) {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		return FunctionMeta{}, fmt.Errorf("Provided function value is not an actual function")
	}

	name := runtime.FuncForPC(funcValue.Pointer()).Name()
	funcInstance, err := reg.Get(name)
	if err == nil {
		return funcInstance, nil
	}
	reg.funcs[name] = FunctionMeta{
		Name:     name,
		function: function,
		params:   reg.resolveParamTypes(function),
	}
	return reg.funcs[name], nil
}

func (reg *FuncRegistry) Get(name string) (FunctionMeta, error) {
	function, ok := reg.funcs[name]
	if ok {
		return function, nil
	}
	return FunctionMeta{}, fmt.Errorf("Function %s not found", name)
}

func (reg *FuncRegistry) Exists(name string) bool {
	_, ok := reg.funcs[name]
	if ok {
		return true
	}
	return false
}

func (reg *FuncRegistry) resolveParamTypes(function Function) map[string]reflect.Type {
	paramTypes := make(map[string]reflect.Type)
	funcType := reflect.TypeOf(function)
	for idx := 0; idx < funcType.NumIn(); idx++ {
		in := funcType.In(idx)
		paramTypes[in.Name()] = in
	}
	return paramTypes
}

func (meta *FunctionMeta) Params() []reflect.Type {
	funcType := reflect.TypeOf(meta.function)
	paramTypes := make([]reflect.Type, funcType.NumIn())
	for idx := 0; idx < funcType.NumIn(); idx++ {
		in := funcType.In(idx)
		paramTypes[idx] = in
	}
	return paramTypes
}
