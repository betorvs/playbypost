package parser

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

const (
	StartSolo = "start-solo"
	NextSolo  = "next-solo"
)

func ParserAutoPlaysSolo(autoPlays []types.AutoPlay) ([]types.GenericIDName, []types.Options) {
	encOptions := []types.GenericIDName{}
	opts := []types.Options{}
	for _, v := range autoPlays {
		encOptions = append(encOptions, types.GenericIDName{ID: v.ID, Name: fmt.Sprintf("%s-%s:%d", StartSolo, v.Text, v.ID)})
		opts = append(opts, types.Options{ID: v.ID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", StartSolo, v.Text, v.ID)})
	}
	return encOptions, opts
}

func ParserAutoPlaysNext(autoPlays []types.AutoPlayNext) ([]types.GenericIDName, []types.Options) {
	encOptions := []types.GenericIDName{}
	opts := []types.Options{}
	for _, v := range autoPlays {
		encOptions = append(encOptions, types.GenericIDName{ID: v.AutoPlayID, Name: fmt.Sprintf("%s-%s:%d", NextSolo, v.Text, v.NextEncounterID)})
		opts = append(opts, types.Options{ID: v.AutoPlayID, Name: v.Text, Value: fmt.Sprintf("%s-%s:%d", NextSolo, v.Text, v.NextEncounterID)})
	}
	return encOptions, opts
}
