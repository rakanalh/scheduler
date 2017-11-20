package task

import (
	"fmt"
	"reflect"
	"runtime"
)

// Function is a pointer to the callback function
type Function interface{}

// Param represents a single function parameter
type Param interface{}

// FunctionMeta holds information about function such as name and parameters.
type FunctionMeta struct {
	Name     string
	function Function
	params   map[string]reflect.Type
}

// FuncRegistry holds the list of all registered task functions.
type FuncRegistry struct {
	funcs map[string]FunctionMeta
}

// NewFuncRegistry will return an instance of the FuncRegistry.
func NewFuncRegistry() *FuncRegistry {
	return &FuncRegistry{
		funcs: make(map[string]FunctionMeta),
	}
}

// Add appends the function to the registry after resolving specific information about this function.
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

// Get returns the FunctionMeta instance which holds all information about any single registered task function.
func (reg *FuncRegistry) Get(name string) (FunctionMeta, error) {
	function, ok := reg.funcs[name]
	if ok {
		return function, nil
	}
	return FunctionMeta{}, fmt.Errorf("Function %s not found", name)
}

// Exists checks if a function with provided name exists.
func (reg *FuncRegistry) Exists(name string) bool {
	_, ok := reg.funcs[name]
	if ok {
		return true
	}
	return false
}

// Params returns the list of parameter types
func (meta *FunctionMeta) Params() []reflect.Type {
	funcType := reflect.TypeOf(meta.function)
	paramTypes := make([]reflect.Type, funcType.NumIn())
	for idx := 0; idx < funcType.NumIn(); idx++ {
		in := funcType.In(idx)
		paramTypes[idx] = in
	}
	return paramTypes
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
