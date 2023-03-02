package rules

import (
	"github.com/betorvs/playbypost/core/types"
)

const (
	NonStackableBonus int = 0
	StackableBonus    int = -1
	MaxInherentBonus  int = 5
	NotApplicableTTL  int = -2
)

type TypeBonus int

const (
	EnhancementBonus TypeBonus = iota
	AlchemicalBonus
	DeflectionBonus
	NaturalArmorBonus
	DodgeBonus // stacks
	ArmorBonus
	ShieldBonus
	CompetenceBonus
	InherentBonus // max 5 per Abilities
	InsightBonus
	LuckBonus
	MoraleBonus
	ProfaneBonus
	RacialBonus
	ResistanceBonus
	SacredBonus
	SizeBonus
	CircumstanceBonus // different context, stack
)

func (c TypeBonus) String() string {
	switch c {
	case EnhancementBonus:
		return "EnhancementBonus"
	case AlchemicalBonus:
		return "AlchemicalBonus"
	case DeflectionBonus:
		return "DeflectionBonus"
	case NaturalArmorBonus:
		return "NaturalArmorBonus"
	case DodgeBonus:
		return "DodgeBonus"
	case ArmorBonus:
		return "ArmorBonus"
	case ShieldBonus:
		return "ShieldBonus"
	case CompetenceBonus:
		return "CompetenceBonus"
	case InherentBonus:
		return "InherentBonus"
	case InsightBonus:
		return "InsightBonus"
	case LuckBonus:
		return "LuckBonus"
	case MoraleBonus:
		return "MoraleBonus"
	case ProfaneBonus:
		return "ProfaneBonus"
	case RacialBonus:
		return "RacialBonus"
	case ResistanceBonus:
		return "ResistanceBonus"
	case SacredBonus:
		return "SacredBonus"
	case SizeBonus:
		return "SizeBonus"
	case CircumstanceBonus:
		return "CircumstanceBonus"
	}
	return types.Unknown
}

type CombatBonus int

const (
	AttackBonus CombatBonus = iota
	DefenseBonus
	DamageBonus
	HitPointsBonus
	FortitudeBonus
	ReflexBonus
	WillBonus
	StrengthBonus
	DexterityBonus
	ConstitutionBonus
	IntelligenceBonus
	WisdomBonus
	CharismaBonus
)

func (c CombatBonus) String() string {
	switch c {
	case AttackBonus:
		return "AttackBonus"
	case DefenseBonus:
		return "DefenseBonus"
	case DamageBonus:
		return "DamageBonus"
	case HitPointsBonus:
		return "HitPointsBonus"
	case FortitudeBonus:
		return "FortitudeBonus"
	case ReflexBonus:
		return "ReflexBonus"
	case WillBonus:
		return "WillBonus"
	case StrengthBonus:
		return "StrengthBonus"
	case DexterityBonus:
		return "DexterityBonus"
	case ConstitutionBonus:
		return "ConstitutionBonus"
	case IntelligenceBonus:
		return "IntelligenceBonus"
	case WisdomBonus:
		return "WisdomBonus"
	case CharismaBonus:
		return "CharismaBonus"
	}
	return types.Unknown
}

/*
TemporaryCombatBonuses struct holds all temporary Bonuses applied to a player or Monster
*/
type TemporaryCombatBonuses struct {
	Attack       int
	Defense      int
	Damage       int
	HitPoints    int
	Fortitude    int
	Reflex       int
	Will         int
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
	Source       map[CombatBonus]DescriptionBonus
}

/*
DescriptionBonus struct holds information about a specific bonus type.
each entry, should be added to Values map and it should be unique in the map
*/
type DescriptionBonus struct {
	Total        int
	Sources      map[string]TypeBonus
	Values       map[string]int
	WithTTL      map[string]int
	TotalPerType map[TypeBonus]int
}

type TypeDescriptionBonus struct {
}

/**/
func NewTemporaryCombatBonuses() *TemporaryCombatBonuses {
	c := TemporaryCombatBonuses{}
	source := make(map[CombatBonus]DescriptionBonus)
	c.Source = source
	// fmt.Printf("factory %+v", c)
	return &c
}

func (c TemporaryCombatBonuses) GetBonus(target CombatBonus) int {
	switch target {
	case AttackBonus:
		return c.Attack
	case DefenseBonus:
		return c.Defense
	case DamageBonus:
		return c.Damage
	case HitPointsBonus:
		return c.HitPoints
	case FortitudeBonus:
		return c.Fortitude
	case ReflexBonus:
		return c.Reflex
	case WillBonus:
		return c.Will
	case StrengthBonus:
		return c.Strength
	case DexterityBonus:
		return c.Dexterity
	case ConstitutionBonus:
		return c.Constitution
	case IntelligenceBonus:
		return c.Intelligence
	case WisdomBonus:
		return c.Wisdom
	case CharismaBonus:
		return c.Charisma
	}
	return 0
}

