package validator

import "github.com/betorvs/playbypost/core/sys/web/types"

func (v *Validator) ValidateStory(story *types.Story, hashID string) {
	valid, a, _ := v.checkStory(story.ID)
	v.UpdateRequest(hashID, valid, "story validated", a)
}

func (v *Validator) checkStory(id int) (bool, Analitycs, []types.Encounter) {
	encounters, err := v.db.GetEncounterByStoryID(v.ctx, id)
	if err != nil {
		v.logger.Error("error getting encounters", "error", err)
	}
	// validate parameters
	valid := false
	shouldHaveEncounters := false
	oneFirstEncounter := false
	atLeastOneLastEncounter := false
	a := Analitycs{}
	if len(encounters) > 0 {
		shouldHaveEncounters = true
		for _, encounter := range encounters {
			if encounter.FirstEncounter {
				oneFirstEncounter = true
			}
			if encounter.LastEncounter {
				atLeastOneLastEncounter = true
			}
		}
	}
	if !shouldHaveEncounters {
		a.Results = append(a.Results, "story should have encounters")
	}
	if !oneFirstEncounter {
		a.Results = append(a.Results, "first encounter not found")
	}
	if !atLeastOneLastEncounter {
		a.Results = append(a.Results, "last encounter not found")
	}
	if oneFirstEncounter && atLeastOneLastEncounter && shouldHaveEncounters {
		valid = true
	}
	return valid, a, encounters
}
