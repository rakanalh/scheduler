package storage

import (
	"encoding/json"
	"fmt"

	"github.com/rakanalh/scheduler/task"
)

type ParamMarshaler interface {
	Marshal(params []task.Param) (string, error)
	Unmarshal(bytes string) ([]task.Param, error)
}

type TaskStore interface {
	Add(task *task.Task) error
	Fetch() ([]*task.Task, error)
	Remove(task *task.Task) error
}

type JsonMarshaler struct{}

func NewMarshaler() JsonMarshaler {
	return JsonMarshaler{}
}

func (m JsonMarshaler) Marshal(params []task.Param) (string, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("Could not marshal params: %s", err)
	}
	return string(b), nil
}

func (m JsonMarshaler) Unmarshal(bytes string) ([]task.Param, error) {
	var params []task.Param
	err := json.Unmarshal([]byte(bytes), &params)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal params: %s", err)
	}
	return params, nil
}
