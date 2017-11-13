package task

import (
	"testing"
	"time"
)

func TestTaskIsDue(t *testing.T) {
	mock := CallbackMock{}
	task := newTestTaskWithSchedule(t, mock.CallNoArgs, []Param{}, Schedule{
		IsRecurring: false,
		LastRun:     time.Now(),
		NextRun:     time.Now(),
		Duration:    0,
	})
	task.NextRun = time.Now()
	if !task.IsDue() {
		t.Error("Task should be due")
	}
	task.NextRun = time.Now().Add(-5 * time.Second)
	if !task.IsDue() {
		t.Error("Task (now - 5 seconds) should be due")
	}

	task.NextRun = time.Now().Add(5 * time.Second)
	if task.IsDue() {
		t.Error("Task (now + 5 seconds) should not be due")
	}
}

func TestTaskRun(t *testing.T) {
	mock := CallbackMock{}
	mock.On("CallNoArgs").Return()

	task := newTestTask(t, mock.CallNoArgs, []Param{})
	task.Run()

	mock.AssertExpectations(t)
}

func TestTaskRunWithArgs(t *testing.T) {
	mock := CallbackMock{}
	mock.On("CallWithArgs", "Test", true).Return()

	task := newTestTask(t, mock.CallWithArgs, []Param{"Test", true})
	task.Run()

	mock.AssertExpectations(t)
}

func TestTaskRunScheduledNextRun(t *testing.T) {
	mock := CallbackMock{}
	mock.On("CallNoArgs").Return()

	timeNow := time.Now()
	task := newTestTask(t, mock.CallNoArgs, []Param{})
	task.IsRecurring = true
	task.NextRun = timeNow
	task.Duration = 5 * time.Second
	task.Run()

	if task.NextRun != timeNow.Add(5*time.Second) {
		t.Fail()
	}

	mock.AssertExpectations(t)
}

func TestGenerateHash(t *testing.T) {
	mock := CallbackMock{}
	task := newTestTask(t, mock.CallNoArgs, []Param{})
	task.IsRecurring = true
	task.NextRun = time.Now()
	task.Duration = 5 * time.Second
	hash := task.Hash()

	if hash == "" {
		t.Fail()
	}
}

func newTestTask(t *testing.T, function Function, params []Param) *Task {
	funcMeta, err := newFuncMeta(function)
	if err != nil {
		t.Error("Failed to register function")
	}
	return New(funcMeta, params)
}

func newTestTaskWithSchedule(t *testing.T, function Function, params []Param, schedule Schedule) *Task {
	funcMeta, err := newFuncMeta(function)
	if err != nil {
		t.Error("Failed to register function")
	}
	return NewWithSchedule(funcMeta, params, schedule)
}