func (c *TemporaryCombatBonuses) CalcTotal(target CombatBonus, value int) {
	switch target {
	case AttackBonus:
		c.Attack = value
	case DefenseBonus:
		c.Defense = value
	case DamageBonus:
		c.Damage = value
	case HitPointsBonus:
		c.HitPoints = value
	case FortitudeBonus:
		c.Fortitude = value
	case ReflexBonus:
		c.Reflex = value
	case WillBonus:
		c.Will = value
	case StrengthBonus:
		c.Strength = value
	case DexterityBonus:
		c.Dexterity = value
	case ConstitutionBonus:
		c.Constitution = value
	case IntelligenceBonus:
		c.Intelligence = value
	case WisdomBonus:
		c.Wisdom = value
	case CharismaBonus:
		c.Charisma = value
	}
}

func (c *TemporaryCombatBonuses) AddBonuses(kind TypeBonus, target CombatBonus, value, ttl int, source string) {
	switch kind {
	case AlchemicalBonus, ArmorBonus, CompetenceBonus, DeflectionBonus, EnhancementBonus, InsightBonus, LuckBonus, MoraleBonus, NaturalArmorBonus, ProfaneBonus, RacialBonus, ResistanceBonus, SacredBonus, ShieldBonus, SizeBonus:
		endValue := calcStackBonus(kind, target, &c.Source, value, ttl, NonStackableBonus, source)
		c.CalcTotal(target, endValue)

	case DodgeBonus, CircumstanceBonus:
		// dogde always stack
		// CircumstanceBonus sometimes.
		endValue := calcStackBonus(kind, target, &c.Source, value, ttl, StackableBonus, source)
		c.CalcTotal(target, endValue)
		// stacks
	case InherentBonus:
		// max 5 per Abilities
		endValue := calcStackBonus(kind, target, &c.Source, value, ttl, MaxInherentBonus, source)
		c.CalcTotal(target, endValue)
	}
}

func calcStackBonus(kind TypeBonus, target CombatBonus, sourceDescription *map[CombatBonus]DescriptionBonus, value, ttl int, maxValue int, source string) int {
	// fmt.Println("calcStackBonus", sourceDescription)
	if maxValue > NonStackableBonus && value > maxValue {
		// only change it in case that maxValue is greather than zero
		value = maxValue
	}
	var endValue int
	if d, ok := (*sourceDescription)[target]; ok {
		d.Values[source] = value
		d.Sources[source] = kind
		if ttl != NotApplicableTTL {
			d.WithTTL[source] = ttl
		}
		if value > d.TotalPerType[kind] {
			d.TotalPerType[kind] = value
		}
		if maxValue == StackableBonus {
			// stackable bonus
			// maxValue equal -1 means that this is a stackable value
			// d.Values[source] = value
			for k, v := range d.Sources {
				if v == kind {
					endValue += d.Values[k]
				}
			}
			d.TotalPerType[kind] = endValue
		}
		allValues := 0
		// fmt.Println("sources", d.Sources)
		for _, v := range d.TotalPerType {
			allValues += v
		}
		d.Total = allValues
		endValue = d.Total
		// fmt.Println("total", d.Total)
		// fmt.Println("total per type", d.TotalPerType)
	} else {
		desc := make(map[CombatBonus]DescriptionBonus)
		// sources
		sources := make(map[string]TypeBonus)
		sources[source] = kind
		// values
		values := make(map[string]int)
		values[source] = value
		// nonstackable bonus
		totalPerType := make(map[TypeBonus]int)
		totalPerType[kind] = value
		// with ttl
		withTTL := make(map[string]int)

		n := DescriptionBonus{
			Total:        value,
			Sources:      sources,
			Values:       values,
			TotalPerType: totalPerType,
		}
		if ttl != NotApplicableTTL {
			withTTL[source] = ttl
			n.WithTTL = withTTL
		}

		desc[target] = n
		(*sourceDescription) = desc
		endValue = value
	}
	return endValue
}

func (c *TemporaryCombatBonuses) DecreaseTTL(kind TypeBonus, target CombatBonus, source string) {
	var markToDelete bool
	if d, ok := c.Source[target]; ok {
		if d.WithTTL[source] != NotApplicableTTL {
			d.WithTTL[source]--
			if d.WithTTL[source] == 0 {
				markToDelete = true
			}
		}
	}
	if markToDelete {
		c.removeBonusReCalc(kind, target, source)
		// call reCalc function for that target
	}
}

// removeBonusReCalc
func (c *TemporaryCombatBonuses) removeBonusReCalc(kind TypeBonus, target CombatBonus, source string) {

	// call reCalc function for that target
	// loop into Temporary source to re calculate it
	// and remove extra bonus
	valuesTemp := c.Source[target]
	delete(valuesTemp.Values, source)
	delete(valuesTemp.Sources, source)
	delete(valuesTemp.WithTTL, source)
	// delete(c.Source, kind)
	// internal := make(map[string]int)
	// desc := TypeDescriptionBonus{
	// 	WithTTL: internal,
	// 	Values:  internal,
	// }
	// c.Source[target] = desc

	for k, v := range valuesTemp.Values {
		if d, ok := valuesTemp.WithTTL[k]; ok {
			c.AddBonuses(kind, target, v, d, k)
			continue
		}
		c.AddBonuses(kind, target, v, NotApplicableTTL, k)
	}
}
