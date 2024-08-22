package parser

import (
	"testing"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func TestParserOptions(t *testing.T) {
	storyteller := true
	running := types.RunningStage{}
	result := ParserOptions(storyteller, running)
	if len(result) != 0 {
		t.Errorf("ParserOptions failed, expected %v, got %v", 0, len(result))
	}
}

func TestChangeEncounterText(t *testing.T) {
	result := changeEncounterText("running")
	if result != ChangeEncounterToRunning {
		t.Errorf("changeEncounterText failed, expected %v, got %v", "test", result)
	}
}
