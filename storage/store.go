package storage

import "github.com/rakanalh/scheduler/task"

type TaskStore interface {
	Store(task *task.Task) error
	Fetch() ([]*task.Task, error)
}
