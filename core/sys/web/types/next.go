package types

import (
	"fmt"
	"slices"
)

const (
	UpstreamKindAutoPlay = "auto_play"
	UpstreamKindStage    = "stage"
)

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

func ValidateNextSlice(s []Next, upstreamKind string) ([]Next, error) {
	n := []Next{}
	if len(s) == 0 {
		return s, fmt.Errorf("next encounters cannot be empty")
	}
	counter := make(map[int]int)
	counterNext := make(map[int]int)
	for k, obj := range s {
		if obj.NextEncounterID == 0 || obj.EncounterID == 0 || obj.UpstreamID == 0 {
			return s, fmt.Errorf("next encounter id, encounter id and auto play id cannot be empty")
		}
		if obj.Objective.Kind != "" && !slices.Contains(Objectives(), obj.Objective.Kind) {
			return s, fmt.Errorf("objective kind %v is not valid", obj.Objective.Kind)
		}
		if obj.Objective.Kind == "" && len(obj.Objective.Values) == 0 {
			obj.Objective.Kind = ObjectiveDefault
			obj.Objective.Values = []int{0}
		}
		if obj.Objective.Kind == ObjectiveDiceRoll && len(s) <= 1 {
			return s, fmt.Errorf("objective dice roll should have at least 2 next encounters")
		}
		switch upstreamKind {
		case "auto_play":
			// check for "dice roll" objective and auto complete with possible values
			if obj.Objective.Kind == ObjectiveDiceRoll {
				obj.Objective.Values = SplitDiceNextObjctive(k, len(s))
			}
			if obj.Objective.Kind == ObjectiveTaskOkay || obj.Objective.Kind == ObjectiveVictory {
				return s, fmt.Errorf("objective task okay and victory are not valid for auto play")
			}
		case "stage":
			if obj.Objective.Kind == ObjectiveTaskOkay || obj.Objective.Kind == ObjectiveVictory {
				if len(s) != 2 {
					return s, fmt.Errorf("objective task okay and victory should have 2 next encounters")
				}
			}
		}
		counter[obj.EncounterID]++
		counterNext[obj.NextEncounterID]++

		n = append(n, obj)
	}
	if len(counter) != 1 {
		return s, fmt.Errorf("encounter id should be unique in all next encounters")
	}
	for k, v := range counterNext {
		if v > 1 {
			return s, fmt.Errorf("next encounter id %v should be unique", k)
		}
	}

	return n, nil
}
