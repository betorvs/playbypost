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
