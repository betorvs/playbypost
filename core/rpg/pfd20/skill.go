package pfd20

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

func (c *PathfinderCharacter) SkillCheck(d rpg.RollInterface, check rules.Check, logger *slog.Logger, lib *library.Library) (rules.Result, error) {
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
	case rpg.GreaterThan:
		result, err := d.Check(diceName)
		if err != nil {
			return response, err
		}
		bonus := c.calcSkillModifier(check.Skill)
		logger.Info("dice result", "result", result.Result+bonus, "rolled", result.Rolled)
		abilityBonus := c.calcAbilityModifier(abilityBase)
		response.Success = result.Result+abilityBonus >= check.Target
		response.Description = result.Description
		response.Result = result.Result + abilityBonus
	}
	return response, nil
}

func (c *PathfinderCharacter) calcSkillModifier(skill string) int {
	v, ok := c.Proficiency[skill]
	if !ok {
		return 0
	}
	return proficiencyRank(v.Level, c.Level)
}
