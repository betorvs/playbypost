package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
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
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	LastUpdateAt time.Time `json:"priority_timestamp"`
	Interactions int       `json:"interactions"`
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
	NextEncounters    []Next          `json:"next_encounters"`
}
