package rules

import (
	"errors"
	"strconv"

	"github.com/betorvs/playbypost/core/sys/diceroll"
	"github.com/betorvs/playbypost/core/types"
)

var ErrAttackTypes = errors.New("cannot convert string to AttackTypes")

type AttackTypes int

const (
	Melee AttackTypes = iota
	Ranged
	MeleeTouch
	RangedTouch
)

func (a AttackTypes) String() string {
	switch a {
	case Melee:
		return "melee"
	case Ranged:
		return "ranged"
	case MeleeTouch:
		return "melee touch"
	case RangedTouch:
		return "ranged touch"
	}
	return types.Unknown
}

func (a AttackTypes) Atoi(s string) (AttackTypes, error) {
	switch s {
	case "melee", "Melee":
		return Melee, nil
	case "ranged", "Ranged":
		return Ranged, nil
	case "melee touch", "MeleeTouch":
		return MeleeTouch, nil
	case "ranged touch", "RangedTouch":
		return RangedTouch, nil
	}
	return -1, ErrAttackTypes
}

type AttackOption struct {
	AttackBonus               int
	Damage                    string
	Critical                  string
	Range                     int
	AttackerSize              CreaturesSizes
	Weapon                    *Weapon
	WeaponPreferredDamageType WeaponDamageType
	Type                      AttackTypes
	Notes                     string
}

type DamageOption struct {
	Successful       bool
	Damage           int
	Type             AttackTypes
	WeaponDamageType WeaponDamageType
}

func (a *AttackOption) Attack(valueFromDice, attackBonus, armorClass int, dice diceroll.DiceRoller) (DamageOption, error) {
	damage := DamageOption{}
	damage.Successful = a.CheckAttack(valueFromDice, attackBonus, armorClass)
	if damage.Successful {
		damage.Damage = a.calcDamage(dice)
		critical, value := a.Weapon.IsCritical(valueFromDice, false)
		if critical {
			damage.Damage = a.calcDamage(dice) * value
		}
	}
	return damage, nil
}

func (a *AttackOption) calcDamage(dice diceroll.DiceRoller) int {

	switch a.AttackerSize {
	case Small:
		damage, _ := strconv.Atoi(a.Weapon.DamageSmall)
		return damage
	case Medium:
		damage, _ := strconv.Atoi(a.Weapon.DamageMedium)
		return damage
	}

	return 0
}

func (a *AttackOption) CheckAttack(valueFromDice, attackBonus, armorClass int) bool {
	total := valueFromDice + attackBonus
	/*
		cases:
			- 20 natural always hit
			- dice + bonus higher than armor class
			- critical, to calculate damage, can be different than 20 for certain weapons
	*/
	switch {
	case valueFromDice == 20:
		return true
	case total >= armorClass:
		return true
	}
	return false
}
