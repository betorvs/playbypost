/*
Package diceroll implements a simple call to simulate a dice roll present in role playing games
*/
package diceroll

import (
	"fmt"

	"github.com/betorvs/dice"
)

// DiceRoller interface
type DiceRoller interface {
	DiceRoll(name, text string) (int, string, error)
}

// RollInternal struct
type RollInternal struct {
}

// DiceRoll func returns a int value, a string text and error
func (r RollInternal) DiceRoll(name, text string) (int, string, error) {
	diceRolled, _, err := dice.Roll(text)
	if err != nil {
		return 0, "No dices to roll", err
	} else {
		message := fmt.Sprintf("Hey %s, your dice rolled %s and result %v with rolls %s", name, diceRolled.Description(), diceRolled.Int(), diceRolled.String())
		return diceRolled.Int(), message, nil
	}
}
