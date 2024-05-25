package types

import "github.com/betorvs/playbypost/core/initiative"

type Scene struct {
	Story      Story
	Encounter  Encounters
	Tasks      map[string]Task
	Combat     bool
	Initiative *initiative.Initiative
}
