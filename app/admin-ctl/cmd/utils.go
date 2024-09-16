package cmd

import (
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (app *application) getEncounters() ([]types.Encounter, error) {
	encounters, err := app.Web.GetEncounters()
	if err != nil {
		app.Logger.Error("get encounters", "error", err.Error())
		return []types.Encounter{}, err
	}
	return encounters, nil
}

func (app *application) findEncounterID(encounterTitle string, encounters []types.Encounter) int {
	for _, v := range encounters {
		if v.Title == encounterTitle {
			return v.ID
		}
	}
	return 0
}

func (app *application) getStories() ([]types.Story, error) {
	stories, err := app.Web.GetStory()
	if err != nil {
		app.Logger.Error("get stories", "error", err.Error())
		return []types.Story{}, err
	}
	return stories, nil
}

func (app *application) findStoryID(storyTitle string, stories []types.Story) int {
	for _, v := range stories {
		if v.Title == storyTitle {
			return v.ID
		}
	}
	return 0
}

func (app *application) getStages() ([]types.Stage, error) {
	stages, err := app.Web.GetStage()
	if err != nil {
		app.Logger.Error("get stages", "error", err.Error())
		return []types.Stage{}, err
	}
	return stages, nil
}

func (app *application) findStageID(stageTitle string, stages []types.Stage) int {
	for _, v := range stages {
		if v.Text == stageTitle {
			return v.ID
		}
	}
	return 0
}
