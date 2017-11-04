package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
)

func TaskWithoutArgs() {
	fmt.Println("TaskWithoutArgs is executed")
}

func TaskWithArgs(message string) {
	fmt.Println("TaskWithArgs is executed. message:", name)
}

func main() {
	storage := storage.NewSqlite3Storage(
		storage.Sqlite3Config{
			DbName: "task_store.db",
		},
		storage.NewMarshaler(),
	)
	if err := storage.Connect(); err != nil {
		log.Fatal("Could not connect to db", err)
	}

	if err := storage.Initialize(); err != nil {
		log.Fatal("Could not intialize database", err)
	}

	s := scheduler.New(storage)

	if err := s.RunAfter(5*time.Second, TaskWithoutArgs); err != nil {
		log.Fatal(err)
	}
	if err := s.RunEvery(5*time.Second, TaskWithArgs, "Hello from recurring task"); err != nil {
		log.Fatal(err)
	}
	s.Start()
	s.Wait()
}
