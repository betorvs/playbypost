package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	initiative string = "initiative"
)

func (c *Cli) CreateInitiative(userID, channel string, encounterID int) ([]byte, error) {
	s := types.Initiative{
		Channel:     channel,
		EncounterID: encounterID,
		UserID:      userID,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(initiative, body)
	return res, err
}
