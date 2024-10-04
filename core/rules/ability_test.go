package rules

import (
	"log/slog"
	"os"
	"testing"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/tests/mock"
)

func TestAbilityD10HM(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rpgSystem := rpg.LoadRPGSystemsDefault(rpg.D10HM)
	rpgSystem.AppendAbilities("strength")
	dice := mock.NewRollMock("d10", rpg.D10HM)
	person1 := NewCreature("test-ability-d10hm-1", rpgSystem)
	err := person1.AddAbility(Ability{Name: "strength", Value: 5})
	if err != nil {
		t.Errorf("error on ability %v", err)
	}
	result, err := person1.AbilityCheck(dice, Check{Ability: "strength", Target: 1}, logger)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if result.Success != true {
		t.Errorf("result.Success %v", result.Success)
	}
}
