package cli

import (
	"encoding/json"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	initiative string = "initiative"
)

func (c *Cli) CreateInitiative(name string, encounterID int, npc bool) ([]byte, error) {
	s := types.Initiative{
		Name:        name,
		EncounterID: encounterID,
		NPC:         npc,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(initiative, body)
	return res, err
}
