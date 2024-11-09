package rules

// func (c *Creature) AddSkill(s Skill) error {
// 	if slices.Contains(c.RPG.Skill.List, s.Name) {
// 		c.Skills[s.Name] = s
// 		return nil
// 	}
// 	return errors.New(SkillInvalid)
// }

// func (c *Creature) SkillCheck(d rpg.RollInterface, check Check, logger *slog.Logger) (Result, error) {
// 	response := Result{}
// 	if !slices.Contains(c.RPG.Skill.List, c.Skills[check.Skill].Name) {
// 		return response, errors.New(SkillInvalid)
// 	}
// 	abilityBase := c.Skills[check.Skill].Base
// 	if check.Override != "" && slices.Contains(c.RPG.Ability.List, check.Override) {
// 		abilityBase = check.Override
// 	}
// 	diceName := fmt.Sprintf("check-skill-%s-%s", check.Skill, c.Name)
// 	switch c.RPG.SuccessRule {
// 	case rpg.GreaterThan:
// 		result, err := d.Check(diceName)
// 		if err != nil {
// 			return response, err
// 		}
// 		bonus := c.calcSkillModifier(check.Skill)
// 		logger.Info("dice result", "result", result.Result+bonus, "rolled", result.Rolled)
// 		abilityBonus := c.calcAbilityModifier(abilityBase)
// 		response.Success = result.Result+abilityBonus >= check.Target
// 		response.Description = result.Description
// 		response.Result = result.Result + abilityBonus
// 		// response.Rolled = result.Rolled

// 	case rpg.CountResults:
// 		abilityBonus := c.Abilities[abilityBase].Value
// 		calcDices := d.FormatDice(c.Skills[check.Skill].Value+abilityBonus, 0)
// 		result, err := d.FreeRoll(diceName, calcDices)
// 		if err != nil {
// 			return response, err
// 		}
// 		logger.Info("dice result", "result", result, "rolled", result.Rolled)
// 		response.Success = result.Result >= check.Target
// 		response.Description = result.Description
// 		response.Result = result.Result
// 		response.Rolled = result.Rolled

// 	case rpg.DifficultAndCount:
// 		abilityBonus := c.Abilities[abilityBase].Value
// 		calcDices := d.FormatDice(c.Skills[check.Skill].Value+abilityBonus, check.Target)
// 		result, err := d.FreeRoll(diceName, calcDices)
// 		if err != nil {
// 			return response, err
// 		}
// 		logger.Info("dice result", "result", result, "rolled", result.Rolled)
// 		response.Success = result.Result > check.Difficult
// 		response.Description = result.Description
// 		response.Result = result.Result
// 		response.Rolled = result.Rolled
// 	}
// 	return response, nil
// }

// // D20 only
// func (c *Creature) calcSkillModifier(skill string) int {
// 	switch c.RPG.SkillRank {
// 	case rpg.SkillRanks:
// 		return c.Skills[skill].Value
// 	case rpg.SkillModifiers:
// 		if c.Skills[skill].Value == 1 {
// 			return 8
// 		}
// 		return 4
// 	case rpg.LevelBasedSkill:
// 		if c.Skills[skill].Value == 1 {
// 			value, _ := c.Extension.SkillBonus(skill)
// 			return value
// 		}
// 		return 0
// 	case rpg.OnePerOne:
// 		return c.Skills[skill].Value
// 	}
// 	return 0
// }
