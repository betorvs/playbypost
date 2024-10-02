package types

type Next struct {
	ID              int       `json:"id"`
	UpstreamID      int       `json:"upstream_id"`
	EncounterID     int       `json:"encounter_id"`
	NextEncounterID int       `json:"next_encounter_id"`
	Text            string    `json:"text"`
	Objective       Objective `json:"objective"`
}

type Objective struct {
	ID     int    `json:"id"`
	Kind   string `json:"kind"`
	Values []int  `json:"values"`
}

func Objectives() []string {
	return []string{ObjectiveDefault, ObjectiveDiceRoll, ObjectiveTaskOkay, ObjectiveVictory}
}
