package d10hm

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
)

func (c *StorytellingCharacter) SkillCheck(d rpg.RollInterface, check rules.Check, logger *slog.Logger) (rules.Result, error) {
	response := rules.Result{}
	if !slices.Contains(c.RPG.Skill.List, c.Skills[check.Skill].Name) {
		return response, errors.New(base.SkillInvalid)
	}
	abilityBase := c.Skills[check.Skill].Base
	if check.Override != "" && slices.Contains(c.RPG.Ability.List, check.Override) {
		abilityBase = check.Override
	}
	diceName := fmt.Sprintf("check-skill-%s-%s", check.Skill, c.Name())
	switch c.RPG.SuccessRule {

	case rpg.CountResults:
		abilityBonus := c.Abilities[abilityBase].Value
		calcDices := d.FormatDice(c.Skills[check.Skill].Value+abilityBonus, 0)
		result, err := d.FreeRoll(diceName, calcDices)
		if err != nil {
			return response, err
		}
		logger.Info("dice result", "result", result, "rolled", result.Rolled)
		response.Success = result.Result >= check.Target
		response.Description = result.Description
		response.Result = result.Result
		response.Rolled = result.Rolled

	}
	return response, nil
}
