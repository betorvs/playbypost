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

		multiclassplayer1 := Player{
			Abilities: Ability{
				Strength: 10,
			},
			Race: Race{
				Size: Medium,
			},
			MultiClass: true,
			Class: []Class{
				fighter,
				mage,
			},
		}
		res1 := multiclassplayer1.CalcAttackTotal(Melee)
		// fmt.Println(res1)
		if res1 != 15 {
			t.Errorf("\t Not expected attack bonus %d", res1)
		}

		multiclassplayer2 := Player{
			Abilities: Ability{
				Strength: 10,
			},
			Race: Race{
				Size: Medium,
			},
			MultiClass: true,
			Class: []Class{
				fighter,
				bard,
			},
		}
		res2 := multiclassplayer2.CalcAttackTotal(Melee)
		// fmt.Println(res2)
		if res2 != 17 {
			t.Errorf("\t Not expected attack bonus %d", res2)
		}
	}
}
