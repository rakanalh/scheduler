package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type TMongoStorage struct {
	storage MongoDBStorage
	init    sync.Once
}

// default credentials for testing database
const user string = ""
const pass string = ""

var mongoStorage *TMongoStorage = &TMongoStorage{}
var mongoConfig MongoDBConfig = MongoDBConfig{
	HostName: "127.0.0.1",
	Port:     27017,
	Db:       "test_tasks",
}

var sampleTask = TaskAttributes{
	Hash:        "A",
	Name:        "B",
	LastRun:     "",
	NextRun:     "",
	Duration:    "",
	IsRecurring: "",
	Params:      "",
}

// Initializes database connection and collection cleanup
func (s *TMongoStorage) Init(config MongoDBConfig, t *testing.T) {
	s.init.Do(func() {
		s.storage = NewMongoDBStorage(config)
		err := s.storage.Connect()
		require.NoError(t, err)
		err = s.storage.clean()
		require.NoError(t, err)
	})
}

// Tests insertion
func TestAdd(t *testing.T) {
	mongoStorage.Init(mongoConfig, t)
	err := mongoStorage.storage.Add(sampleTask)
	require.NoError(t, err)
}

// Tests fetching all elements
func TestFetch(t *testing.T) {
	mongoStorage.Init(mongoConfig, t)
	fetchTask := sampleTask
	fetchTask.Hash = "C"

	err := mongoStorage.storage.Add(fetchTask)
	require.NoError(t, err)
	tasks, err := mongoStorage.storage.Fetch()

	var found bool = false
	for _, v := range tasks {
		if v.Hash == "C" {
			found = true
		}
	}

	require.False(t, !found)
}

// Tests removing an element
func TestRemove(t *testing.T) {
	mongoStorage.Init(mongoConfig, t)
	fetchTask := sampleTask
	fetchTask.Hash = "D"

	err := mongoStorage.storage.Add(fetchTask)
	require.NoError(t, err)
	err = mongoStorage.storage.Remove(fetchTask)
	require.NoError(t, err)
}
