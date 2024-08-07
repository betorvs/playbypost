package types

type Initiative struct {
	Channel     string `json:"channel"`
	EncounterID int    `json:"encounter_id"`
	UserID      string `json:"user_id,omitempty"`
	// StoryID     int      `json:"story_id"`
	// PlayersID   []int `json:"players_id"`
	// NPC         bool     `json:"is_npc"`
	// NonPlayerID []int `json:"npcs"`
}

type InitiativeShort struct {
	ID           int      `json:"id,omitempty"`
	Name         string   `json:"name"`
	NextPlayer   string   `json:"next_player"`
	Participants []string `json:"participants"`
}
