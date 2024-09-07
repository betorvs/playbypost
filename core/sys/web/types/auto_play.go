package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type AutoPlayStart struct {
	StoryID int    `json:"story_id"`
	Text    string `json:"text"`
	Solo    bool   `json:"solo"`
}

type AutoPlay struct {
	ID      int    `json:"id"`
	StoryID int    `json:"story_id"`
	Text    string `json:"text"`
	Solo    bool   `json:"solo"`
}

type AutoPlayGroup struct {
	ID     int    `json:"id"`
	UserID string `json:"user_id"`
}

func (a AutoPlayGroup) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *AutoPlayGroup) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type AutoPlayNext struct {
	ID              int    `json:"id"`
	AutoPlayID      int    `json:"auto_play_id"`
	EncounterID     int    `json:"encounter_id"`
	NextEncounterID int    `json:"next_encounter_id"`
	Text            string `json:"text"`
}

type AutoPlayEncounterList struct {
	EncounterList []Options                   `json:"encounter_list"`
	Link          []AutoPlayEncounterWithNext `json:"link"`
}

type AutoPlayEncounterWithNext struct {
	Encounter     string `json:"encounter"`
	NextEncounter string `json:"next_encounter"`
}

type AutoPlayOptions struct {
	AutoPlay
	EncodingKey       string          `json:"encoding_key"`
	AutoPlayChannelID int             `json:"auto_play_channel_id"`
	ChannelID         string          `json:"channel_id"`
	Group             []AutoPlayGroup `json:"group"`
	NextEncounters    []AutoPlayNext  `json:"next_encounters"`
}

type AutoPlayActivities struct {
	ID          int     `json:"id"`
	Actions     Actions `json:"actions"`
	AutoPlayID  int     `json:"auto_play_id"`
	EncounterID int     `json:"encounter_id"`
	Processed   bool    `json:"processed"`
}
