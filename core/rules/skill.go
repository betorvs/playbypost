package rules

import (
	"errors"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
)

func (c *Creature) AddSkill(s Skill) error {
	if slices.Contains(c.RPG.Skill.List, s.Name) {
		c.Skills[s.Name] = s
		return nil
	}
	return errors.New(SkillInvalid)
}

func (c *Creature) SkillCheck(d *rpg.Roll, check Check) (Result, error) {
	response := Result{}
	if !slices.Contains(c.RPG.Skill.List, c.Skills[check.Skill].Name) {
		return response, errors.New(SkillInvalid)
	}
	abilityBase := c.Skills[check.Skill].Base
	if check.Override != "" && slices.Contains(c.RPG.Ability.List, check.Override) {
		abilityBase = check.Override
	}
	switch c.RPG.SuccessRule {
	case rpg.GreaterThan:
		result, res, err := d.Check("check skill")
		if err != nil {
			return response, err
		}
		bonus := c.calcSkillModifier(check.Skill)
		c.RPG.Logger.Info("dice result", "result", result+bonus, "rolled", res)
		abilityBonus := c.calcAbilityModifier(abilityBase)
		response.Success = result+abilityBonus >= check.Target
		response.Text = res
		response.Result = result + abilityBonus

	case rpg.CountResults:
		abilityBonus := c.Abilities[abilityBase].Value
		calcDices := d.Dice(c.Skills[check.Skill].Value+abilityBonus, 0)
		result, res, err := d.FreeRoll("check abiliy", calcDices)
		if err != nil {
			return response, err
		}
		c.RPG.Logger.Info("dice result", "result", result, "rolled", res)
		response.Success = result >= check.Target
		response.Text = res
		response.Result = result

	case rpg.DifficultAndCount:
		abilityBonus := c.Abilities[abilityBase].Value
		calcDices := d.Dice(c.Skills[check.Skill].Value+abilityBonus, check.Target)
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

// D20 only
func (c *Creature) calcSkillModifier(skill string) int {
	switch c.RPG.SkillRank {
	case rpg.SkillRanks:
		return c.Skills[skill].Value
	case rpg.SkillModifiers:
		if c.Skills[skill].Value == 1 {
			return 8
		}
		return 4
	case rpg.LevelBasedSkill:
		if c.Skills[skill].Value == 1 {
			value, _ := c.Extension.SkillBonus(skill)
			return value
		}
		return 0
	case rpg.OnePerOne:
		return c.Skills[skill].Value
	}
	return 0
}
