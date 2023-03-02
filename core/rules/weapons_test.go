package rules

import "testing"

func TestIsCritical(t *testing.T) {
	listWeaponsShortCritical := []Weapon{
		{
			Name:     "Hammer, gnome hooked",
			Critical: `x3/x4`,
		},
		{
			Name:     "warhammer",
			Critical: "x3",
		},
		{
			Name:     "Pick, heavy",
			Critical: "x4",
		},
	}
	shortDicesRolls := []int{20, 19}
	for _, w := range listWeaponsShortCritical {
		for _, v := range shortDicesRolls {
			valid, _ := w.IsCritical(v, false)
			if v == 20 && !valid {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid)
			}
			if v == 19 && valid {
				t.Errorf("\t Weapon %s Should NOT be critical %d but %v", w.Name, v, valid)
			}
			valid1, _ := w.IsCritical(v, true)
			if v == 20 && !valid1 {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid1)
			}
			if v == 19 && !valid1 {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid1)
			}
		}
	}
	listWeaponsLargerCritical := []Weapon{
		{
			Name:     "Longsword",
			Critical: "19-20/x2",
		},
		{
			Name:     "Scimitar",
			Critical: "18-20/x2",
		},
	}
	largerDicesRolls := []int{20, 17, 15}
	for _, w := range listWeaponsLargerCritical {
		for _, v := range largerDicesRolls {
			valid, _ := w.IsCritical(v, false)
			if v == 20 && !valid {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid)
			}
			if v == 17 && valid {
				t.Errorf("\t Weapon %s Should NOT be critical %d but %v", w.Name, v, valid)
			}
			if v == 15 && valid {
				t.Errorf("\t Weapon %s Should NOT be critical %d but %v", w.Name, v, valid)
			}
			valid1, _ := w.IsCritical(v, true)
			if v == 20 && !valid1 {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid1)
			}
			if v == 17 && !valid1 {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid1)
			}
			if v == 15 && !valid1 && w.Name == "Scimitar" {
				t.Errorf("\t Weapon %s Should be critical %d but %v", w.Name, v, valid1)
			}
		}
	}
}
