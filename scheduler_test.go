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
	scheduler := New(storage.NewMemoryStorage())
	taskID, err := scheduler.RunAt(timeNow, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	_, err = scheduler.RunAt(timeNow, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if len(scheduler.tasks) > 1 {
		t.Error("There should only be one task")
	}

	if scheduler.tasks[taskID].NextRun != timeNow {
		t.Error("The task's NextRun should be equal to passed parameter")
	}
}

func TestRunAfter(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NewMemoryStorage())
	_, err := scheduler.RunAfter(5, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}
	_, err = scheduler.RunAfter(5, "InvalidFunction")
	if err == nil {
		t.Error()
	}
}

func TestRunEvery(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NewMemoryStorage())
	taskID, err := scheduler.RunEvery(5, mock.CallNoArgs)
	if err != nil {
		t.Error("Creating a task should succeed")
	}

	_, err = scheduler.RunEvery(5, "InvalidFunction")
	if err == nil {
		t.Error("InvalidFunction should have failed RunAt")
	}

	if !scheduler.tasks[taskID].IsRecurring {
		t.Error()
	}
}

func TestRunPending(t *testing.T) {
	mock := task.CallbackMock{}
	scheduler := New(storage.NewMemoryStorage())
	_, err := scheduler.RunAt(time.Now(), mock.CallNoArgs)
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
	_, _ = scheduler.RunEvery(5, mock.CallNoArgs)

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

	scheduler := New(storage.NewMemoryStorage())
	_, err := scheduler.RunAt(time.Now(), mock.CallNoArgs)
	if err != nil {
		t.Error("Should not fail")
	}
	scheduler.Start()

	time.AfterFunc(2*time.Second, func() {
		scheduler.Stop()
	})
	scheduler.Wait()
	mock.On("CallNoArgs").Return()

	// Task should be executed and then rescheduled
	mock.AssertExpectations(t)
}

func TestCancelTask(t *testing.T) {
	scheduler := New(storage.NewNoOpStorage())
	mock := task.CallbackMock{}

	err := scheduler.Cancel(task.ID("123456"))
	if err == nil {
		t.Error("Should fail because task does not exist")
	}

	taskID, _ := scheduler.RunAfter(5*time.Second, mock.CallNoArgs)
	err = scheduler.Cancel(taskID)
	if err != nil {
		t.Error("Clearing an existent task should not fail")
	}
}

func TestClearTask(t *testing.T) {
	scheduler := New(storage.NewNoOpStorage())
	mock := task.CallbackMock{}

	scheduler.RunAfter(5*time.Second, mock.CallNoArgs)
	scheduler.RunAfter(5*time.Second, mock.CallWithArgs)

	scheduler.Clear()

	if len(scheduler.tasks) > 0 {
		t.Error("Clearing tasks didn't take effect.")
	}
}

func TestPopulateTasks(t *testing.T) {
	mock := task.CallbackMock{}

	taskAttributes := storage.TaskAttributes{
		Hash:        "TestHash",
		LastRun:     "2017-11-10T12:00:00Z",
		NextRun:     "2017-11-10T12:00:00Z",
		Duration:    "5s",
		IsRecurring: "0",
		Name:        "github.com/rakanalh/scheduler/task.(*CallbackMock).CallNoArgs-fm",
		Params:      "",
	}

	memStore := storage.NewMemoryStorage()
	memStore.Add(taskAttributes)
	scheduler := New(memStore)
	scheduler.RunAfter(5, mock.CallNoArgs)
	err := scheduler.populateTasks()
	if err != nil {
		t.Error("Failed to populate tasks: ", err)
	}
}
