package main

import (
	"log"
	"time"

	"github.com/rakanalh/scheduler"
	"github.com/rakanalh/scheduler/storage"
)

func TaskWithoutArgs() {
	log.Println("TaskWithoutArgs is executed")
}

func TaskWithArgs(message string) {
	log.Println("TaskWithArgs is executed. message:", message)
}

func main() {
	storage := storage.NewSqlite3Storage(
		storage.Sqlite3Config{
			DbName: "task_store.db",
		},
	)
	if err := storage.Connect(); err != nil {
		log.Fatal("Could not connect to db", err)
	}

	if err := storage.Initialize(); err != nil {
		log.Fatal("Could not intialize database", err)
	}

	s := scheduler.New(storage)

	// Start a task without arguments
	if _, err := s.RunAfter(30*time.Second, TaskWithoutArgs); err != nil {
		log.Fatal(err)
	}

	// Start a task with arguments
	if _, err := s.RunEvery(5*time.Second, TaskWithArgs, "Hello from recurring task 1"); err != nil {
		log.Fatal(err)
	}

	// Start the same task as above with a different argument
	if _, err := s.RunEvery(10*time.Second, TaskWithArgs, "Hello from recurring task 2"); err != nil {
		log.Fatal(err)
	}
	s.Start()
	s.Wait()
}
