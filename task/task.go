package task

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type Function interface{}

type Param interface{}

type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}

type ScheduledTask struct {
	Schedule
	Name   string
	Func   Function
	Params []Param
}

func NewTask(function Function, params ...Param) (*ScheduledTask, error) {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("Provided function value is not an actual function")
	}

	name := runtime.FuncForPC(funcValue.Pointer()).Name()
	return &ScheduledTask{
		Name:   name,
		Func:   function,
		Params: params,
		Schedule: Schedule{
			IsRecurring: false,
		},
	}, nil
}

func (task *ScheduledTask) IsDue() bool {
	timeNow := time.Now()
	return timeNow == task.NextRun || timeNow.After(task.NextRun)
}

func (task *ScheduledTask) Run() {
	function := reflect.ValueOf(task.Func)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)

	task.scheduleNextRun()
}

func (task *ScheduledTask) scheduleNextRun() {
	if !task.IsRecurring {
		return
	}

	task.LastRun = task.NextRun
	task.NextRun = task.NextRun.Add(task.Duration)
}
