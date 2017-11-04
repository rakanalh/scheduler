package storage

import (
	"github.com/rakanalh/scheduler/task"
)

type NoOpStorage struct {
}

func NewNoOpStorage() NoOpStorage {
	return NoOpStorage{}
}

func (noop NoOpStorage) Add(task *task.Task) error {
	return nil
}

func (noop NoOpStorage) Fetch() ([]*task.Task, error) {
	return []*task.Task{}, nil
}

func (noop NoOpStorage) Remove(task *task.Task) error {
	return nil
}
