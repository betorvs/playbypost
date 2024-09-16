package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	task string = "task"
)

func (c *Cli) GetTask() ([]types.Task, error) {
	var t []types.Task
	body, err := c.getGeneric(task)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (c *Cli) CreateTask(description, ability, skill string, kind, target int) ([]byte, error) {
	s := types.Task{
		Description: description,
		Kind:        types.TaskKindAtoi(kind),
		Ability:     ability,
		Skill:       skill,
		Target:      target,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(task, body)
	return res, err
}
