package types

import "testing"

func TestActions(t *testing.T) {
	t.Log("Checking all actions in loop")
	for i := DoAttack; i <= DoTotalDefense; i++ {
		res := i.String()
		if res == "" && res != Unknown {
			t.Errorf("\t Actions should not be empty %s", res)
		}
	}
}

func TestEffect(t *testing.T) {
	t.Log("Checking all effects in loop")
	for i := NoneEffect; i <= ChangeConditionEffect; i++ {
		res := i.String()
		if res == "" && res != Unknown {
			t.Errorf("\t Effect should not be empty %s", res)
		}
	}
}
