package storage

import (
	"github.com/rakanalh/scheduler/task"
)

type NoOpStorage struct {
}

func NewNoOpStorage() NoOpStorage {
	return NoOpStorage{}
}

func (noop NoOpStorage) Store(task *task.Task) error {
	return nil
}

func (noop NoOpStorage) Fetch() ([]*task.Task, error) {
	return nil, nil
}
