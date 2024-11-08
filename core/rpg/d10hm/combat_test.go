package d10hm

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
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
	person1 := New("test-combat-p1-d10hm-1", rpgSystem)
	err := person1.AddAbility(base.Ability{Name: "strength", Value: 5})
	if err != nil {
		t.Errorf("error add ability %v", err)
	}
	err = person1.AddSkill(base.Skill{Name: "weaponry", Value: 5, Base: "strength"})
	if err != nil {
		t.Errorf("error add skill %v", err)
	}
	// extended := rpg.NewExtended()
	_ = json.Unmarshal(rawData1, &person1.D10Extented)
	// person1.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
	evilperson1 := New("test-combat-ep1-d10hm-1", rpgSystem)
	err = evilperson1.AddAbility(base.Ability{Name: "strength", Value: 5})
	if err != nil {
		t.Errorf("error add ability %v", err)
	}
	err = evilperson1.AddSkill(base.Skill{Name: "weaponry", Value: 5, Base: "strength"})
	if err != nil {
		t.Errorf("error add skill %v", err)
	}
	// extended2 := rpg.NewExtended()
	_ = json.Unmarshal(rawData2, &evilperson1.D10Extented)
	// evilperson1.Extension = rpg.NewExtendedSystem(rpgSystem, extended2)
	attack := rules.NewAttack("test-combat-attack-1", "longsword", rules.Melee, person1, evilperson1, dice, logger)
	attack.Call()
	if attack.Response.Success != true {
		t.Errorf("attack.Response.Success %v", attack.Response.Success)
	}
	if attack.Defensor.IsDead() != true {
		t.Errorf("attack.Defensor.IsDead() %v", attack.Defensor.IsDead())
	}

}
