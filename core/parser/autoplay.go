package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	StartSolo    = "start-solo"
	NextSolo     = "next-solo"
	DiceRollSolo = "dice-roll-solo"
)

func ParserAutoPlaysSolo(autoPlays []types.AutoPlay) []types.Options {
	opts := []types.Options{}
	for _, v := range autoPlays {
		opts = append(opts, types.Options{ID: v.ID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", StartSolo, v.Text, v.ID)})
	}
	return opts
}

func ParserAutoPlaysNext(autoPlays []types.AutoPlayNext) ([]types.Options, bool) {
	ok := false
	opts := []types.Options{}
	for _, v := range autoPlays {
		switch v.Objective.Kind {
		case types.ObjectiveDiceRoll:
			if !ok {
				// it should add only one dice roll option
				opts = append(opts, types.Options{ID: v.AutoPlayID, Name: DiceRollSolo, Value: fmt.Sprintf("%s:%d", DiceRollSolo, v.EncounterID)})
				ok = true
			}
		case types.ObjectiveDefault:
			opts = append(opts, types.Options{ID: v.AutoPlayID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", NextSolo, v.Text, v.NextEncounterID)})
		}

	}
	return opts, ok
}
