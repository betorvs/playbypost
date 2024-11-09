package rpg

import (
	"fmt"

	"github.com/betorvs/dice"
)

type RollInterface interface {
	FreeRoll(name, text string) (DiceRoll, error)
	Check(name string) (DiceRoll, error)
	FormatDice(m, target int) string
}

// Roll struct
type Roll struct {
	RPGSystem *RPGSystem
}

type DiceRoll struct {
	RequestedBy string
	Description string
	Result      int
	Rolled      string
}

func NewRollMock(rpgSystem *RPGSystem) RollInterface {
	return &Roll{RPGSystem: rpgSystem}
}

// FreeRoll func returns a int value, a string text and error
func (r Roll) FreeRoll(name, text string) (DiceRoll, error) {
	res := DiceRoll{}
	res.RequestedBy = name
	diceRolled, _, err := dice.Roll(text)
	if err != nil {
		res.Result = 0
		res.Description = "No dices to roll"
		return res, err
	} else {
		// message := fmt.Sprintf("%s rolled %s and result %v with rolls %s", name, diceRolled.Description(), diceRolled.Int(), diceRolled.String())
		res.Result = diceRolled.Int()
		res.Description = diceRolled.Description()
		res.Rolled = diceRolled.String()
		return res, nil
	}
}

func (r Roll) Check(name string) (DiceRoll, error) {
	res := DiceRoll{}
	diceRolled, _, err := dice.Roll(r.RPGSystem.BaseDice)
	if err != nil {
		res.Result = 0
		res.Description = "No dices to roll"
		return res, err
	} else {
		// message := fmt.Sprintf("%s rolled %s and result %v with rolls %s", name, diceRolled.Description(), diceRolled.Int(), diceRolled.String())
		res.Result = diceRolled.Int()
		res.Description = diceRolled.Description()
		res.Rolled = diceRolled.String()
		return res, nil
	}
}

func (r Roll) FormatDice(m, target int) string {
	dices := 1
	if m > 0 {
		dices = m
	}
	var dice string
	switch r.RPGSystem.Name {
	case D10HM:
		dice = fmt.Sprintf("%d%srv8", dices, r.RPGSystem.BaseDice)
	}
	return dice
}
