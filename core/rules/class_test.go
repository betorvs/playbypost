package rules

import (
	"testing"
)

func TestClassAttackBase(t *testing.T) {
	t.Log("Testing Class Attack Bonus Calculation")
	{
		class1 := Class{
			AttackBonus: Good,
		}
		class2 := Class{
			AttackBonus: Average,
		}
		class3 := Class{
			AttackBonus: Poor,
		}
		var res1, res2, res3 int
		for i := 1; i < 21; i++ {
			class1.Level = i
			class2.Level = i
			class3.Level = i
			// t.Log("level", i, "attack base", "good", class1.AttackBase(), "average", class2.AttackBase(), "poor", class3.AttackBase())
			res1 = class1.AttackBase()
			res2 = class2.AttackBase()
			res3 = class3.AttackBase()
		}
		if res1 != 20 {
			t.Errorf("\t Not expected attack bonus %d", res1)
		}
		if res2 != 15 {
			t.Errorf("\t Not expected attack bonus %d", res2)
		}
		if res3 != 10 {
			t.Errorf("\t Not expected attack bonus %d", res3)
		}
	}
}

func TestCalcXP(t *testing.T) {
	t.Log("Testing XP Calculation")
	{
		class1 := Class{
			Level: 9,
		}
		level10 := class1.calcNextLevelXP()
		if level10 != 45000 {
			t.Errorf("\t Not expected xp %d", level10)
		}
		class2 := Class{
			Level: 14,
		}
		level15 := class2.calcNextLevelXP()
		if level15 != 105000 {
			t.Errorf("\t Not expected xp %d", level15)
		}
		class3 := Class{
			Level: 19,
		}
		level20 := class3.calcNextLevelXP()
		if level20 != 190000 {
			t.Errorf("\t Not expected xp %d", level20)
		}
	}
}
