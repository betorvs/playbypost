package parser

import (
	"fmt"
	"strings"

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
	CloseStage                = "close-stage"
	RollInitiative            = "roll-initiative"
	CurrentInitiative         = "current-initiative"
	Task                      = "task"
)

func ParserOptions(storyteller bool, running types.RunningStage) []types.Options {
	opts := []types.Options{}

	if storyteller {
		// storyteller options
		if len(running.Encounter.NPC) > 0 {
			if running.Encounter.InitiativeID == 0 {
				opts = append(opts, types.Options{ID: running.Stage.StorytellerID, Name: strings.ToTitle(RollInitiative), Value: fmt.Sprintf("%s:%d", RollInitiative, running.Encounter.ID)})
				for _, v := range running.Encounter.NPC {
					opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(ActAsNPC), Value: fmt.Sprintf("%s-%s:%d", ActAsNPC, v.Name, v.ID)})
				}
			} else {
				opts = append(opts, types.Options{ID: running.Encounter.InitiativeID, Name: strings.ToTitle(CurrentInitiative), Value: fmt.Sprintf("%s:%d", CurrentInitiative, running.Encounter.ID)})
				for _, v := range running.Encounter.NPC {
					for _, p := range running.Encounter.PC {
						opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(fmt.Sprintf("%s %s", AttackPlayer, p.Name)), Value: fmt.Sprintf("%s-%s-npc-%s:%d", AttackPlayer, p.Name, v.Name, p.ID)})
					}
					// healt status for npc
					opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(fmt.Sprintf("%s %s", HealthStatus, v.Name)), Value: fmt.Sprintf("%s-npc-%s:%d", HealthStatus, v.Name, v.ID)})
				}
			}

		}
		if len(running.Encounters) > 0 {
			for _, v := range running.Encounters {
				p := types.PhaseAtoi(v.Phase)
				opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(changeEncounterText(p.NextPhase().String())), Value: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), v.ID)})
			}
		}
		p := types.PhaseAtoi(running.Encounter.Phase)
		if p == types.Running {
			opts = append(opts, types.Options{ID: running.Encounter.ID, Name: strings.ToTitle(changeEncounterText(p.NextPhase().String())), Value: fmt.Sprintf("%s:%d", changeEncounterText(p.NextPhase().String()), running.Encounter.ID)})
		}
		return opts
	}
	// player options
	count := len(running.Encounter.Options)
	if count > 0 {
		for _, v := range running.Encounter.Options {
			opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(fmt.Sprintf("%s %s", Task, v.Name)), Value: fmt.Sprintf("%s-%s:%d", Task, v.Name, v.ID)})
		}
	}
	if running.Encounter.InitiativeID != 0 {
		for _, v := range running.Encounter.NPC {
			opts = append(opts, types.Options{ID: v.ID, Name: strings.ToTitle(fmt.Sprintf("%s %s", AttackNPC, v.Name)), Value: fmt.Sprintf("%s-%s:%d", AttackNPC, v.Name, v.ID)})
		}
		opts = append(opts, types.Options{ID: running.Encounter.InitiativeID, Name: strings.ToTitle(CurrentInitiative), Value: fmt.Sprintf("%s:%d", CurrentInitiative, running.Encounter.ID)})
		// health status for player
		opts = append(opts, types.Options{ID: running.Encounter.InitiativeID, Name: strings.ToTitle(fmt.Sprintf("%s %s", HealthStatus, running.Players.Name)), Value: fmt.Sprintf("%s-%s:%d", HealthStatus, running.Players.Name, running.Encounter.InitiativeID)})
	}

	return opts
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
