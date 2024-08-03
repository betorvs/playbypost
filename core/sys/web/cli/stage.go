package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	stage string = "stage"
)

func (c *Cli) GetStage() ([]types.Stage, error) {
	var t []types.Stage
	body, err := c.getGeneric(stage)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (c *Cli) CreateStage(text, userID string, storyID int) ([]byte, error) {
	s := types.Stage{
		StoryID: storyID,
		Text:    text,
		UserID:  userID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(stage, body)
	return res, err
}
