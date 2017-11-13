package task

import (
	"reflect"
	"testing"
)

func TestRegistryFunc(t *testing.T) {
	mock := CallbackMock{}
	_, err := newFuncMeta(mock.CallNoArgs)
	if err != nil {
		t.Error("Failed to register function")
	}
}

func TestRegistryFake(t *testing.T) {
	fakeCallback := "String"
	_, err := newFuncMeta(fakeCallback)

	if err == nil {
		t.Error("New did not fail when passing a non-function value")
	}
}

func TestGet(t *testing.T) {
	mock := CallbackMock{}

	funcRegistry := NewFuncRegistry()
	funcMeta, err := funcRegistry.Add(mock.CallNoArgs)
	_, err = funcRegistry.Add(mock.CallWithArgs)

	if err != nil {
		t.Error("Failed to register function")
	}

	getResult, err := funcRegistry.Get("github.com/rakanalh/scheduler/task.(*CallbackMock).CallNoArgs-fm")

	if err != nil || funcMeta.Name != getResult.Name {
		t.Error("Could not find registered function")
	}
}

func TestAddExistingFunc(t *testing.T) {
	mock := CallbackMock{}

	funcRegistry := NewFuncRegistry()
	_, err := funcRegistry.Add(mock.CallWithArgs)
	if err != nil {
		t.Error("Failed to add function")
	}
	_, err = funcRegistry.Add(mock.CallWithArgs)
	if err != nil {
		t.Error("Failed to add existing function")
	}
}

func TestExists(t *testing.T) {
	mock := CallbackMock{}

	funcRegistry := NewFuncRegistry()
	_, err := funcRegistry.Add(mock.CallNoArgs)
	_, err = funcRegistry.Add(mock.CallWithArgs)

	if err != nil {
		t.Error("Failed to register function")
	}

	found := funcRegistry.Exists("CallNoArgs")

	if found {
		t.Error("Found a non-registered function")
	}
	found = funcRegistry.Exists("github.com/rakanalh/scheduler/task.(*CallbackMock).CallNoArgs-fm")
	if !found {
		t.Error("Couldn't find a registered function")
	}
}

func TestFunctionMetaParams(t *testing.T) {
	mock := CallbackMock{}
	funcMeta, _ := newFuncMeta(mock.CallWithArgs)
	params := funcMeta.Params()

	expectedParams := []reflect.Type{
		reflect.TypeOf(""),
		reflect.TypeOf(true),
	}
	for idx, param := range params {
		if expectedParams[idx].Name() != param.Name() {
			t.Error("Types don't match")
		}
	}
}

func newFuncMeta(function Function) (FunctionMeta, error) {
	funcRegistry := NewFuncRegistry()
	return funcRegistry.Add(function)
}
