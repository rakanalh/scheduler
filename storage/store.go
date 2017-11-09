package storage

type TaskAttributes map[string]string

type TaskStore interface {
	Add(TaskAttributes) error
	Fetch() ([]TaskAttributes, error)
	Remove(TaskAttributes) error
}
