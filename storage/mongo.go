// +build cgo

package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

const COLLECTION_NAME string = "task_store"

// MongoDBConfig is the config structure holding information about mongo db.
type MongoDBConfig struct {
	ConnectionUrl string
	Db            string
}

// MongoDBStorage is the structure responsible for handling mongo storage.
type MongoDBStorage struct {
	config MongoDBConfig
	client *mongo.Client
}

// NewMongoDBStorage returns a new instance of MongoDBStorage.
func NewMongoDBStorage(config MongoDBConfig) *MongoDBStorage {
	return &MongoDBStorage{
		config: config,
	}
}

// Connect creates the database file.
func (mongodb *MongoDBStorage) Connect() error {
	var client *mongo.Client

	client, err := mongo.NewClient(mongodb.config.ConnectionUrl)
	if err != nil {
		return err
	}

	mongodb.client = client
	err = mongodb.client.Connect(context.TODO())

	return err
}

// Close will close the open DB file.
func (mongodb *MongoDBStorage) Close() error {
	return mongodb.client.Disconnect(context.Background())
}

// Initialize mongodb collection
func (mongodb *MongoDBStorage) Initialize() error {
	task_store := mongodb.client.
		Database(mongodb.config.Db).Collection(COLLECTION_NAME)

	if task_store == nil {
		log.Printf("could not initialize collection")
		return errors.New("mongo error")
	}

	return nil
}

// Stores the task to mongo
func (mongodb MongoDBStorage) Add(task TaskAttributes) error {
	task_store := mongodb.client.Database(mongodb.config.Db).Collection(COLLECTION_NAME)

	if task_store == nil {
		return errors.New("could not get collection")
	}

	// filter := bson.NewDocument(bson.EC.String("hash", task.Hash))
	filter := bsonx.Doc{{"hash", bsonx.String(task.Hash)}}
	res, err := task_store.Count(context.Background(), filter)

	if err != nil {
		return errors.New(fmt.Sprintf("%v", err))
	}

	if res == 0 {
		res, err := task_store.InsertOne(context.Background(),
			map[string]string{
				"name":         task.Name,
				"params":       task.Params,
				"duration":     task.Duration,
				"last_run":     task.LastRun,
				"next_run":     task.NextRun,
				"is_recurring": task.IsRecurring,
				"hash":         task.Hash,
			})
		if res == nil {
			return errors.New("element not inserted")
		}
		return err
	}
	return nil
}

// Remove will delete the task from storage.
func (mongodb MongoDBStorage) Remove(task TaskAttributes) error {
	task_store := mongodb.client.Database(mongodb.config.Db).Collection(COLLECTION_NAME)

	if task_store == nil {
		return errors.New("could not get collection")
	}

	// filter := bson.NewDocument(bson.EC.String("hash", task.Hash))
	filter := bsonx.Doc{{"hash", bsonx.String(task.Hash)}}

	_, err := task_store.DeleteOne(context.Background(), filter)

	return err
}

// Fetch will return the list of all stored tasks.
func (mongodb MongoDBStorage) Fetch() ([]TaskAttributes, error) {
	task_store := mongodb.client.Database(mongodb.config.Db).Collection(COLLECTION_NAME)

	if task_store == nil {
		log.Println("taskstore empty")
		return nil, errors.New("could not get collection")
	}

	var tasks []TaskAttributes

	cur, err := task_store.Find(context.Background(), nil)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var elem bsonx.Doc
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		task := TaskAttributes{
			Name:        elem.Lookup("name").StringValue(),
			Params:      elem.Lookup("params").StringValue(),
			LastRun:     elem.Lookup("last_run").StringValue(),
			NextRun:     elem.Lookup("next_run").StringValue(),
			Duration:    elem.Lookup("duration").StringValue(),
			IsRecurring: elem.Lookup("is_recurring").StringValue(),
			Hash:        elem.Lookup("hash").StringValue(),
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}
