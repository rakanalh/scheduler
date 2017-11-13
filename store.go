package scheduler

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/rakanalh/scheduler/storage"
	"github.com/rakanalh/scheduler/task"
)

type storeBridge struct {
	store        storage.TaskStore
	funcRegistry *task.FuncRegistry
}

func (sb *storeBridge) Add(task *task.Task) error {
	attributes, err := sb.getTaskAttributes(task)
	if err != nil {
		return err
	}
	return sb.store.Add(attributes)
}

func (sb *storeBridge) Fetch() ([]*task.Task, error) {
	storedTasks, err := sb.store.Fetch()
	if err != nil {
		return []*task.Task{}, err
	}
	var tasks []*task.Task
	for _, storedTask := range storedTasks {
		lastRun, err := time.Parse(time.RFC3339, storedTask.LastRun)
		if err != nil {
			return nil, err
		}

		nextRun, err := time.Parse(time.RFC3339, storedTask.NextRun)
		if err != nil {
			return nil, err
		}

		duration, err := time.ParseDuration(storedTask.Duration)
		if err != nil {
			return nil, err
		}

		isRecurring, err := strconv.Atoi(storedTask.IsRecurring)
		if err != nil {
			return nil, err
		}

		funcMeta, err := sb.funcRegistry.Get(storedTask.Name)
		if err != nil {
			return nil, err
		}

		params, err := paramsFromString(funcMeta, storedTask.Params)
		if err != nil {
			return nil, err
		}

		t := task.NewWithSchedule(funcMeta, params, task.Schedule{
			IsRecurring: isRecurring == 1,
			Duration:    time.Duration(duration),
			LastRun:     lastRun,
			NextRun:     nextRun,
		})
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (sb *storeBridge) Remove(task *task.Task) error {
	attributes, err := sb.getTaskAttributes(task)
	if err != nil {
		return err
	}
	return sb.store.Remove(attributes)
}

func (sb *storeBridge) getTaskAttributes(task *task.Task) (storage.TaskAttributes, error) {
	params, err := paramsToString(task.Params)
	if err != nil {
		return storage.TaskAttributes{}, err
	}

	isRecurring := 0
	if task.IsRecurring {
		isRecurring = 1
	}

	return storage.TaskAttributes{
		Hash:        string(task.Hash()),
		Name:        task.Func.Name,
		LastRun:     task.LastRun.Format(time.RFC3339),
		NextRun:     task.NextRun.Format(time.RFC3339),
		Duration:    task.Duration.String(),
		IsRecurring: strconv.Itoa(isRecurring),
		Params:      params,
	}, nil
}

func paramsToString(params []task.Param) (string, error) {
	var paramsList []string
	for _, param := range params {
		paramStr, err := json.Marshal(param)
		if err != nil {
			return "", err
		}
		paramsList = append(paramsList, string(paramStr))
	}
	data, err := json.Marshal(paramsList)
	return string(data), err
}

func paramsFromString(funcMeta task.FunctionMeta, payload string) ([]task.Param, error) {
	var params []task.Param
	if strings.TrimSpace(payload) == "" {
		return params, nil
	}
	paramTypes := funcMeta.Params()
	var paramsStrings []string
	err := json.Unmarshal([]byte(payload), &paramsStrings)
	if err != nil {
		return params, err
	}
	for i, paramStr := range paramsStrings {
		paramType := paramTypes[i]
		target := reflect.New(paramType)
		err := json.Unmarshal([]byte(paramStr), target.Interface())
		if err != nil {
			return params, err
		}
		param := reflect.Indirect(target).Interface().(task.Param)
		params = append(params, param)
	}

	return params, nil
}
