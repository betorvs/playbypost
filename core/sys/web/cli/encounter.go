package cli

import (
	"encoding/json"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	encounter string = "encounter"
)

func (c *Cli) GetEncounters() ([]types.Encounter, error) {
	var list []types.Encounter
	body, err := c.getGeneric(encounter)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (c *Cli) GetEncounterByID(id int) (types.Encounter, error) {
	var list types.Encounter
	enc := fmt.Sprintf("%s/%d", encounter, id)
	body, err := c.getGeneric(enc)
	if err != nil {
		return list, err
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (c *Cli) CreateEncounter(title, announcement, notes string, storyID, writerID int, first, last bool) ([]byte, error) {
	s := types.Encounter{
		Title:          title,
		Announcement:   announcement,
		Notes:          notes,
		StoryID:        storyID,
		WriterID:       writerID,
		FirstEncounter: first,
		LastEncounter:  last,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	res, err := c.postGeneric(encounter, body)
	return res, err
}

func (c *Cli) ChangeEncounterPhase(id, phase int) error {
	enc := fmt.Sprintf("%s/%d/%d", encounter, id, phase)
	err := c.putGeneric(enc)
	return err
}

func (c *Cli) AddParticipants(encounterID int, npc bool, IDs []int) ([]byte, error) {
	s := types.Participants{
		Identifies:  IDs,
		EncounterID: encounterID,
		NPC:         npc,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	enc := fmt.Sprintf("%s/participants", encounter)
	res, err := c.postGeneric(enc, body)
	return res, err
}

func (c *Cli) UpdateEncounter(title, announcement, notes string, id, storyID, writerID int, first, last bool) ([]byte, error) {
	s := types.Encounter{
		ID:             id,
		Title:          title,
		Announcement:   announcement,
		Notes:          notes,
		StoryID:        storyID,
		WriterID:       writerID,
		FirstEncounter: first,
		LastEncounter:  last,
	}
	body, err := json.Marshal(s)
	if err != nil {
		return []byte{}, err
	}
	kind := fmt.Sprintf("%s/%d", encounter, id)
	res, err := c.putGenericWithHeaders(kind, body)
	return res, err
}
