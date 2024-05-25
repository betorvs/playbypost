package types

type Initiative struct {
	Name        string   `json:"name"`
	StoryID     int      `json:"story_id"`
	EncounterID int      `json:"encounter_id"`
	PlayersID   []int    `json:"players_id"`
	NPC         bool     `json:"is_npc"`
	NonPlayerID []string `json:"npcs"`
}

type InitiativeShort struct {
	ID           int      `json:"id,omitempty"`
	Name         string   `json:"name"`
	NextPlayer   string   `json:"next_player"`
	Participants []string `json:"participants"`
}
