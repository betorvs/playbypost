package cli

import (
	"encoding/json"
	"fmt"

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

func (c *Cli) CreateStage(text, userID string, storyID, creatorID int) ([]byte, error) {
	s := types.Stage{
		StoryID:   storyID,
		CreatorID: creatorID,
		Text:      text,
		UserID:    userID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(stage, body)
	return res, err
}

func (c *Cli) AddEncounterToStage(text string, storyId, stageID, encounterID int) ([]byte, error) {
	s := types.EncounterAssociation{
		StageID:     stageID,
		EncounterID: encounterID,
		StoryID:     storyId,
		Text:        text,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	encounter := fmt.Sprintf("%s/encounter", stage)
	res, err := c.postGeneric(encounter, body)
	return res, err
}

func (c *Cli) StartStage(stageID int, channelID string) ([]byte, error) {
	s := types.Channel{
		StageID: stageID,
		Channel: channelID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	start := fmt.Sprintf("%s/channel", stage)
	res, err := c.postGeneric(start, body)
	return res, err
}

func (c *Cli) AssignTask(text string, stageID, taskID, storytellerID, encounterID int) ([]byte, error) {
	s := types.RunningTask{
		StageID:       stageID,
		TaskID:        taskID,
		StorytellerID: storytellerID,
		EncounterID:   encounterID,
		Text:          text,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	task := fmt.Sprintf("%s/encounter/task", stage)
	res, err := c.postGeneric(task, body)
	return res, err
}
