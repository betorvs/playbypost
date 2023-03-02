package rules

import (
	"testing"

	"github.com/betorvs/playbypost/core/types"
)

func TestCombatModifiers(t *testing.T) {
	t.Log("Checking all combat modifiers in loop")
	for i := Zero; i <= KneelingSitting; i++ {
		res := i.String()
		if res == "" && res != types.Unknown {
			t.Errorf("\t combat modifiers should not be empty %s", res)
		}
	}
}
