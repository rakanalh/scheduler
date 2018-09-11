package main

import (
"github.com/rakanalh/scheduler"
"github.com/rakanalh/scheduler/storage"
"io"
"log"
"time"
)

func TaskWithoutArgs() {
	log.Println("TaskWithoutArgs is executed")
}

func TaskWithArgs(message string) {
	log.Println("TaskWithArgs is executed. message:", message)
}

func main() {
	storage,err := storage.NewPostgresStorage(
		storage.PostgresConfig{
			DbURL: "postgresql://<db-username>:<db-password>@localhost:5432/scheduler?sslmode=disable",
		},
	)
	if err != nil{
		log.Fatalf("Couldn't create scheduler storage : %v",err)
	}

	s := scheduler.New(storage)

	go func(s scheduler.Scheduler,store io.Closer){
		time.Sleep(time.Second *10)
		// store.Close()
		s.Stop()
	}(s,storage)
	// Start a task without arguments
	if _, err := s.RunAfter(60*time.Second, TaskWithoutArgs); err != nil {
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
