package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	ActAsNPC                  = "act-as-npc"
	AttackNPC                 = "attack-npc"
	AttackPlayer              = "attack-player"
	HealthStatus              = "health-status"
	ChangeEncounter           = "change-encounter"
	ChangeEncounterToStarted  = "change-encounter-to-started"
	ChangeEncounterToRunning  = "change-encounter-to-running"
	ChangeEncounterToFinished = "change-encounter-to-finished"
	RollInitiative            = "roll-initiative"
	CurrentInitiative         = "current-initiative"
	Task                      = "task"
)

func ParserOptions(storyteller bool, running types.RunningStage) []types.GenericIDName {
	encOptions := []types.GenericIDName{}

	if storyteller {
		// storyteller options
		if len(running.Encounter.NPC) > 0 {
			if running.Encounter.InitiativeID == 0 {
				encOptions = append(encOptions, types.GenericIDName{ID: running.Stage.StorytellerID, Name: fmt.Sprintf("%s:%d", RollInitiative, running.Encounter.ID)})
				for _, v := range running.Encounter.NPC {
					encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s:%s", ActAsNPC, v.Name)})
				}
			} else {
				encOptions = append(encOptions, types.GenericIDName{ID: running.Encounter.InitiativeID, Name: fmt.Sprintf("%s:%d", CurrentInitiative, running.Encounter.ID)})
				for _, v := range running.Encounter.NPC {
					for _, p := range running.Encounter.PC {
						encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s-%s-npc-%s:%d", AttackPlayer, p.Name, v.Name, p.ID)})
					}
					// healt status for npc
					encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s-npc-%s:%d", HealthStatus, v.Name, v.ID)})
				}
			}

		}
		if len(running.Encounters) > 0 {
			for _, v := range running.Encounters {
				p := types.PhaseAtoi(v.Phase)
				encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), v.ID)})
			}
		}
		p := types.PhaseAtoi(running.Encounter.Phase)
		if p == types.Running {
			encOptions = append(encOptions, types.GenericIDName{ID: running.Encounter.ID, Name: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), running.Encounter.ID)})
		}
		return encOptions
	}
	// player options
	count := len(running.Encounter.Options)
	if count > 0 {
		for _, v := range running.Encounter.Options {
			encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s-%s:%d", Task, v.Name, v.ID)})
		}
	}
	if running.Encounter.InitiativeID != 0 {
		for _, v := range running.Encounter.NPC {
			encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s-%s:%d", AttackNPC, v.Name, v.ID)})
		}
		encOptions = append(encOptions, types.GenericIDName{ID: running.Encounter.InitiativeID, Name: fmt.Sprintf("%s:%d", CurrentInitiative, running.Encounter.ID)})
		// health status for player
		encOptions = append(encOptions, types.GenericIDName{ID: running.Encounter.InitiativeID, Name: fmt.Sprintf("%s-%s:%d", HealthStatus, running.Players.Name, running.Encounter.InitiativeID)})
	}

	return encOptions
}

func changeEncounterText(phase string) string {
	switch phase {
	case types.Started.String():
		return ChangeEncounterToStarted
	case types.Running.String():
		return ChangeEncounterToRunning
	case types.Finished.String():
		return ChangeEncounterToFinished
	}
	return ""
}
