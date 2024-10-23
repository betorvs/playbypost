package validator

import "github.com/betorvs/playbypost/core/sys/web/types"

func (v *Validator) ValidateStage(stage *types.StageAggregated, hashID string) {
	ok := false
	a := Analitycs{}
	// get all next encounters from auto play
	if stage.Story.ID != 0 {
		valid, storyAnalitics, encounters := v.checkStory(stage.Story.ID)
		a.Results = append(a.Results, storyAnalitics.Results...)
		if valid {
			ok = valid
			v.logger.Debug("story is valid", "story_id", stage.Story.ID)
		}
		stageEncounter, err := v.db.GetStageEncountersByStageID(v.ctx, stage.Stage.ID)
		if err != nil {
			v.logger.Error("error getting stage encounters", "error", err)
		}
		shouldHaveStageEncounters := false
		if len(stageEncounter) > 0 {
			shouldHaveStageEncounters = true
			firstEncounter, listLastEncounters := parserEncountersFirstAndListOfLast(encounters)
			if firstEncounter == 0 {
				a.Results = append(a.Results, "stage should have a first encounter")
				ok = false
			}
			if len(listLastEncounters) == 0 {
				a.Results = append(a.Results, "stage should have a last encounter")
				ok = false
			}
		}
		if !shouldHaveStageEncounters {
			a.Results = append(a.Results, "stage should have encounters added")
			ok = false
		}
	}

	v.UpdateRequest(hashID, ok, "stage validated", a)
}
