package types

type Encounters struct {
	ID            int    `json:"id,omitempty"`
	Title         string `json:"title"`
	Announcement  string `json:"announcement"`
	Notes         string `json:"notes"`
	StoryID       int    `json:"story_id"`
	// Phase         Phase  `json:"phase"`
	// Finished      bool   `json:"finished"`
	// Reward        string `json:"reward"`
	// XP            int    `json:"xp"`
	StorytellerID int    `json:"storyteller_id"`
}

type Phase int

const (
	Waiting Phase = iota
	Started
	Running
	Finished
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

type Participants struct {
	PlayersID   []int `json:"players_id"`
	EncounterID int   `json:"encounter_id"`
	NPC         bool  `json:"is_npc"`
}
