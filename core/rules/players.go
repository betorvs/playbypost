package rules

import (
	"fmt"

	"github.com/betorvs/playbypost/core/types"
)

type Player struct {
	Name                  string
	PlayerName            string
	MultiClass            bool
	Class                 []Class
	Level                 int
	Race                  Race
	Aligment              Aligment
	Metadata              Metadata
	Abilities             Ability
	HitPoints             int
	ArmorClass            int
	SpellResistance       int
	PreferenceAttackIndex int
	AttackOption          []AttackOption
	Mounted               bool
	SkillsList            []Skills
	Feats                 []string
	SpecialAbilities      []string
	Spells                map[int][]string
	Languages             []string
	Gear                  Gear
	Bonuses               *TemporaryCombatBonuses
	totalDefense          bool
	state                 types.State
}

type Metadata struct {
	Deity  string
	Age    int
	Gender int
	Height int
	Weight int
	Eyes   string
	Hair   string
	Skin   string
}

type Ability struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

/*
Gear

https://www.dandwiki.com/wiki/SRD:Additional_Magic_Item_Space
*/
type Gear struct {
	Head        string // one headband, hat, or helmet
	Eyes        string // one pair of eye lenses or goggles
	Back        string // one cloak, cape, or mantle
	Neck        string // one amulet, brooch, medallion, necklace, periapt, or scarab
	Armor       string // one suit of armor; one robe; one vest, vestment, or shirt
	Arm         string // one pair of bracers or bracelets
	Hand        string // one pair of gloves or gauntlets
	RightRing   string
	LeftRing    string
	Waist       string // one belt
	Foot        string // and one pair of boots.
	Shield      string
	Possessions []string
	WeightTotal int
}

func (p Player) CalcSavingThrowsTotal(r string) int {
	switch r {
	case Fortitude:
		// Constitution
		bonus := p.Class[0].SavingThrowsCalcBonus(Fortitude)
		if p.MultiClass {
			for _, c := range p.Class {
				bonus += c.SavingThrowsCalcBonus(Fortitude)
			}
		}
		return bonus + CalcAbilityModifier(p.Abilities.Constitution)
	case Reflex:
		// Dexterity
		bonus := p.Class[0].SavingThrowsCalcBonus(Dexterity)
		if p.MultiClass {
			for _, c := range p.Class {
				bonus += c.SavingThrowsCalcBonus(Dexterity)
			}
		}
		return bonus + CalcAbilityModifier(p.Abilities.Dexterity)
	case Will:
		// Wisdom
		bonus := p.Class[0].SavingThrowsCalcBonus(Wisdom)
		if p.MultiClass {
			for _, c := range p.Class {
				bonus += c.SavingThrowsCalcBonus(Wisdom)
			}
		}
		return bonus + CalcAbilityModifier(p.Abilities.Wisdom)
	}
	return 0
}

func (p Player) CalcSkillTotal(s Skills) int {
	var value int
	ability := s.AbilityKey()
	if ability != types.Unknown {
		switch ability {
		case Strength:
			value = p.Abilities.Strength
		case Dexterity:
			value = p.Abilities.Strength
		case Constitution:
			value = p.Abilities.Strength
		case Intelligence:
			value = p.Abilities.Strength
		case Wisdom:
			value = p.Abilities.Strength
		case Charisma:
			value = p.Abilities.Strength
		default:
			value = 0
		}
	}
	skillBonus := p.Class[0].SkillCalcBonus(s)
	if p.MultiClass {
		for _, c := range p.Class {
			if c.SkillCalcBonus(s) > skillBonus {
				skillBonus = c.SkillCalcBonus(s)
			}
		}
	}
	return skillBonus + CalcAbilityModifier(value)
}

/*
CalcAttackTotal func

Your attack bonus with a melee weapon is:

Base attack bonus + Strength modifier + size modifier
With a ranged weapon, your attack bonus is:

Base attack bonus + Dexterity modifier + size modifier + range penalty
*/
func (p Player) CalcAttackTotal(kind AttackTypes) int {
	bonus := p.Class[0].AttackBase()
	if p.MultiClass {
		multiclassBonus := 0
		for _, c := range p.Class {
			multiclassBonus += c.AttackBase()
		}
		bonus = multiclassBonus
	}
	fmt.Println("bonus", bonus, "calc", CalcAbilityModifier(p.Abilities.Strength), "size", p.Race.Size.AttackModifier())
	switch kind {
	case Melee:
		return bonus + CalcAbilityModifier(p.Abilities.Strength) + p.Race.Size.AttackModifier()
	case Ranged:
		return bonus + CalcAbilityModifier(p.Abilities.Dexterity) + p.Race.Size.AttackModifier()
	}

	return bonus
}

func (p *Player) ChangeCondition(state types.State) error {
	p.state = state
	return nil
}

func (p *Player) SetTotalDefense() {
	p.totalDefense = true
}

func (p *Player) RemoveTotalDefense() {
	p.totalDefense = false
}
