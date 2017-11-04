package task

import (
	"crypto/sha1"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"time"
)

type TaskID string

type Function interface{}

type Param interface{}

type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}

type Task struct {
	Schedule
	FuncName string
	Func     Function
	Params   []Param
}

func New(function Function, params ...Param) (*Task, error) {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		return nil, fmt.Errorf("Provided function value is not an actual function")
	}

	name := runtime.FuncForPC(funcValue.Pointer()).Name()
	return &Task{
		FuncName: name,
		Func:     function,
		Params:   params,
		Schedule: Schedule{
			IsRecurring: false,
		},
	}, nil
}

func (task *Task) IsDue() bool {
	timeNow := time.Now()
	return timeNow == task.NextRun || timeNow.After(task.NextRun)
}

func (task *Task) Run() {
	function := reflect.ValueOf(task.Func)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)

	task.scheduleNextRun()
}

func (task *Task) Hash() TaskID {
	hash := sha1.New()
	io.WriteString(hash, task.FuncName)
	io.WriteString(hash, fmt.Sprintf("%+v", task.Params))
	io.WriteString(hash, fmt.Sprintf("%s", task.Schedule.Duration))
	io.WriteString(hash, fmt.Sprintf("%t", task.Schedule.IsRecurring))
	return TaskID(hash.Sum(nil))
}

func (task *Task) scheduleNextRun() {
	if !task.IsRecurring {
		return
	}

	task.LastRun = task.NextRun
	task.NextRun = task.NextRun.Add(task.Duration)
}
