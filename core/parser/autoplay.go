package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	StartSolo = "start-solo"
	NextSolo  = "next-solo"
)

func ParserAutoPlaysSolo(autoPlays []types.AutoPlay) []types.Options {
	opts := []types.Options{}
	for _, v := range autoPlays {
		opts = append(opts, types.Options{ID: v.ID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", StartSolo, v.Text, v.ID)})
	}
	return opts
}

func ParserAutoPlaysNext(autoPlays []types.AutoPlayNext) []types.Options {
	opts := []types.Options{}
	for _, v := range autoPlays {
		opts = append(opts, types.Options{ID: v.AutoPlayID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", NextSolo, v.Text, v.NextEncounterID)})
	}
	return opts
}
