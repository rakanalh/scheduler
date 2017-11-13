package scheduler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rakanalh/scheduler/storage"
	"github.com/rakanalh/scheduler/task"
)

type Scheduler struct {
	funcRegistry *task.FuncRegistry
	stopChan     chan bool
	tasks        map[task.TaskID]*task.Task
	taskStore    storeBridge
}

func New(store storage.TaskStore) Scheduler {
	funcRegistry := task.NewFuncRegistry()
	return Scheduler{
		funcRegistry: funcRegistry,
		stopChan:     make(chan bool),
		tasks:        make(map[task.TaskID]*task.Task),
		taskStore: storeBridge{
			store:        store,
			funcRegistry: funcRegistry,
		},
	}
}

func (scheduler *Scheduler) RunAt(time time.Time, function task.Function, params ...task.Param) (task.TaskID, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := task.New(funcMeta, params)

	task.NextRun = time

	scheduler.registerTask(task)
	return task.Hash(), nil
}

func (scheduler *Scheduler) RunAfter(duration time.Duration, function task.Function, params ...task.Param) (task.TaskID, error) {
	return scheduler.RunAt(time.Now().Add(duration), function, params...)
}

func (scheduler *Scheduler) RunEvery(duration time.Duration, function task.Function, params ...task.Param) (task.TaskID, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := task.New(funcMeta, params)

	task.IsRecurring = true
	task.Duration = duration
	task.NextRun = time.Now().Add(duration)

	scheduler.registerTask(task)
	return task.Hash(), nil
}

func (scheduler *Scheduler) Start() error {
	log.Println("Scheduler is starting...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Populate tasks from storage
	if err := scheduler.populateTasks(); err != nil {
		return nil
	}
	if err := scheduler.persistRegisteredTasks(); err != nil {
		return nil
	}
	scheduler.runPending()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				scheduler.runPending()
			case <-sigChan:
				scheduler.stopChan <- true
			case <-scheduler.stopChan:
				close(scheduler.stopChan)
			}
		}
	}()

	return nil
}

func (scheduler *Scheduler) Stop() {
	scheduler.stopChan <- true
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
		// If we can't find the function, it's been changed/removed by user
		exists := scheduler.funcRegistry.Exists(dbTask.Func.Name)
		if !exists {
			log.Printf("%s was not found, it will be removed\n", dbTask.Func.Name)
			_ = scheduler.taskStore.Remove(dbTask)
			continue
		}

		// If the task instance is still registered with the same computed hash then move on.
		// Otherwise, one of the attributes changed and therefore, the task instance should
		// be added to the list of tasks to be executed with the stored params
		registeredTask, ok := scheduler.tasks[dbTask.Hash()]
		if !ok {
			log.Printf("Detected a change in attributes of one of the instances of task %s, \n",
				dbTask.Func.Name)
			dbTask.Func, _ = scheduler.funcRegistry.Get(dbTask.Func.Name)
			registeredTask = dbTask
			scheduler.tasks[dbTask.Hash()] = registeredTask
		}

		// Skip task which is not a recurring one and the NextRun has already passed
		if !dbTask.IsRecurring && dbTask.NextRun.Before(time.Now()) {
			// We might have a task instance which was executed already.
			// In this case, delete it.
			_ = scheduler.taskStore.Remove(dbTask)
			delete(scheduler.tasks, dbTask.Hash())
			continue
		}

		// Duration may have changed for recurring tasks
		if dbTask.IsRecurring && registeredTask.Duration != dbTask.Duration {
			// Reschedule NextRun based on dbTask.LastRun + registeredTask.Duration
			registeredTask.NextRun = dbTask.LastRun.Add(registeredTask.Duration)
		}
	}
	return nil
}

func (scheduler *Scheduler) persistRegisteredTasks() error {
	for _, task := range scheduler.tasks {
		err := scheduler.taskStore.Add(task)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scheduler *Scheduler) runPending() {
	for _, task := range scheduler.tasks {
		if task.IsDue() {
			go task.Run()

			if !task.IsRecurring {
				_ = scheduler.taskStore.Remove(task)
				delete(scheduler.tasks, task.Hash())
			}
		}
	}
}

func (scheduler *Scheduler) registerTask(task *task.Task) {
	_, _ = scheduler.funcRegistry.Add(task.Func)
	scheduler.tasks[task.Hash()] = task
}
