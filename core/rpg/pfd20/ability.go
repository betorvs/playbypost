package pfd20

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/library"
)

func (c *PathfinderCharacter) calcAbilityModifier(ability string) int {
	if c.RPG.BaseSystem == rpg.D20 {
		result := math.Floor((float64(c.Abilities[ability].Value) - 10) / 2)
		return int(result)
	}
	return c.Abilities[ability].Value
}

func (c *PathfinderCharacter) AbilityCheck(d rpg.RollInterface, check rules.Check, logger *slog.Logger, lib *library.Library) (rules.Result, error) {
	response := rules.Result{}
	if !slices.Contains(lib.Ability.List, c.Abilities[check.Ability].Name) {
		return response, errors.New(base.AbilityInvalid)
	}
	diceName := fmt.Sprintf("check-ability-%s-%s", check.Ability, c.Name())
	switch c.RPG.SuccessRule {
	case rpg.GreaterThan:
		result, err := d.Check(diceName)
		if err != nil {
			return response, err
		}
		bonus := c.calcAbilityModifier(check.Ability)
		logger.Info("dice result", "result", result.Result+bonus, "rolled", result.Rolled)
		response.Success = result.Result+bonus >= check.Target
		if result.Result == 20 {
			logger.Info("20 natural roll")
			response.Success = true
		}
		response.Description = result.Description
		response.Result = result.Result + bonus
	}
	return response, nil
}
