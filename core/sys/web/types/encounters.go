package types

const (
	ObjectiveDefault  = "no_action"
	ObjectiveDiceRoll = "dice_roll"
	ObjectiveTaskOkay = "task_okay"
	ObjectiveVictory  = "victory"
)

type Encounter struct {
	ID             int    `json:"id,omitempty"`
	Title          string `json:"title"`
	Announcement   string `json:"announcement"`
	Notes          string `json:"notes"`
	StoryID        int    `json:"story_id"`
	WriterID       int    `json:"writer_id"`
	FirstEncounter bool   `json:"first_encounter"`
	LastEncounter  bool   `json:"last_encounter"`
}
