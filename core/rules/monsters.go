package rules

import "github.com/betorvs/playbypost/core/types"

type Monster struct {
	Name  string
	Title string
	Creature
	// Type                  Creatures
	// Size                  CreaturesSizes
	// HitDice               string
	// HitPoints             int
	Initiative int
	// Speed                 int
	// ArmorClass            int
	BaseAttack int
	// PreferenceAttackIndex int // should be validated with len(AttackOption)
	// AttackOption          []AttackOption
	// Mounted               bool
	SpecialAttack    []AttackOption
	SpecialQualities []string
	// SkillsList            []Skills
	// Feats                 []string
	ChallengeRate float64
	Aligment      Aligment
	// Bonuses               *TemporaryCombatBonuses
	// state                 types.State
	// totalDefense          bool
	// Ability
	// SavingThrows
}

func (c *Monster) ChangeCondition(state types.State) error {
	c.state = state
	return nil
}

func (c *Monster) SetTotalDefense() {
	c.totalDefense = true
}

func (c *Monster) RemoveTotalDefense() {
	c.totalDefense = false
}

func (c *Monster) GetTotalDefense() bool {
	return c.totalDefense
}
