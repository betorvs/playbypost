package d10hm

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/library"
)

func (c *StorytellingCharacter) SkillCheck(d rpg.RollInterface, check rules.Check, logger *slog.Logger, lib *library.Library) (rules.Result, error) {
	response := rules.Result{}
	if !slices.Contains(lib.Skill.List, c.Skills[check.Skill].Name) {
		return response, errors.New(base.SkillInvalid)
	}
	abilityBase := c.Skills[check.Skill].Base
	if check.Override != "" && slices.Contains(lib.Ability.List, check.Override) {
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
