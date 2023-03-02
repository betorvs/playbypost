package rules

import "github.com/betorvs/playbypost/core/types"

type WeaponDamageType int

const (
	Bludgeoning WeaponDamageType = iota
	Piercing
	Slashing
)

func (w WeaponDamageType) String() string {
	switch w {
	case Bludgeoning:
		return "bludgeoning"
	case Piercing:
		return "piercing"
	case Slashing:
		return "slashing"
	}
	return types.Unknown
}

type WeaponCategory int

const (
	Simple WeaponCategory = iota
	Martial
	Exotic
)

func (w WeaponCategory) String() string {
	switch w {
	case Simple:
		return "simple"
	case Martial:
		return "martial"
	case Exotic:
		return "exotic"
	}
	return types.Unknown
}

// Weapon struct
type Weapon struct {
	Name                  string           `json:"name"`
	Title                 string           `json:"title"`
	AttackType            AttackTypes      `json:"attack_type"`
	Category              WeaponCategory   `json:"category"`
	Cost                  int              `json:"cost"`
	CoinType              CoinsType        `json:"coin_type"`
	RangeIncrement        int              `json:"range_increment"`
	RangeIncrementMeasure string           `json:"range_increment_measure"`
	DamageMedium          string           `json:"damage_medium"`
	DamageSmall           string           `json:"damage_small"`
	Critical              string           `json:"critical,omitempty"`
	CriticalSecondary     string           `json:"critical_secondary,omitempty"`
	DamageType            WeaponDamageType `json:"damage_type"`
	SecondaryDamageType   WeaponDamageType `json:"secondary_damage_type"`
	Weight                int              `json:"weight"`
	Measure               string           `json:"measure"`
	Properties            string           `json:"properties"`
}

func (w Weapon) IsCritical(v int, improvedCritical bool) (bool, int) {
	critical := w.Critical
	// if player has improved critical for it
	if improvedCritical {
		switch w.Critical {
		case "x3/x4":
			critical = "19-20/x3/x4"
		case "x2":
			critical = "19-20/x2"
		case "x3":
			critical = "19-20/x3"
		case "x4":
			critical = "19-20/x4"
		case "18-20/x2":
			critical = "15-20/x2"
		case "19-20/x2":
			critical = "17-20/x2"
		}
	}
	switch critical {
	case "x3/x4":
		// need to check based on which attack was used
		return isTwentyNatural(v), 3
	case "x2":
		return isTwentyNatural(v), 2
	case "x3":
		return isTwentyNatural(v), 3
	case "x4":
		return isTwentyNatural(v), 4
	case "18-20/x2":
		if v >= 18 && v <= 20 {
			return true, 2
		}
	case "19-20/x2":
		if v == 19 || v == 20 {
			return true, 2
		}
	// Improved Critical
	case "19-20/x3/x4", "19-20/x3":
		if v == 19 || v == 20 {
			return true, 3
		}
	// Improved Critical
	case "19-20/x4":
		if v == 19 || v == 20 {
			return true, 4
		}
	// Improved Critical
	case "17-20/x2":
		if v >= 17 && v <= 20 {
			return true, 2
		}
	// Improved Critical
	case "15-20/x2":
		if v >= 15 && v <= 20 {
			return true, 2
		}
	}
	return false, 1
}

func isTwentyNatural(v int) bool {
	return v == 20
}
