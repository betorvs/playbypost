package rules

import (
	"fmt"
	"log/slog"

	"github.com/betorvs/playbypost/core/rpg"
)

type AttackTypes int

const (
	Unarmed AttackTypes = iota
	Melee
	Ranged
	MeleeTouch
	RangedTouch
)

type Combat interface {
	Call()
	Undo()
}

type Attack struct {
	Round    string
	Weapon   string
	Kind     AttackTypes
	Attacker *Creature
	Defensor *Creature
	Dice     rpg.RollInterface
	Response Response
	Logger   *slog.Logger
}

type Response struct {
	Success bool
	Text    string
	Damage  int
}

func (a *Attack) Call() {
	switch a.Kind {
	case Unarmed:
		a.singleUnarmedAttack()
	case Melee:
		a.singleMeleeAttack()
	case Ranged:
		a.singleRangedAttack()

		// case MeleeTouch:
		// 	//
		// case RangedTouch:
		// 	//
	}
}

func (a *Attack) singleMeleeAttack() {
	a.Logger.Info("single melee attack call", "system_name", a.Attacker.RPG.Name)
	strenght := "strenght"
	diceName := fmt.Sprintf("attack-roll-%s-%s", a.Attacker.Name, strenght)
	switch a.Attacker.RPG.Name {
	case rpg.D2035:
		a.Logger.Info("d20 single melee attack call")
		// level bonus + str against AC
		// damage require rolls : dependency weapon dices + bonus
		result, _ := a.Dice.Check(diceName)
		// if err != nil {
		// 	return
		// }
		bonus, _ := a.Attacker.Extension.AttackBonus(a.Weapon)
		// if err != nil {
		// 	return
		// }
		abilityBonus := a.Attacker.calcAbilityModifier(strenght)
		weaponBonus, diceDamage, _ := a.Attacker.Extension.WeaponBonus(a.Weapon)
		a.Logger.Info("dice result", "result", result.Result+bonus+weaponBonus, "rolled", result.Rolled)
		defensor, _ := a.Defensor.Extension.DefenseBonus("")
		a.Response.Success = result.Result+bonus+abilityBonus+weaponBonus >= defensor
		if result.Result == 20 {
			a.Logger.Info("20 natural roll")
			a.Response.Success = true
		}
		a.Response.Text = fmt.Sprint("result", result.Result+bonus+abilityBonus+weaponBonus, "rolled", result.Rolled)
		if a.Response.Success {
			damage, _ := a.Dice.FreeRoll("damage", diceDamage)
			a.Response.Damage = damage.Result + abilityBonus
			a.Response.Text += fmt.Sprint("damage roll ", result.Rolled)
			err := a.Defensor.Extension.Damage(damage.Result + abilityBonus)
			if err != nil {
				a.Logger.Error("error on damage calc", "error", err)
			}
			if a.Defensor.Extension.HealthStatus() <= 0 {
				_ = a.Defensor.Destroy()
			}
		}
		a.Logger.Info("results", "round", a.Round, "details", a.Response.Text)

	case rpg.D10HM:
		// str + weaponry - defense (lower between dex or wits)
		// damage is equal success
		weaponry := "weaponry"
		abilityBonus := a.Attacker.calcAbilityModifier(strenght)
		bonus := a.Attacker.calcSkillModifier(weaponry)
		weaponBonus, _, _ := a.Attacker.Extension.WeaponBonus(a.Weapon)
		defense, _ := a.Defensor.Extension.DefenseBonus("melee")
		calcDices := a.Dice.FormatDice(abilityBonus+bonus+weaponBonus-defense, 0)
		a.Logger.Info("calcDices", "dices", calcDices)
		result, _ := a.Dice.FreeRoll(diceName, calcDices)
		a.Logger.Info("dice result", "result", abilityBonus+bonus+weaponBonus-defense, "rolled", result.Rolled)
		a.Response.Text = result.Rolled
		if result.Result > 0 {
			a.Response.Success = true
			a.Response.Damage = result.Result
			err := a.Defensor.Extension.Damage(result.Result)
			if err != nil {
				a.Logger.Error("error on damage calc", "error", err)
			}
			if a.Defensor.Extension.HealthStatus() <= 0 {
				_ = a.Defensor.Destroy()
			}
		}
		a.Logger.Info("results", "round", a.Round, "details", a.Response.Text)
		// if err != nil {
		// 	return response, err
		// }

	case rpg.D10OS:
		// dex + melee - against - dex + melee || dodge
		// damage require rolls = str + weapon bonus
		melee := "melee"
		dexterity := "dexterity"
		abilityBonus := a.Attacker.calcAbilityModifier(dexterity)
		bonus := a.Attacker.calcSkillModifier(melee)
		calcDices := a.Dice.FormatDice(abilityBonus+bonus, 6)
		result, _ := a.Dice.FreeRoll(diceName, calcDices)

		defenseAbilityBonus := a.Attacker.calcAbilityModifier(dexterity)
		defensebonus := a.Attacker.calcSkillModifier(melee)
		defenseCalcDices := a.Dice.FormatDice(defenseAbilityBonus+defensebonus, 6)
		defenseDiceName := fmt.Sprintf("defense-roll-%s-%s", a.Defensor.Name, dexterity)
		defenseResult, _ := a.Dice.FreeRoll(defenseDiceName, defenseCalcDices)
		a.Logger.Info("results", "round", a.Round, "attack_details", result.Rolled, "defense_details", defenseResult.Rolled)

		if result.Result >= defenseResult.Result {
			a.Response.Success = true
			damageAbility := a.Attacker.calcAbilityModifier(strenght)
			weaponBonus, _, _ := a.Attacker.Extension.WeaponBonus(a.Weapon)
			calcDices := a.Dice.FormatDice(damageAbility+weaponBonus, 6)
			damageDiceName := fmt.Sprintf("damage-roll-%s-%s", a.Attacker.Name, strenght)
			damageResult, _ := a.Dice.FreeRoll(damageDiceName, calcDices)
			a.Response.Damage = damageResult.Result
			if a.Defensor.Extension.HealthStatus() <= 0 {
				_ = a.Defensor.Destroy()
			}
			a.Logger.Info("results", "round", a.Round, "damage_details", damageResult.Rolled)
		}
	}
}

func (a *Attack) singleRangedAttack() {
	switch a.Attacker.RPG.BaseSystem {
	case rpg.D2035:
		// level bonus + dex against AC
		// damage require rolls : dependency weapon dices + bonus

	case rpg.D10HM:
		// dex + firearms - armor
		// damage is equal success

	case rpg.D10OS:
		// dex + firearms - against - dex + dodge
		// damage require rolls = weapon damage + accuracy
	}
}

func (a *Attack) singleUnarmedAttack() {
	switch a.Attacker.RPG.BaseSystem {
	case rpg.D2035:
		// level bonus + str against AC
		// damage require rolls
		// monk class features

	case rpg.D10HM:
		// str + brawl - defense + armor
		// damage is equal success

	case rpg.D10OS:
		// dex + brawl - against - dex + dodge
		// damage require rolls = Str
	}
}

// func (a *Attack) totalAttack() {}

func NewAttack(round, weapon string, kind AttackTypes, attacker *Creature, defensor *Creature, dice rpg.RollInterface, logger *slog.Logger) *Attack {
	a := Attack{
		Round:    round,
		Weapon:   weapon,
		Kind:     kind,
		Attacker: attacker,
		Defensor: defensor,
		Dice:     dice,
		Response: Response{},
		Logger:   logger,
	}
	return &a
}
