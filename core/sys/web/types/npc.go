package types

type GenerateNPC struct {
	StorytellerID int    `json:"storyteller_id,omitempty"`
	StageID       int    `json:"stage_id"`
	EncounterID   int    `json:"encounter_id,omitempty"`
	Name          string `json:"name"`
}

func NewNPC() *Players {
	return &Players{
		Abilities: make(map[string]int),
		Skills:    make(map[string]int),
		Extension: map[string]int{},
	}
}
