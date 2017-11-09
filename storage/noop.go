package storage

type NoOpStorage struct {
}

func NewNoOpStorage() NoOpStorage {
	return NoOpStorage{}
}

func (noop NoOpStorage) Add(task TaskAttributes) error {
	return nil
}

func (noop NoOpStorage) Fetch() ([]TaskAttributes, error) {
	return []TaskAttributes{}, nil
}

func (noop NoOpStorage) Remove(task TaskAttributes) error {
	return nil
}
