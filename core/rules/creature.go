package rules

import "github.com/betorvs/playbypost/core/types"

type Creature struct {
	Race                  string
	HitPoints             int
	ArmorClass            int
	Speed                 int
	SpellResistance       int
	PreferenceAttackIndex int
	AttackOption          []AttackOption
	Mounted               bool
	SkillsList            []Skills
	Feats                 []string
	SpecialAbilities      []string
	Bonuses               *TemporaryCombatBonuses
	Type                  Creatures
	Size                  CreaturesSizes
	Ability
	SavingThrows
	totalDefense bool
	state        types.State
}
