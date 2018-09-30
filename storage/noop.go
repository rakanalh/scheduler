package storage

// NoOpStorage is an ineffective storage which can be used to prevent storing tasks altogether.
type NoOpStorage struct {
}

// NewNoOpStorage returns an instance of NoOpStorage
func NewNoOpStorage() NoOpStorage {
	return NoOpStorage{}
}

// Add does nothing
func (noop NoOpStorage) Add(task TaskAttributes) error {
	return nil
}

// Fetch returns an empty list of tasks
func (noop NoOpStorage) Fetch() ([]TaskAttributes, error) {
	return []TaskAttributes{}, nil
}

// Remove does nothing
func (noop NoOpStorage) Remove(task TaskAttributes) error {
	return nil
}

func (noop NoOpStorage) Close() error {
	return nil
}
