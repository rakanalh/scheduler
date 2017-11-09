package scheduler

import (
	"encoding/json"
	"log"
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
	tasks_maps, err := sb.store.Fetch()
	if err != nil {
		return []*task.Task{}, err
	}
	var tasks []*task.Task
	for _, task_map := range tasks_maps {
		lastRun, err := time.Parse(time.RFC3339, task_map["last_run"])
		if err != nil {
			return nil, err
		}

		nextRun, err := time.Parse(time.RFC3339, task_map["next_run"])
		if err != nil {
			return nil, err
		}

		duration, err := time.ParseDuration(task_map["duration"])
		if err != nil {
			return nil, err
		}

		isRecurring, err := strconv.Atoi(task_map["is_recurring"])
		if err != nil {
			return nil, err
		}

		funcMeta, err := sb.funcRegistry.Get(task_map["name"])
		if err != nil {
			return nil, err
		}

		params, err := paramsFromString(funcMeta, task_map["params"])
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
		return nil, err
	}

	isRecurring := 0
	if task.IsRecurring {
		isRecurring = 1
	}

	return storage.TaskAttributes{
		"hash":         string(task.Hash()),
		"name":         task.Func.Name,
		"last_run":     task.LastRun.Format(time.RFC3339),
		"next_run":     task.NextRun.Format(time.RFC3339),
		"duration":     task.Duration.String(),
		"is_recurring": strconv.Itoa(isRecurring),
		"params":       params,
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
	return strings.Join(paramsList, ","), nil
}

func paramsFromString(funcMeta task.FunctionMeta, payload string) ([]task.Param, error) {
	var params []task.Param
	if strings.TrimSpace(payload) == "" {
		return params, nil
	}
	paramTypes := funcMeta.Params()
	paramsStrings := strings.Split(payload, ",")
	for i, paramStr := range paramsStrings {
		paramType := paramTypes[i]
		target := reflect.New(paramType)
		err := json.Unmarshal([]byte(paramStr), target.Interface())
		if err != nil {
			log.Fatal(err)
		}
		param := reflect.Indirect(target).Interface().(task.Param)
		params = append(params, param)
	}

	return params, nil
}
