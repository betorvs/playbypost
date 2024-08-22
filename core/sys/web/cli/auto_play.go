package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	autoPlay string = "autoplay"
)

func (c *Cli) CreateAutoPlay(text string, storyID int, solo bool) ([]byte, error) {
	a := types.AutoPlayStart{
		StoryID: storyID,
		Text:    text,
		Solo:    solo,
	}
	body, err := json.Marshal(a)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(autoPlay, body)
	return res, err
}

func (c *Cli) AddNextEncounter(autoPlayID, encounterID, nextEncounterID int, text string) ([]byte, error) {
	a := types.AutoPlayNext{
		AutoPlayID:      autoPlayID,
		EncounterID:     encounterID,
		NextEncounterID: nextEncounterID,
		Text:            text,
	}
	body, err := json.Marshal(a)
	if err != nil {
		return []byte{}, err
	}
	next := fmt.Sprintf("%s/next", autoPlay)
	res, err := c.postGeneric(next, body)
	return res, err
}
