package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	story string = "story"
)

func (c *Cli) GetStory() ([]types.Story, error) {
	var s []types.Story
	body, err := c.getGeneric(story)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (c *Cli) CreateStory(title, announcement, notes string, writerID int) ([]byte, error) {
	s := types.Story{
		Title:        title,
		Announcement: announcement,
		Notes:        notes,
		WriterID:     writerID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(story, body)
	return res, err
}
