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

func (c *StorytellingCharacter) AbilityCheck(d rpg.RollInterface, check rules.Check, logger *slog.Logger, lib *library.Library) (rules.Result, error) {
	response := rules.Result{}
	if !slices.Contains(lib.Ability.List, c.Abilities[check.Ability].Name) {
		return response, errors.New(base.AbilityInvalid)
	}
	diceName := fmt.Sprintf("check-ability-%s-%s", check.Ability, c.Name())
	switch c.RPG.SuccessRule {

	case rpg.CountResults:
		calcDices := d.FormatDice(c.Abilities[check.Ability].Value, 0)
		result, err := d.FreeRoll(diceName, calcDices)
		if err != nil {
			return response, err
		}
		logger.Info("dice result", "result", result, "rolled", result.Description)
		response.Success = result.Result >= check.Target
		response.Description = result.Description
		response.Result = result.Result

	}
	return response, nil
}
