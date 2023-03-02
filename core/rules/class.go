package rules

import (
	"math"

	"github.com/betorvs/playbypost/core/types"
)

const (
	Fortitude string = "fortitude"
	Reflex    string = "reflex"
	Will      string = "will"
)

type KindBonus int

const (
	Poor KindBonus = iota
	Average
	Good
)

func (k KindBonus) String() string {
	switch k {
	case Poor:
		return "poor"
	case Average:
		return "average"
	case Good:
		return "good"
	}
	return types.Unknown
}

type HitDices int

const (
	D4 HitDices = iota
	D6
	D8
	D10
	D12
)

func (h HitDices) String() string {
	switch h {
	case D4:
		return "d4"
	case D6:
		return "d6"
	case D8:
		return "d8"
	case D10:
		return "d10"
	case D12:
		return "d12"
	}
	return types.Unknown
}

// Barbarian
// Bard
// Cleric
// Druid
// Fighter
// Monk
// Paladin
// Ranger
// Rogue
// Sorcerer
// Wizard

/*
Class

Base Attack
Base Savings bonus fortitude reflex will
skill list
Weapon Proficiency
Armor Proficiency
Special features per level
*/

type Class struct {
	Level             int
	SkillList         []Skills
	WeaponProficiency []string
	ArmorProficiency  []string
	SpecialFeatures   []string
	HitDie            HitDices
	AttackBonus       KindBonus
	SavingThrows
}

type ByLevel struct {
	AttackBonus int
}

type SavingThrows struct {
	Fortitude KindBonus
	Reflex    KindBonus
	Will      KindBonus
}

func (c Class) AttackBase() int {
	return c.baseAttackCalc(c.AttackBonus)
}

/*
baseAttackCalc func

https://nwn2.fandom.com/wiki/Base_attack_bonus
*/
func (c Class) baseAttackCalc(kind KindBonus) int {
	// fmt.Println("base calc", c.Level, kind)
	var result float64
	switch kind {
	case Good:
		result = float64(c.Level)
	case Average:
		result = math.Floor(float64(c.Level) * 0.75)
	case Poor:
		result = math.Floor(float64(c.Level) * 0.5)
	}
	return int(result)
}

/*
SkillCalcBonus

Class Skills: 1d20 + character's level + modifiers

Cross-Class Skills: 1d20 + modifiers
*/
func (c Class) SkillCalcBonus(s Skills) int {
	if c.search(s) {
		return c.Level
	}
	return 0
}

func (c Class) search(s Skills) bool {
	for _, v := range c.SkillList {
		if v == s {
			return true
		}
	}
	return false
}

/*
SavingThrowsCalcBonus func
*/
func (c Class) SavingThrowsCalcBonus(s string) int {
	switch s {
	case Fortitude:
		return c.baseSaveCalc(c.Fortitude)
	case Reflex:
		return c.baseSaveCalc(c.Reflex)
	case Will:
		return c.baseSaveCalc(c.Will)
	}
	return 0
}

/*
BaseSaveCalc func

https://nwn2.fandom.com/wiki/Base_save_bonus
*/
func (c Class) baseSaveCalc(s KindBonus) int {
	var result float64
	switch s {
	case Good:
		result = math.Floor((float64(c.Level) / 2) + 2)
	case Poor:
		result = math.Floor((float64(c.Level) / 3))
	}

	return int(result)
}

/*
https://forum.rpg.net/index.php?threads/d-d-3-3-5-xp-formula.228600/#:~:text=To%20get%20to%20a%20new,N%2D1))%20*%201000.
N = N*(N-1)*500
*/
func (c Class) calcNextLevelXP() int {
	level := c.Level + 1
	result := level * (level - 1) * 500
	return result
}
