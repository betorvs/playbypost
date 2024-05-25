package rules

import (
	"errors"
	"math"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
)

func (c *Creature) calcAbilityModifier(ability string) int {
	if c.RPG.BaseSystem == rpg.D20 {
		result := math.Floor((float64(c.Abilities[ability].Value) - 10) / 2)
		return int(result)
	}
	return c.Abilities[ability].Value
}

func (c *Creature) AddAbility(a Ability) error {
	if slices.Contains(c.RPG.Ability.List, a.Name) {
		c.Abilities[a.Name] = a
		return nil
	}
	return errors.New(AbilityInvalid)
}

func (c *Creature) AbilityCheck(d *rpg.Roll, check Check) (Result, error) {
	response := Result{}
	if !slices.Contains(c.RPG.Ability.List, c.Abilities[check.Ability].Name) {
		return response, errors.New(AbilityInvalid)
	}
	switch c.RPG.SuccessRule {
	case rpg.GreaterThan:
		result, res, err := d.Check("check abiliy")
		if err != nil {
			return response, err
		}
		bonus := c.calcAbilityModifier(check.Ability)
		c.RPG.Logger.Info("dice result", "result", result+bonus, "rolled", res)
		response.Success = result+bonus >= check.Target
		if result == 20 {
			c.RPG.Logger.Info("20 natural roll")
			response.Success = true
		}
		response.Text = res
		response.Result = result + bonus

	case rpg.CountResults:
		calcDices := d.Dice(c.Abilities[check.Ability].Value, 0)
		result, res, err := d.FreeRoll("check abiliy", calcDices)
		if err != nil {
			return response, err
		}
		c.RPG.Logger.Info("dice result", "result", result, "rolled", res)
		response.Success = result >= check.Target
		response.Text = res
		response.Result = result

	case rpg.DifficultAndCount:
		calcDices := d.Dice(c.Abilities[check.Ability].Value, check.Target)
		result, res, err := d.FreeRoll("check abiliy", calcDices)
		if err != nil {
			return response, err
		}
		c.RPG.Logger.Info("dice result", "result", result, "rolled", res)
		response.Success = result > check.Difficult
		response.Text = res
		response.Result = result
	}
	return response, nil
}
