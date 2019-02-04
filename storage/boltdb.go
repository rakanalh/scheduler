package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

const TaskStoreBucket string = "task_store"

// BoltDBConfig is the config structure holding information about boltdb db.
type BoltDBConfig struct {
	DBPath string
}

// BoltDBStorage is the structure responsible for handling boltdb storage.
type BoltDBStorage struct {
	config BoltDBConfig
	db     *bolt.DB
}

// NewBoltDBStorage returns a new instance of BoltDBStorage.
func NewBoltDBStorage(config BoltDBConfig) *BoltDBStorage {
	return &BoltDBStorage{
		config: config,
	}
}

// Connect opens the database file, or creates it if it does not exist.
func (b *BoltDBStorage) Connect() error {
	db, err := bolt.Open(b.config.DBPath, 0600, &bolt.Options{Timeout: time.Second * 1})
	if err != nil {
		log.Printf("could not connnect to database")
		return fmt.Errorf("boldb error: %s ", err)
	}
	b.db = db
	return nil
}

// initialize boltdb task_store bucket
func (b *BoltDBStorage) initialize() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(TaskStoreBucket))
		if err != nil {
			log.Printf("could not initialize bucket")
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return err
}

// Add stores the task to boltdb.
func (b *BoltDBStorage) Add(task TaskAttributes) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TaskStoreBucket))
		var taskBuf bytes.Buffer
		var id bytes.Buffer
		if err := gob.NewEncoder(&id).Encode(task.Hash); err != nil {
			return fmt.Errorf("task hash encoding : %s", err)
		}
		if err := gob.NewEncoder(&taskBuf).Encode(task); err != nil {
			return fmt.Errorf("task encoding : %s", err)
		}
		err := b.Put(id.Bytes(), taskBuf.Bytes())
		if err != nil {
			return fmt.Errorf("db update : %s", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Fetch will return the list of all stored tasks.
func (b *BoltDBStorage) Fetch() ([]TaskAttributes, error) {
	var tasks []TaskAttributes
	getTask := func(b []byte) (TaskAttributes, error) {
		var task TaskAttributes
		d := bytes.NewBuffer(b)
		if err := gob.NewDecoder(d).Decode(&task); err != nil {
			return TaskAttributes{}, fmt.Errorf("task decoding error : %s", err)
		}
		return task, nil
	}
	if err := b.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(TaskStoreBucket))
		if err := buck.ForEach(func(k, v []byte) error {
			task, err := getTask(v)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return []TaskAttributes{}, err
	}
	return tasks, nil
}

// Remove will delete the task from boltdb storage.
func (b *BoltDBStorage) Remove(task TaskAttributes) error {
	if err := b.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(TaskStoreBucket))
		if err := buck.Delete([]byte(task.Hash)); err != nil {
			return fmt.Errorf("task removal err: %s", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Close will close the open DB file.
func (b *BoltDBStorage) Close() error {
	return b.db.Close()
}
