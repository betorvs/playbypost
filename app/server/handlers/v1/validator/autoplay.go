package validator

import (
	"fmt"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (v *Validator) ValidateAutoPlay(autoPlay *types.AutoPlay, hashID string) {
	ok := false
	a := Analitycs{}
	// get all next encounters from auto play
	story, err := v.getStoryByID(autoPlay.StoryID)
	if err != nil {
		v.logger.Error("error getting story", "error", err)
	}
	shouldHaveNextEncounters := false
	if story.ID != 0 {
		valid, storyAnalitics, encounters := v.checkStory(story.ID)
		a.Results = append(a.Results, storyAnalitics.Results...)
		if valid {
			ok = valid
			// v.logger.Info("story is valid", "story_id", story.ID)
			// check if next encounters are valid
			next, err := v.db.GetNextEncounterByAutoPlayID(v.ctx, autoPlay.ID)
			if err != nil {
				v.logger.Error("error getting next encounters", "error", err)
			}
			if len(next) > 0 {
				shouldHaveNextEncounters = true
			}
			firstEncounter, listLastEncounters := parserEncountersFirstAndListOfLast(encounters)
			result := parserAutoPlayNext(next, firstEncounter, listLastEncounters)
			if !result.LastEncounters {
				a.Results = append(a.Results, fmt.Sprintf("expected %d last encounters, found %d paths and these encounter IDs %v", len(listLastEncounters), result.TotalPaths, result.LastEncounterFound))
				ok = false
			}
			if !result.ObjectivesMatch {
				a.Results = append(a.Results, "objectives does not match. In case you use dice roll objective in encounter ID 7, all encounters linked to 7, should have a dice roll objective")
				ok = false
			}
		}
	}
	if !shouldHaveNextEncounters {
		a.Results = append(a.Results, "auto play should have next encounters")
		ok = false
	}

	v.UpdateRequest(hashID, ok, "autoplay validated", a)
}

type AutoPlayNextResult struct {
	Path               map[string][]int
	TotalPaths         int
	LastEncounters     bool
	LastEncountersUsed bool
	LastEncounterFound []int
	ObjectivesMatch    bool
}

func parserAutoPlayNext(next []types.Next, first int, last []int) AutoPlayNextResult {
	var result AutoPlayNextResult
	path := make(map[string][]int)
	objectiveCount := make(map[int]int)
	linkedCount := make(map[int]int)
	for k, n := range next {
		if n.EncounterID == first {
			loop := fmt.Sprintf("loop-%d", k)
			path[loop] = append(path[loop], n.EncounterID)
			recursivePath(k, next, n.EncounterID, path, loop)
		}
		linkedCount[n.EncounterID]++
		if n.Objective.Kind == types.ObjectiveDiceRoll {
			objectiveCount[n.EncounterID]++
		}
	}
	result.Path = path
	result.TotalPaths = len(path)
	result.LastEncounters = true
	result.LastEncountersUsed = true
	result.ObjectivesMatch = true

	if len(path) != len(last) {
		result.LastEncounters = false
	}
	for _, v := range path {
		result.LastEncounterFound = append(result.LastEncounterFound, v[len(v)-1])
		if !slices.Contains(last, v[len(v)-1]) {
			result.LastEncountersUsed = false
		}
	}
	if len(objectiveCount) > 0 {
		for k, v := range objectiveCount {
			if linkedCount[k] != v {
				result.ObjectivesMatch = false
			}
		}
	}
	return result
}

func recursivePath(position int, next []types.Next, encounterID int, path map[string][]int, name string) {
	// fmt.Println("position", position, "encounterID", encounterID)
	if position == len(next) {
		// path[name] = append(path[name], next[position].NextEncounterID)
		return
	}
	if next[position].EncounterID == encounterID {
		path[name] = append(path[name], next[position].NextEncounterID)
		recursivePath(position+1, next, next[position].NextEncounterID, path, name)
		return
	}

	recursivePath(position+1, next, encounterID, path, name)

}

func (v *Validator) getStoryByID(id int) (*types.Story, error) {
	story, err := v.db.GetStoryByID(v.ctx, id)
	if err != nil {
		v.logger.Error("error getting story", "error", err)
		return nil, err
	}
	return &story, nil
}

func parserEncountersFirstAndListOfLast(encounters []types.Encounter) (int, []int) {
	firstEncounter := 0
	listLastEncounters := []int{}
	for _, encounter := range encounters {
		if encounter.FirstEncounter {
			firstEncounter = encounter.ID
		}
		if encounter.LastEncounter {
			listLastEncounters = append(listLastEncounters, encounter.ID)
		}
	}

	return firstEncounter, listLastEncounters
}
