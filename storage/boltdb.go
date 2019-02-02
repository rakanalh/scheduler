package storage

// BoltDBConfig is the config structure holding information about boltdb db.
type BoltDBConfig struct {
	DBPath string
}

// BoltDBStorage is the structure responsible for handling boltdb storage.
type BoltDBStorage struct {
	Config BoltDBConfig
}

// NewBoltDBStorage returns a new instance of BoltDBStorage.
func NewBoltDBStorage(config BoltDBConfig) *BoltDBStorage {
	return &BoltDBStorage{
		Config: config,
	}
}

// Connect opens the database file, or creates it if it does not exist.
func (b *BoltDBStorage) Connect() {

}

// Add stores the task to boltdb.
func (b *BoltDBStorage) Add(TaskAttributes) error {
	panic("implement me")
}

// Fetch will return the list of all stored tasks.
func (b *BoltDBStorage) Fetch() ([]TaskAttributes, error) {
	panic("implement me")
}

// Remove will delete the task from boltdb storage.
func (b *BoltDBStorage) Remove(TaskAttributes) error {
	panic("implement me")
}

// Close will close the open DB file.
func (b *BoltDBStorage) Close() error {
	panic("implement me")
}
