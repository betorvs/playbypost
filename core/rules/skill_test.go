//go:build unit

package rules

import (
	"log/slog"
	"os"
	"testing"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/tests/mock"
)

func TestSkillD10HM(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rpgSystem := rpg.LoadRPGSystemsDefault(rpg.D10HM)
	rpgSystem.AppendAbilities("dexterity")
	rpgSystem.AppendSkills("athletics")
	dice := mock.NewRollMock("d10", rpg.D10HM)
	person1 := NewCreature("test-athletics-d10hm-1", rpgSystem)
	err := person1.AddAbility(Ability{Name: "dexterity", Value: 5})
	if err != nil {
		t.Errorf("error on ability %v", err)
	}
	err = person1.AddSkill(Skill{Name: "athletics", Value: 5, Base: "dexterity"})
	if err != nil {
		t.Errorf("error on skill %v", err)
	}
	result, err := person1.SkillCheck(dice, Check{Skill: "athletics", Target: 5}, logger)
	if err != nil {
		t.Errorf("error %v", err)
	}
	if result.Success != true {
		t.Errorf("result.Success %v", result.Success)
	}

}
