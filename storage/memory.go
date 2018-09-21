package storage

// MemoryStorage is a memory task store
type MemoryStorage struct {
	tasks []TaskAttributes
}

// NewMemoryStorage returns an instance of MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

// Add adds a task to the memory store.
func (memStore *MemoryStorage) Add(task TaskAttributes) error {
	memStore.tasks = append(memStore.tasks, task)
	return nil
}

// Fetch will return all tasks stored.
func (memStore *MemoryStorage) Fetch() ([]TaskAttributes, error) {
	return memStore.tasks, nil
}

// Remove will remove task from store
func (memStore *MemoryStorage) Remove(task TaskAttributes) error {
	var newTasks []TaskAttributes
	for _, existingTask := range memStore.tasks {
		if task.Hash == existingTask.Hash {
			continue
		}
		newTasks = append(newTasks, existingTask)
	}
	memStore.tasks = newTasks
	return nil
}

func (memStore *MemoryStorage) Close() error {
	return nil
}
