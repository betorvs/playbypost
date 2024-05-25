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
	Dice     *rpg.Roll
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
	switch a.Attacker.RPG.Name {
	case rpg.D2035:
		a.Logger.Info("d20 single melee attack call")
		// level bonus + str against AC
		// damage require rolls : dependency weapon dices + bonus
		result, res, _ := a.Dice.Check("attack roll")
		// if err != nil {
		// 	return
		// }
		bonus, _ := a.Attacker.Extension.AttackBonus(a.Weapon)
		// if err != nil {
		// 	return
		// }
		abilityBonus := a.Attacker.calcAbilityModifier(strenght)
		weaponBonus, diceDamage, _ := a.Attacker.Extension.WeaponBonus(a.Weapon)
		a.Logger.Info("dice result", "result", result+bonus+weaponBonus, "rolled", res)
		defensor, _ := a.Defensor.Extension.DefenseBonus("")
		a.Response.Success = result+bonus+abilityBonus+weaponBonus >= defensor
		if result == 20 {
			a.Logger.Info("20 natural roll")
			a.Response.Success = true
		}
		a.Response.Text = fmt.Sprint("result", result+bonus+abilityBonus+weaponBonus, "rolled", res)
		if a.Response.Success {
			damage, res, _ := a.Dice.FreeRoll("damage", diceDamage)
			a.Response.Damage = damage + abilityBonus
			a.Response.Text += fmt.Sprint("damage roll ", res)
			a.Defensor.Extension.Damage(damage + abilityBonus)
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
		calcDices := a.Dice.Dice(abilityBonus+bonus+weaponBonus-defense, 0)
		result, res, _ := a.Dice.FreeRoll("attack roll", calcDices)
		a.Logger.Info("dice result", "result", abilityBonus+bonus+weaponBonus-defense, "rolled", res)
		if result > 0 {
			a.Response.Success = true
			a.Response.Text = res
			a.Defensor.Extension.Damage(result)
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
		calcDices := a.Dice.Dice(abilityBonus+bonus, 6)
		result, res, _ := a.Dice.FreeRoll("attack roll", calcDices)

		defenseAbilityBonus := a.Attacker.calcAbilityModifier(dexterity)
		defensebonus := a.Attacker.calcSkillModifier(melee)
		defenseCalcDices := a.Dice.Dice(defenseAbilityBonus+defensebonus, 6)
		defenseResult, defenseRes, _ := a.Dice.FreeRoll("defense roll", defenseCalcDices)
		a.Logger.Info("results", "round", a.Round, "attack_details", res, "defense_details", defenseRes)

		if result >= defenseResult {
			a.Response.Success = true
			damageAbility := a.Attacker.calcAbilityModifier(strenght)
			weaponBonus, _, _ := a.Attacker.Extension.WeaponBonus(a.Weapon)
			calcDices := a.Dice.Dice(damageAbility+weaponBonus, 6)
			damageResult, damageRes, _ := a.Dice.FreeRoll("damage roll", calcDices)
			a.Response.Damage = damageResult
			if a.Defensor.Extension.HealthStatus() <= 0 {
				_ = a.Defensor.Destroy()
			}
			a.Logger.Info("results", "round", a.Round, "damage_details", damageRes)
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

func NewAttack(round, weapon string, kind AttackTypes, attacker *Creature, defensor *Creature, dice *rpg.Roll, logger *slog.Logger) *Attack {
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
