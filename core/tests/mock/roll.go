package mock

import (
	"fmt"

	"github.com/betorvs/playbypost/core/rpg"
)

const (
	D10HM string = "D10HomeMade"  // D10 based on World of Darkness
	D10OS string = "D10OldSchool" // D10 based on Old School
)

// MockRoll struct
type MockRoll struct {
	RPGSystem RPGSystem
}

type RPGSystem struct {
	Name     string
	BaseDice string
}

func NewRollMock(d, n string) rpg.RollInterface {
	rpgSystem := RPGSystem{Name: n, BaseDice: d}
	return &MockRoll{RPGSystem: rpgSystem}
}

// FreeRoll func returns a int value, a string text and error
func (r MockRoll) FreeRoll(name, text string) (rpg.DiceRoll, error) {
	fmt.Println(name)
	res := rpg.DiceRoll{}
	res.RequestedBy = name
	res.Result = 0
	res.Description = ""
	res.Rolled = "0 [0]"
	switch name {
	case "check-ability-strength-test-ability-d10hm-1":
		res.Result = 2
		res.Description = "2"
		res.Rolled = "2 [9 9 3 3 3]"
	case "check-skill-athletics-test-athletics-d10hm-1":
		res.Result = 5
		res.Description = "5"
		res.Rolled = "5 [9 9 9 9 9 3 3 3 3 3]"
	case "attack-roll-test-combat-p1-d10hm-1-strenght":
		res.Result = 10
		res.Description = "10"
		res.Rolled = "10 [9 9 9 9 9 10 10 9 9 9 3 3]"
	}
	return res, nil
}

func (r MockRoll) Check(name string) (rpg.DiceRoll, error) {
	res := rpg.DiceRoll{}

	res.Result = 0
	res.Description = ""
	res.Rolled = "[0]"
	return res, nil
}

func (r MockRoll) FormatDice(m, target int) string {
	dices := 1
	if m > 0 {
		dices = m
	}
	var dice string
	switch r.RPGSystem.Name {
	case D10HM:
		dice = fmt.Sprintf("%d%srv8", dices, r.RPGSystem.BaseDice)
	case D10OS:
		dice = fmt.Sprintf("%d%srv%d", dices, r.RPGSystem.BaseDice, target)
	}
	return dice
}
