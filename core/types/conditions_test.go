package types

import "testing"

func TestConditions(t *testing.T) {
	t.Log("Checking all conditions in loop")
	for i := Alive; i <= Unconscious; i++ {
		res := i.String()
		if res == "" && res != Unknown {
			t.Errorf("\t Condition should not be empty %s", res)
		}
	}
}
