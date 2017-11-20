package task

import (
	"crypto/sha1"
	"fmt"
	"io"
	"reflect"
	"time"
)

// ID is returned upon scheduling a task to be executed
type ID string

// Schedule holds information about the execution times of a specific task
type Schedule struct {
	IsRecurring bool
	LastRun     time.Time
	NextRun     time.Time
	Duration    time.Duration
}

// Task holds information about task
type Task struct {
	Schedule
	Func   FunctionMeta
	Params []Param
}

// New returns an instance of task
func New(function FunctionMeta, params []Param) *Task {
	return &Task{
		Func:   function,
		Params: params,
	}
}

// NewWithSchedule creates an instance of task with the provided schedule information
func NewWithSchedule(function FunctionMeta, params []Param, schedule Schedule) *Task {
	return &Task{
		Func:     function,
		Params:   params,
		Schedule: schedule,
	}
}

// IsDue returns a boolean indicating whether the task should execute or not
func (task *Task) IsDue() bool {
	timeNow := time.Now()
	return timeNow == task.NextRun || timeNow.After(task.NextRun)
}

// Run will execute the task and schedule it's next run.
func (task *Task) Run() {
	// Reschedule task first to prevent running the task
	// again in case the execution time takes more than the
	// task's duration value.
	task.scheduleNextRun()

	function := reflect.ValueOf(task.Func.function)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)
}

// Hash will return the SHA1 representation of the task's data.
func (task *Task) Hash() ID {
	hash := sha1.New()
	_, _ = io.WriteString(hash, task.Func.Name)
	_, _ = io.WriteString(hash, fmt.Sprintf("%+v", task.Params))
	_, _ = io.WriteString(hash, fmt.Sprintf("%s", task.Schedule.Duration))
	_, _ = io.WriteString(hash, fmt.Sprintf("%t", task.Schedule.IsRecurring))
	return ID(fmt.Sprintf("%x", hash.Sum(nil)))
}

func (task *Task) scheduleNextRun() {
	if !task.IsRecurring {
		return
	}

	task.LastRun = task.NextRun
	task.NextRun = task.NextRun.Add(task.Duration)
}
