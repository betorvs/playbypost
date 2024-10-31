package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	StartSolo    = "start-solo"
	NextSolo     = "next-solo"
	DiceRollSolo = "dice-roll-solo"
	StartDidatic = "didatic-start"
	NextDidatic  = "didatic-next"
	JoinDidatic  = "didatic-join"
)

func ParserAutoPlays(autoPlays []types.AutoPlay, kind string) []types.Options {
	opts := []types.Options{}
	cmd := StartSolo
	solo := true
	switch kind {
	case StartDidatic:
		cmd = StartDidatic
		solo = false
	case JoinDidatic:
		cmd = JoinDidatic
		solo = false
	}

	for _, v := range autoPlays {
		if v.Solo == solo {
			opts = append(opts, types.Options{ID: v.ID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", cmd, v.Text, v.ID)})
		}
	}
	return opts
}

func ParserAutoPlaysNext(autoPlays []types.Next, solo bool) ([]types.Options, bool) {
	ok := false
	opts := []types.Options{}
	cmd := NextSolo
	if !solo {
		cmd = NextDidatic
	}
	for _, v := range autoPlays {
		switch v.Objective.Kind {
		case types.ObjectiveDiceRoll:
			if !ok {
				// it should add only one dice roll option
				opts = append(opts, types.Options{ID: v.UpstreamID, Name: DiceRollSolo, Value: fmt.Sprintf("%s:%d", DiceRollSolo, v.EncounterID)})
				ok = true
			}
		case types.ObjectiveDefault:
			opts = append(opts, types.Options{ID: v.UpstreamID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", cmd, v.Text, v.NextEncounterID)})
		}

	}
	return opts, ok
}
