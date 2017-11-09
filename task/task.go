package task

import (
	"crypto/sha1"
	"fmt"
	"io"
	"reflect"
	"time"
)

type TaskID string

type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}

type Task struct {
	Schedule
	Func   FunctionMeta
	Params []Param
}

func New(function FunctionMeta, params []Param) *Task {
	return &Task{
		Func:   function,
		Params: params,
	}
}

func NewWithSchedule(function FunctionMeta, params []Param, schedule Schedule) *Task {
	return &Task{
		Func:     function,
		Params:   params,
		Schedule: schedule,
	}
}

func (task *Task) IsDue() bool {
	timeNow := time.Now()
	return timeNow == task.NextRun || timeNow.After(task.NextRun)
}

func (task *Task) Run() {
	function := reflect.ValueOf(task.Func.function)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)

	task.scheduleNextRun()
}

func (task *Task) Hash() TaskID {
	hash := sha1.New()
	io.WriteString(hash, task.Func.Name)
	io.WriteString(hash, fmt.Sprintf("%+v", task.Params))
	io.WriteString(hash, fmt.Sprintf("%s", task.Schedule.Duration))
	io.WriteString(hash, fmt.Sprintf("%t", task.Schedule.IsRecurring))
	return TaskID(fmt.Sprintf("%x", hash.Sum(nil)))
}

func (task *Task) scheduleNextRun() {
	if !task.IsRecurring {
		return
	}

	task.LastRun = task.NextRun
	task.NextRun = task.NextRun.Add(task.Duration)
}
