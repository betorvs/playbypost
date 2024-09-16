//go:build unit

package rules

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/tests/mock"
)

var (
	rawData1 = []byte(`{
  "size": 5,
  "armor": 3,
  "health": 8,
  "weapon": {
    "longsword": {
      "name": "longsword",
      "value": 6,
      "description": "magical longsword"
    }
  },
  "defense": 2,
  "willpower": 4,
  "initiative": 5
}`)
	rawData2 = []byte(`{
	"size": 5,
	"armor": 3,
	"health": 7,
	"weapon": {
	  "longsword": {
		"name": "longsword",
		"value": 3,
		"description": "a longsword"
	  }
	},
	"defense": 2,
	"willpower": 4,
	"initiative": 5
  }`)
)

func TestCombatD10HM(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	rpgSystem := rpg.LoadRPGSystemsDefault(rpg.D10HM)
	rpgSystem.AppendAbilities("strength")
	rpgSystem.AppendSkills("weaponry")
	dice := mock.NewRollMock("d10", rpg.D10HM)
	person1 := NewCreature("test-combat-p1-d10hm-1", rpgSystem)
	err := person1.AddAbility(Ability{Name: "strength", Value: 5})
	if err != nil {
		t.Errorf("error add ability %v", err)
	}
	err = person1.AddSkill(Skill{Name: "weaponry", Value: 5, Base: "strength"})
	if err != nil {
		t.Errorf("error add skill %v", err)
	}
	extended := rpg.NewExtended()
	_ = json.Unmarshal(rawData1, &extended)
	person1.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
	evilperson1 := NewCreature("test-combat-ep1-d10hm-1", rpgSystem)
	err = evilperson1.AddAbility(Ability{Name: "strength", Value: 5})
	if err != nil {
		t.Errorf("error add ability %v", err)
	}
	err = evilperson1.AddSkill(Skill{Name: "weaponry", Value: 5, Base: "strength"})
	if err != nil {
		t.Errorf("error add skill %v", err)
	}
	extended2 := rpg.NewExtended()
	_ = json.Unmarshal(rawData2, &extended2)
	evilperson1.Extension = rpg.NewExtendedSystem(rpgSystem, extended2)
	attack := NewAttack("test-combat-attack-1", "longsword", Melee, person1, evilperson1, dice, logger)
	attack.Call()
	if attack.Response.Success != true {
		t.Errorf("attack.Response.Success %v", attack.Response.Success)
	}
	if attack.Defensor.IsDead() != true {
		t.Errorf("attack.Defensor.IsDead() %v", attack.Defensor.IsDead())
	}

}
