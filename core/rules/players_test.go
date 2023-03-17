package rules

import (
	"testing"
)

func TestPlayerAttackBase(t *testing.T) {
	t.Log("Testing Player Attack Base without any modifiers")
	{
		fighter := Class{
			Level:       10,
			AttackBonus: Good,
		}
		mage := Class{
			Level:       10,
			AttackBonus: Poor,
		}
		bard := Class{
			Level:       10,
			AttackBonus: Average,
		}
		multiclassplayer1 := NewPlayer("1", "1", "human", true, []Class{fighter, mage})
		// multiclassplayer1 := Player{
		// 	Name:       "",
		// 	PlayerName: "",
		// 	MultiClass: true,
		// 	Class:      []Class{fighter, mage},
		// 	Level:      0,
		// 	Race:       Race{Size: Medium},
		// 	Aligment:   0,
		// 	Metadata:   Metadata{},
		// 	Creature:   {Ability{Strength: 10}},
		// 	Spells:     map[int][]string{},
		// 	Languages:  []string{},
		// 	Gear:       Gear{},
		// }
		res1 := multiclassplayer1.CalcAttackTotal(Melee)
		// fmt.Println(res1)
		if res1 != 15 {
			t.Errorf("\t Not expected attack bonus %d", res1)
		}
		multiclassplayer2 := NewPlayer("2", "2", "human", true, []Class{fighter, bard})
		// multiclassplayer2 := Player{
		// 	Ability{
		// 		Strength: 10,
		// 	},
		// 	Race: Race{
		// 		Size: Medium,
		// 	},
		// 	MultiClass: true,
		// 	Class: []Class{
		// 		fighter,
		// 		bard,
		// 	},
		// }
		res2 := multiclassplayer2.CalcAttackTotal(Melee)
		// fmt.Println(res2)
		if res2 != 17 {
			t.Errorf("\t Not expected attack bonus %d", res2)
		}
	}
}
