package types

type GeneratePlayer struct {
	PlayerID int    `json:"player_id,omitempty"`
	StoryID  int    `json:"story_id"`
	Name     string `json:"name"`
}

type Players struct {
	ID        int               `json:"id,omitempty"`
	Name      string            `json:"name"`
	Abilities map[string]int    `json:"abilities"`
	Skills    map[string]int    `json:"skills"`
	RPG       string            `json:"rpg"`
	Extension map[string]int    `json:"extension"`
	Details   map[string]string `json:"details"`
	Destroyed bool              `json:"destroyed"`
}
