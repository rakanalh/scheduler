package storage

type TaskAttributes struct {
	Hash        string
	Name        string
	LastRun     string
	NextRun     string
	Duration    string
	IsRecurring string
	Params      string
}

type TaskStore interface {
	Add(TaskAttributes) error
	Fetch() ([]TaskAttributes, error)
	Remove(TaskAttributes) error
}
