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
		result, err := d.Check("check abiliy")
		if err != nil {
			return response, err
		}
		bonus := c.calcAbilityModifier(check.Ability)
		c.RPG.Logger.Info("dice result", "result", result.Result+bonus, "rolled", result.Rolled)
		response.Success = result.Result+bonus >= check.Target
		if result.Result == 20 {
			c.RPG.Logger.Info("20 natural roll")
			response.Success = true
		}
		response.Description = result.Description
		response.Result = result.Result + bonus

	case rpg.CountResults:
		calcDices := d.FormatDice(c.Abilities[check.Ability].Value, 0)
		result, err := d.FreeRoll("check abiliy", calcDices)
		if err != nil {
			return response, err
		}
		c.RPG.Logger.Info("dice result", "result", result, "rolled", result.Description)
		response.Success = result.Result >= check.Target
		response.Description = result.Description
		response.Result = result.Result

	case rpg.DifficultAndCount:
		calcDices := d.FormatDice(c.Abilities[check.Ability].Value, check.Target)
		result, err := d.FreeRoll("check abiliy", calcDices)
		if err != nil {
			return response, err
		}
		c.RPG.Logger.Info("dice result", "result", result.Result, "rolled", result.Description)
		response.Success = result.Result > check.Difficult
		response.Description = result.Description
		response.Result = result.Result
	}
	return response, nil
}
