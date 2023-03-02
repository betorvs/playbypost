package rules

import (
	"testing"

	"github.com/betorvs/playbypost/core/types"
)

func TestBonuses(t *testing.T) {
	t.Log("Checking all Bonus types in loop")
	{
		for i := EnhancementBonus; i <= CircumstanceBonus; i++ {
			res := i.String()
			if res == "" && res != types.Unknown {
				t.Errorf("\t Bonus types should not be empty %s", res)
			}
		}
	}
	t.Log("Check TemporaryCombatBonuses methods")
	{
		t.Log("Check factory function")
		temp1 := NewTemporaryCombatBonuses()
		if temp1 == nil {
			t.Fatalf("\t Factory function should create a pointer")
		}
		// t.Log("Example", temp1.Source)
		t.Log("Check AddBonuses method")
		t1 := 3
		temp1.AddBonuses(EnhancementBonus, AttackBonus, t1, NotApplicableTTL, "magical-item-1")
		if t1 != temp1.GetBonus(AttackBonus) {
			// t.Log("Example", temp1.Source)
			t.Fatalf("\t AddBonuses did not update total value in TemporaryCombatBonuses about attack with magical-item-1 %d", temp1.GetBonus(AttackBonus))
		}
		t.Log("Check AddBonuses method with same type not stackable")
		t2 := 2
		temp1.AddBonuses(EnhancementBonus, AttackBonus, t2, NotApplicableTTL, "magical-item-2")
		if t1 != temp1.GetBonus(AttackBonus) {
			t.Fatalf("\t AddBonuses did not update total value in TemporaryCombatBonuse about attack with magical-item-2 %d", temp1.GetBonus(AttackBonus))
		}
		t.Log("Total attack bonus is", temp1.GetBonus(AttackBonus))
		t.Log("Check AddBonuses method with stackable bonus")
		temp1.AddBonuses(DodgeBonus, DefenseBonus, t1, NotApplicableTTL, "class-feature-1")
		if t1 != temp1.GetBonus(DefenseBonus) {
			t.Fatalf("\t AddBonuses did not update total value in TemporaryCombatBonuses about Defense class-feature-1")
		}
		temp1.AddBonuses(DodgeBonus, DefenseBonus, t2, NotApplicableTTL, "magical-item-3")
		dodgeTotal := t1 + t2
		if dodgeTotal != temp1.GetBonus(DefenseBonus) {
			t.Fatalf("\t AddBonuses did not update total value in TemporaryCombatBonuses about Defense magical-item-3")
		}
	}
	t.Log("Check TemporaryCombatBonuses methods related to stackable bonus")
	{
		t.Log("Check factory function again")
		temp2 := NewTemporaryCombatBonuses()
		if temp2 == nil {
			t.Fatalf("\t Factory function should create a pointer")
		}
		t.Log("Add 3 different Attack bonuses")
		a1 := 2
		temp2.AddBonuses(EnhancementBonus, AttackBonus, a1, NotApplicableTTL, "magical-sword-1")
		a2 := 3
		temp2.AddBonuses(LuckBonus, AttackBonus, a2, NotApplicableTTL, "magical-ring-1")
		a3 := 1
		temp2.AddBonuses(MoraleBonus, AttackBonus, a3, NotApplicableTTL, "magical-belt-1")
		aTotal := a1 + a2 + a3
		if aTotal != temp2.GetBonus(AttackBonus) {
			// t.Log("Example", temp2.Source)
			t.Fatalf("\t Total attack bonus should be %d not %d", aTotal, temp2.GetBonus(AttackBonus))
		}
		t.Log("Add 3 different defense bonuses")
		temp2.AddBonuses(EnhancementBonus, DefenseBonus, a1, NotApplicableTTL, "magical-armor-1")
		temp2.AddBonuses(DodgeBonus, DefenseBonus, a2, NotApplicableTTL, "magical-ring-2")
		temp2.AddBonuses(DodgeBonus, DefenseBonus, a3, NotApplicableTTL, "magical-belt-2")
		if aTotal != temp2.GetBonus(DefenseBonus) {
			// t.Log("Example", temp2.Source)
			t.Fatalf("\t Total defens bonus should be %d not %d", aTotal, temp2.GetBonus(DefenseBonus))
		}
	}
}
