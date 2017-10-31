package scheduler

import (
	"time"

	"github.com/rakanalh/scheduler/storage"
	"github.com/rakanalh/scheduler/task"
)

type Scheduler struct {
	tasks     map[string]*task.ScheduledTask
	taskStore storage.TaskStore
	stopChan  chan bool
}

func New(store storage.TaskStore) Scheduler {
	return Scheduler{
		taskStore: store,
		stopChan:  make(chan bool),
		tasks:     make(map[string]*task.ScheduledTask),
	}
}

func (scheduler *Scheduler) RunAt(time time.Time, function task.Function, params ...task.Param) error {
	task, err := scheduler.makeTask(function, params...)
	if err != nil {
		return err
	}

	task.NextRun = time

	return nil
}

func (scheduler *Scheduler) RunAfter(duration time.Duration, function task.Function, params ...task.Param) error {
	return scheduler.RunAt(time.Now().Add(duration), function, params...)
}

func (scheduler *Scheduler) RunEvery(duration time.Duration, function task.Function, params ...task.Param) error {
	task, err := scheduler.makeTask(function, params...)
	if err != nil {
		return err
	}

	task.IsRecurring = true
	task.Duration = duration
	task.NextRun = time.Now().Add(duration)

	return nil
}

func (scheduler *Scheduler) Start() {
	// TODO: Implement signal handling

	// Populate tasks from storage

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				scheduler.runPending()
			case <-scheduler.stopChan:
				return
			}
		}
	}()
}

func (scheduler *Scheduler) Wait() {
	<-scheduler.stopChan
}

func (scheduler *Scheduler) populateTasks() error {
	tasks, err := scheduler.taskStore.Fetch()
	if err != nil {
		return err
	}

	for _, dbTask := range tasks {
		// If we can't find the task, it's been changed/removed by user
		registeredTask, ok := scheduler.tasks[dbTask.Name]
		if !ok {
			continue
		}
		// Duration may have changed
		if registeredTask.Duration != dbTask.Duration {
			// Reschedule NextRun based on dbTask.LastRun + registeredTask.Duration
			registeredTask.NextRun = dbTask.LastRun.Add(registeredTask.Duration)
		}
	}

	return nil
}

func (scheduler *Scheduler) runPending() {
	for _, task := range scheduler.tasks {
		if task.IsDue() {
			go task.Run()

			if !task.IsRecurring {
				delete(scheduler.tasks, task.Name)
			}
		}
	}
}

func (scheduler *Scheduler) makeTask(function task.Function, params ...task.Param) (*task.ScheduledTask, error) {
	task, err := task.NewTask(function, params...)
	if err != nil {
		return nil, err
	}

	scheduler.tasks[task.Name] = task
	return task, nil
}
