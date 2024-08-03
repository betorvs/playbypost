package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	ActAsNPC                  = "act-as-npc"
	AttackNPC                 = "attack-npc"
	AttackPlayer              = "attack-player"
	ChangeEncounterToStarted  = "change-encounter-to-started"
	ChangeEncounterToRunning  = "change-encounter-to-running"
	ChangeEncounterToFinished = "change-encounter-to-finished"
	RollInitiative            = "roll-initiative"
	Task                      = "task"
)

func ParserOptions(storyteller bool, running types.RunningStage) []types.GenericIDName {
	encOptions := []types.GenericIDName{}

	if storyteller {
		// storyteller options
		count := 1
		if len(running.Encounter.NPC) > 0 {
			encOptions = append(encOptions, types.GenericIDName{ID: count, Name: RollInitiative})
			for _, v := range running.Encounter.NPC {
				count++
				encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("%s:%s", ActAsNPC, v.Name)})
			}
		}
		if len(running.Encounters) > 0 {
			for _, v := range running.Encounters {
				count++
				p := types.PhaseAtoi(v.Phase)
				encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), v.ID)})
			}
		}
		p := types.PhaseAtoi(running.Encounter.Phase)
		if p == types.Running {
			encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), running.Encounter.ID)})
		}
		return encOptions
	}
	// player options
	count := len(running.Encounter.Options)
	if count > 0 {
		for _, v := range running.Encounter.Options {
			count++
			encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("%s-%s:%d", Task, v.Name, v.ID)})
		}
		for _, v := range running.Encounter.NPC {
			count++
			encOptions = append(encOptions, types.GenericIDName{ID: count, Name: fmt.Sprintf("%s:%s", ActAsNPC, v.Name)})
		}
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
