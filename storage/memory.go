package storage

type MemoryStorage struct {
	tasks []TaskAttributes
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (memStore *MemoryStorage) Add(task TaskAttributes) error {
	memStore.tasks = append(memStore.tasks, task)
	return nil
}

func (memStore *MemoryStorage) Fetch() ([]TaskAttributes, error) {
	return memStore.tasks, nil
}

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
