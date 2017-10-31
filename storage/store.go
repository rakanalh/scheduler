package storage

import "github.com/rakanalh/scheduler/task"

type TaskStore interface {
	Store(task *task.ScheduledTask) error
	Fetch() ([]*task.ScheduledTask, error)
}
