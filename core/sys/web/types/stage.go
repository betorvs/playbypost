package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Stage struct {
	ID            int    `json:"id"`
	StoryID       int    `json:"story_id"`
	Text          string `json:"text"`
	UserID        string `json:"user_id,omitempty"`
	StorytellerID int    `json:"storyteller_id,omitempty"`
}

type StageAggregated struct {
	Stage   Stage   `json:"stage"`
	Story   Story   `json:"story"`
	Channel Channel `json:"channel"`
}

type Channel struct {
	StageID int    `json:"stage_id"`
	Channel string `json:"channel"`
	Active  bool   `json:"active"`
}

type EncounterAssociation struct {
	StageID     int    `json:"stage_id"`
	StoryID     int    `json:"story_id"`
	EncounterID int    `json:"encounter_id"`
	Text        string `json:"text"`
}

type StageEncounter struct {
	Encounter
	Phase         int
	Text          string    `json:"text"`
	StorytellerID int       `json:"storyteller_id,omitempty"`
	InitiativeID  int       `json:"initiative_id"`
	PC            []Options `json:"pc"`
	NPC           []Options `json:"npc"`
	Options       []Options `json:"options"`
}

type Participants struct {
	Identifies  []int `json:"identifies"`
	EncounterID int   `json:"encounter_id"`
	NPC         bool  `json:"is_npc"`
}

type RunningStage struct {
	StageAggregated
	Players    Players
	Encounter  StageEncounter
	Encounters []StageEncounter
}

// RunningTask represents a task that is currently running.
type RunningTask struct {
	StageID       int    `json:"stage_id"`
	TaskID        int    `json:"task_id"`
	StorytellerID int    `json:"storyteller_id"`
	EncounterID   int    `json:"encounter_id"`
	Text          string `json:"text"`
}

type Actions map[string]string

func NewActions() Actions {
	return make(Actions)
}

func (a Actions) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Actions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Phase int

const (
	Waiting  Phase = iota // 0
	Started               // 1
	Running               // 2
	Finished              // 3
)

func PhaseAtoi(i int) Phase {
	return Phase(i)
}

func (p Phase) String() string {
	switch p {
	case Waiting:
		return "waiting"
	case Started:
		return "started"
	case Running:
		return "running"
	case Finished:
		return "finished"
	}
	return "waiting"
}

func (p Phase) NextPhase() Phase {
	switch p {
	case Waiting:
		return Started
	case Started:
		return Running
	case Running:
		return Finished
	case Finished:
		return Finished
	}
	return Waiting
}
