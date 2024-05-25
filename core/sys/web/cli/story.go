package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	story string = "story"
)

func (c *Cli) GetStory() ([]types.Story, error) {
	var users []types.Story
	body, err := c.getGeneric(story)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func (c *Cli) CreateStory(title, announcement, notes string, storytellerID int) ([]byte, error) {
	s := types.Story{
		Title:         title,
		Announcement:  announcement,
		Notes:         notes,
		StorytellerID: storytellerID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(story, body)
	return res, err
}
