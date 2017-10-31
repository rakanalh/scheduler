package scheduler

import (
	"testing"
	"time"

	"github.com/rakanalh/scheduler/storage"
	"github.com/rakanalh/scheduler/task"
)

const TestTaskName = "github.com/rakanalh/scheduler/task.(*CallbackMock).CallNoArgs-fm"

func TestRunAt(t *testing.T) {
	mock := task.CallbackMock{}

	timeNow := time.Now()
	scheduler := New(storage.NoOpStorage{})
	err := scheduler.RunAt(timeNow, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	err = scheduler.RunAt(timeNow, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if len(scheduler.tasks) > 1 {
		t.Error("There should only be one task")
	}

	if scheduler.tasks[TestTaskName].NextRun != timeNow {
		t.Error("The task's NextRun should be equal to passed parameter")
	}
}

func TestRunAfter(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NoOpStorage{})
	err := scheduler.RunAfter(5, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}
	err = scheduler.RunAfter(5, "InvalidFunction")
	if err == nil {
		t.Error()
	}
}

func TestRunEvery(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NoOpStorage{})
	err := scheduler.RunEvery(5, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	err = scheduler.RunEvery(5, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if !scheduler.tasks[TestTaskName].IsRecurring {
		t.Error()
	}
}

func TestRunPending(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NoOpStorage{})
	err := scheduler.RunAt(time.Now(), mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	mock.On("CallNoArgs").Return()

	scheduler.runPending()

	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectations(t)

	if len(scheduler.tasks) > 0 {
		t.Error("Non-recurring task should be removed once executed")
	}

	// Test again with a recurring task
	_ = scheduler.RunEvery(5, mock.CallNoArgs)

	mock.On("CallNoArgs").Return()

	// Task should be executed and then rescheduled
	scheduler.runPending()
	time.Sleep(100 * time.Millisecond)
	mock.AssertExpectations(t)
	if len(scheduler.tasks) == 0 {
		t.Error("The recurring task should still exist")
	}
}

func TestStart(t *testing.T) {
	mock := task.CallbackMock{}
	mock.On("CallNoArgs").Return()

	scheduler := New(storage.NoOpStorage{})
	err := scheduler.RunAt(time.Now(), mock.CallNoArgs)
	if err != nil {
		t.Error("Should not fail")
	}
	scheduler.Start()

	time.AfterFunc(2*time.Second, func() {
		close(scheduler.stopChan)
	})
	scheduler.Wait()
	mock.On("CallNoArgs").Return()

	// Task should be executed and then rescheduled
	mock.AssertExpectations(t)
}
