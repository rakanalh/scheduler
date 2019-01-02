package storage

import (
	"context"
	"sync"
	"testing"

	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/stretchr/testify/require"
)

type TMongoStorage struct {
	storage *MongoDBStorage
	init    sync.Once
	uninit  sync.Once
}

// default credentials for testing database
const user string = ""
const pass string = ""

var mongoStorage *TMongoStorage = &TMongoStorage{}
var mongoConfig MongoDBConfig = MongoDBConfig{
	ConnectionUrl: "mongodb://localhost/test",
	Db:            "test",
}

var sampleTask = TaskAttributes{
	Hash:        "A",
	Name:        "B",
	LastRun:     "2018-09-30T20:00:00+02:00",
	NextRun:     "2018-09-30T20:00:05+02:00",
	Duration:    "5s",
	IsRecurring: "0",
	Params:      "null",
}

// Initializes database connection and collection cleanup
func (s *TMongoStorage) Init(config MongoDBConfig, t *testing.T) {
	s.init.Do(func() {
		s.storage = NewMongoDBStorage(config)
		err := s.storage.Connect()
		require.NoError(t, err)
		_, err = s.storage.client.Database(config.Db).
			Collection(COLLECTION_NAME).
			DeleteMany(context.Background(), bsonx.Doc{})
		require.NoError(t, err)
		s.uninit = sync.Once{}
	})
}

func (s *TMongoStorage) Uninit(t *testing.T) {
	s.uninit.Do(func() {
		s.init = sync.Once{}
		require.NoError(t, s.storage.Close())
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

	errAdd := mongoStorage.storage.Add(fetchTask)
	require.NoError(t, errAdd)
	tasks, err := mongoStorage.storage.Fetch()
	require.NoError(t, err)

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

// Test closing
func TestClose(t *testing.T) {
	mongoStorage.Init(mongoConfig, t)
	mongoStorage.Uninit(t)
}
