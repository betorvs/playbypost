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

// type AutoPlayEncounterList struct {
// 	EncounterList []Options                   `json:"encounter_list"`
// 	Link          []AutoPlayEncounterWithNext `json:"link"`
// }

// type AutoPlayEncounterWithNext struct {
// 	ID            int    `json:"id"`
// 	EncounterID   int    `json:"encounter_id"`
// 	NextID        int    `json:"next_id"`
// 	Encounter     string `json:"encounter"`
// 	NextEncounter string `json:"next_encounter"`
// }

type AutoPlayOptions struct {
	AutoPlay
	EncodingKey       string          `json:"encoding_key"`
	AutoPlayChannelID int             `json:"auto_play_channel_id"`
	ChannelID         string          `json:"channel_id"`
	Group             []AutoPlayGroup `json:"group"`
	NextEncounters    []Next          `json:"next_encounters"`
}

func SplitDiceNextObjctive(loop, size int) []int {
	switch size {
	case 2:
		switch loop {
		case 0:
			return []int{1, 3, 5}
		case 1:
			return []int{2, 4, 6}
		}
	case 3:
		switch loop {
		case 0:
			return []int{1, 2}
		case 1:
			return []int{3, 4}
		case 2:
			return []int{5, 6}
		}
	case 4:
		switch loop {
		case 0:
			return []int{1, 2}
		case 1:
			return []int{3, 4}
		case 2:
			return []int{5}
		case 3:
			return []int{6}
		}
	case 5:
		switch loop {
		case 0:
			return []int{1, 2}
		case 1:
			return []int{3}
		case 2:
			return []int{4}
		case 3:
			return []int{5}
		case 4:
			return []int{6}
		}
	case 6:
		switch loop {
		case 0:
			return []int{1}
		case 1:
			return []int{2}
		case 2:
			return []int{3}
		case 3:
			return []int{4}
		case 4:
			return []int{5}
		case 5:
			return []int{6}
		}
	}
	return nil
}
