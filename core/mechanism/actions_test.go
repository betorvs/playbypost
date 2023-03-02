package mechanism

import (
	"testing"

	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/types"
)

// MockRollInternal struct
type MockRollInternal struct {
}

func (r MockRollInternal) DiceRoll(name, text string) (int, string, error) {
	switch name {
	case "pc1":
		return 15, "Rolled 15", nil
	case "npc1":
		return 5, "Rolled 5", nil
	default:
		return 1, "Rolled 1", nil
	}
}

func TestActions(t *testing.T) {
	longSword := rules.Weapon{
		Name:         "longsword",
		DamageMedium: "4",
		DamageType:   rules.Slashing,
	}
	pc1Attack := rules.AttackOption{
		AttackBonus:  1,
		Damage:       "1d8+2",
		Type:         rules.Melee,
		AttackerSize: rules.Medium,
		Weapon:       &longSword,
	}
	fighter := rules.Class{
		Level:       1,
		AttackBonus: rules.Good,
	}
	player1 := &rules.Player{
		Name:                  "pc1",
		Level:                 1,
		ArmorClass:            13,
		PreferenceAttackIndex: 0,
		HitPoints:             10,
		Class:                 []rules.Class{fighter},
		AttackOption: []rules.AttackOption{
			pc1Attack,
		},
	}
	pc1 := rules.NewPCCombatParticipant(player1)
	morningStar := rules.Weapon{
		Name:         "morning-start",
		DamageMedium: "3",
		DamageType:   rules.Bludgeoning,
	}
	npc1Attack := rules.AttackOption{
		AttackBonus:  2,
		Damage:       "1d6",
		AttackerSize: rules.Small,
		Weapon:       &morningStar,
		Type:         rules.Melee,
	}
	monster1 := &rules.Monster{
		Name:                  "npc1",
		Title:                 "Goblin",
		ArmorClass:            15,
		HitPoints:             5,
		PreferenceAttackIndex: 0,
		AttackOption: []rules.AttackOption{
			npc1Attack,
		},
	}
	// "npc1"
	npc1 := rules.NewNPCCombatParticipant(monster1)
	// test attack
	diceTest := &MockRollInternal{}
	t.Log("Given a encounter with 2 participants pc1 and npc1")
	{
		testCombat := rules.NewCombatSingle(pc1, npc1)
		testAttack := NewCombatCommand(testCombat, types.DoAttack, 0, diceTest)
		testAttack.Call()
		if npc1.ActorNPC.HitPoints == 10 {
			t.Errorf("\t Not expected npc1 hit points %d", npc1.ActorNPC.HitPoints)
		}
		t.Log("Counter attack happen")
		testCombat2 := rules.NewCombatSingle(npc1, pc1)
		testAttack2 := NewCombatCommand(testCombat2, types.DoAttack, 0, diceTest)
		testAttack2.Call()
		t.Log("Ops, npc1 should take a potion not attack, calling Undo")
		testAttack2.Undo()
		if pc1.ActorPC.HitPoints != 10 {
			t.Errorf("\t Not expected pc1 hit points %d", pc1.ActorPC.HitPoints)
		}
		t.Log("PC hit points restored", pc1.ActorPC.HitPoints)
		//
		t.Log("NPC did not take that action, instead will use total defense")
		// npc1.ActorNPC.SetTotalDefense()
		testCombat2a := rules.NewCombatSingle(npc1, pc1)
		testAttack2a := NewCombatCommand(testCombat2a, types.DoTotalDefense, 0, diceTest)
		testAttack2a.Call()
		// npc1 should be in total defense now
		testCombat3 := rules.NewCombatSingle(pc1, npc1)
		testAttack3 := NewCombatCommand(testCombat3, types.DoAttack, 0, diceTest)
		testAttack3.Call()
		t.Log("Rolled 15, but pc1 cannot hit")
		if !npc1.ActorNPC.GetTotalDefense() {
			t.Errorf("\t Not expected npc1 without total defense")
		}
	}

}
