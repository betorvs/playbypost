package d10hm

import (
	"log/slog"
	"os"
	"testing"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/library"
	"github.com/betorvs/playbypost/core/tests/mock"
)

func TestAbilityD10HM(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rpgSystem := rpg.LoadRPGSystemsDefault(rpg.D10HM)
	lib := library.New()
	lib.AppendAbilities("strength")
	dice := mock.NewRollMock("d10", rpg.D10HM)
	person1 := New("test-ability-d10hm-1", rpgSystem)
	err := person1.AddAbility(base.Ability{Name: "strength", Value: 5}, lib)
	if err != nil {
		t.Errorf("error on ability %v", err)
	}
	result, err := person1.AbilityCheck(dice, rules.Check{Ability: "strength", Target: 1}, logger, lib)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if result.Success != true {
		t.Errorf("result.Success %v", result.Success)
	}
}
