package storage

import (
	"github.com/rakanalh/scheduler/task"
)

type NoOpStorage struct {
}

func NewNoOpStorage() NoOpStorage {
	return NoOpStorage{}
}

func (noop NoOpStorage) Store(task *task.ScheduledTask) error {
	return nil
}

func (noop NoOpStorage) Fetch() ([]*task.ScheduledTask, error) {
	return nil, nil
}
