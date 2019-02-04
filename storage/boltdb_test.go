package storage

import (
	"errors"
	"reflect"
	"testing"
	"time"

	bolt "go.etcd.io/bbolt"
)

var config = BoltDBConfig{DBPath: "test.db"}
var db, _ = bolt.Open(config.DBPath, 0600, &bolt.Options{Timeout: time.Second * 1})
var testTask = TaskAttributes{
	Hash:        "A",
	Name:        "B",
	LastRun:     "2018-09-30T20:00:00+02:00",
	NextRun:     "2018-09-30T20:00:05+02:00",
	Duration:    "5s",
	IsRecurring: "0",
	Params:      "null",
}

func TestNewBoltDBStorage(t *testing.T) {
	type args struct {
		config BoltDBConfig
	}
	tests := []struct {
		name string
		args args
		want *BoltDBStorage
	}{
		{name: "Constructor should return instance with correct db path set", args: args{config: config}, want: NewBoltDBStorage(BoltDBConfig{config.DBPath})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBoltDBStorage(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBoltDBStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltDBStorage_Connect(t *testing.T) {
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Should set a valid *bolt.DB object on the struct", fields: fields{config: config, db: db}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			if err := b.Connect(); (err != nil) != tt.wantErr && b.db.Path() != config.DBPath {
				t.Errorf("BoltDBStorage.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltDBStorage_initialize(t *testing.T) {
	getBucket := func(store *BoltDBStorage) string {
		if err := store.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(TaskStoreBucket))
			if b == nil {
				return errors.New("bucket does not exist")
			}
			return nil
		}); err != nil {
			return ""
		}
		return TaskStoreBucket
	}
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Task_Store bucket should be properly created", fields: fields{config: config, db: db}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			if err := b.initialize(); (err != nil) != tt.wantErr && getBucket(b) != TaskStoreBucket {
				t.Errorf("BoltDBStorage.initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltDBStorage_Add(t *testing.T) {
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	type args struct {
		task TaskAttributes
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Should insert the task without error", fields: fields{config: config, db: db}, args: args{task: testTask}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			if err := b.Add(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("BoltDBStorage.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltDBStorage_Fetch(t *testing.T) {
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []TaskAttributes
		wantErr bool
	}{
		{name: "Should return all tasks", fields: fields{config: config, db: db}, want: []TaskAttributes{testTask}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			got, err := b.Fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("BoltDBStorage.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BoltDBStorage.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltDBStorage_Remove(t *testing.T) {
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	type args struct {
		task TaskAttributes
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Task should not exist after removal", fields: fields{config: config, db: db}, args: args{testTask}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			task, _ := b.Fetch()
			if err := b.Remove(tt.args.task); (err != nil) != tt.wantErr && len(task) != 0 {
				t.Errorf("BoltDBStorage.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltDBStorage_Close(t *testing.T) {
	type fields struct {
		config BoltDBConfig
		db     *bolt.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Should return without error", fields: fields{config: config, db: db}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BoltDBStorage{
				config: tt.fields.config,
				db:     tt.fields.db,
			}
			if err := b.Close(); (err != nil) != tt.wantErr {
				t.Errorf("BoltDBStorage.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
