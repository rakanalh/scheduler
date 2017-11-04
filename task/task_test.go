package task

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	mock := CallbackMock{}
	_, err := New(mock.CallNoArgs)
	if err != nil {
		t.Error("New returned an error when it should succeed", err)
	}

	fakeCallback := "String"
	_, err = New(fakeCallback)
	if err == nil {
		t.Error("New did not fail when passing a non-function value")
	}
}

func TestTaskIsDue(t *testing.T) {
	mock := CallbackMock{}
	task, _ := New(mock.CallNoArgs)
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

	task, _ := New(mock.CallNoArgs)
	task.Run()

	mock.AssertExpectations(t)
}

func TestTaskRunWithArgs(t *testing.T) {
	mock := CallbackMock{}
	mock.On("CallWithArgs", "Test", true).Return()

	task, _ := New(mock.CallWithArgs, "Test", true)
	task.Run()

	mock.AssertExpectations(t)
}

func TestTaskRunScheduledNextRun(t *testing.T) {
	mock := CallbackMock{}
	mock.On("CallNoArgs").Return()

	timeNow := time.Now()
	task, _ := New(mock.CallNoArgs)
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
	task, _ := New(mock.CallNoArgs)
	task.IsRecurring = true
	task.NextRun = time.Now()
	task.Duration = 5 * time.Second
	hash := task.Hash()

	if hash == "" {
		t.Fail()
	}
}
