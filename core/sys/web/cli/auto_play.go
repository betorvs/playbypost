package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	autoPlay string = "autoplay"
)

func (c *Cli) GetAutoPlay() ([]types.AutoPlay, error) {
	var a []types.AutoPlay
	body, err := c.getGeneric(autoPlay)
	if err != nil {
		return a, err
	}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return a, err
	}
	return a, nil
}

func (c *Cli) CreateAutoPlay(text string, storyID, creatorID int, solo bool) ([]byte, error) {
	a := types.AutoPlayStart{
		StoryID:   storyID,
		CreatorID: creatorID,
		Text:      text,
		Solo:      solo,
	}
	body, err := json.Marshal(a)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(autoPlay, body)
	return res, err
}

func (c *Cli) AddNextEncounter(a types.Next) ([]byte, error) {
	sliceNext := []types.Next{a}
	body, err := json.Marshal(sliceNext)
	if err != nil {
		return []byte{}, err
	}
	next := fmt.Sprintf("%s/next", autoPlay)
	res, err := c.postGeneric(next, body)
	return res, err
}

func (c *Cli) PublishAutoPlay(id int) ([]byte, error) {
	publish := fmt.Sprintf("%s/publish/%d", autoPlay, id)
	res, err := c.putEmptyBodyGeneric(publish)
	return res, err
}
