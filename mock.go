package scheduler

import (
	"fmt"

	"github.com/rakanalh/scheduler/storage"
)

type failureMode int

const (
	fail failureMode = iota
	failOnLastRun
	failOnNextRun
	failOnDuration
	failOnIsRecurring
	failOnFuncMeta
	failOnEmptyParams
	failOnEmptyListParams
)

type storeMock struct {
	Mode failureMode
}

func newStoreMockWithMode(mode failureMode) *storeMock {
	return &storeMock{
		Mode: mode,
	}
}

func (s *storeMock) Add(task storage.TaskAttributes) error {
	return nil
}

func (s *storeMock) Fetch() ([]storage.TaskAttributes, error) {
	if s.Mode == fail {
		return []storage.TaskAttributes{}, fmt.Errorf("Error")
	}

	taskAttributes := storage.TaskAttributes{
		Hash: "TestHash",
	}

	if s.Mode == failOnLastRun {
		taskAttributes.LastRun = "SomeCorruptString"
	} else {
		taskAttributes.LastRun = "2017-11-10T12:00:00Z"
	}

	if s.Mode == failOnNextRun {
		taskAttributes.NextRun = "SomeCorruptString"
	} else {
		taskAttributes.NextRun = "2017-11-10T12:00:00Z"
	}

	if s.Mode == failOnDuration {
		taskAttributes.Duration = "SomeCorruptString"
	} else {
		taskAttributes.Duration = "5s"
	}

	if s.Mode == failOnIsRecurring {
		taskAttributes.IsRecurring = "SomeCorruptString"
	} else {
		taskAttributes.IsRecurring = "1"
	}

	if s.Mode == failOnFuncMeta {
		taskAttributes.Name = "NonExistentName"
	} else {
		taskAttributes.Name = "github.com/rakanalh/scheduler.mockFunction"
	}

	if s.Mode == failOnEmptyParams {
		taskAttributes.Params = ""
	} else {
		taskAttributes.Params = "[]"
	}

	if s.Mode == failOnEmptyListParams {
		taskAttributes.Params = "SomeCorruptString"
	} else if s.Mode != failOnEmptyParams {
		taskAttributes.Params = "[]"
	}

	return []storage.TaskAttributes{taskAttributes}, nil
}

func (s *storeMock) Remove(task storage.TaskAttributes) error {
	return nil
}

func mockFunction(a string, b int) {

}
