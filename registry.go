package scheduler

import (
	"fmt"

	"github.com/rakanalh/scheduler/task"
)

type FuncRegistry struct {
	funcs map[string]task.Function
}

func newFuncRegistry() *FuncRegistry {
	return &FuncRegistry{
		funcs: make(map[string]task.Function),
	}
}

func (reg *FuncRegistry) Add(name string, function task.Function) error {
	if reg.Exists(name) {
		return fmt.Errorf("Function already exists: %v", function)
	}
	reg.funcs[name] = function
	return nil
}

func (reg *FuncRegistry) Get(name string) (task.Function, error) {
	function, ok := reg.funcs[name]
	if ok {
		return function, nil
	}
	return nil, fmt.Errorf("Function %s not found", name)
}

func (reg *FuncRegistry) Exists(name string) bool {
	_, ok := reg.funcs[name]
	if ok {
		return true
	}
	return false
}
