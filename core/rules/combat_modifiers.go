package rules

import "github.com/betorvs/playbypost/core/types"

type CombatModifiers int

const (
	Zero CombatModifiers = iota
	Cover
	Concealment
	TotalConcealment
	Flanking
	Helpless
	HigherGround
	SqueezingThroughASpace
	KneelingSitting
)

func (c CombatModifiers) String() string {
	switch c {
	case Zero:
		return "Zero"
	case Cover:
		return "Cover"
	case Concealment:
		return "Concealment"
	case TotalConcealment:
		return "TotalConcealment"
	case Flanking:
		return "Flanking"
	case Helpless:
		return "Helpless"
	case HigherGround:
		return "HigherGround"
	case SqueezingThroughASpace:
		return "SqueezingThroughASpace"
	case KneelingSitting:
		return "KneelingSitting"
	}
	return types.Unknown
}

func IsAttacker(cond types.State, kind AttackTypes, combatMod CombatModifiers) int {
	switch cond {
	case types.Dazzled:
		return -1
	case types.Entangled:
		return -2
	case types.Invisible:
		return +2
	case types.Prone:
		if kind == Melee {
			return -4
		}
		return 0
	case types.Shaken, types.Frightened:
		return -2
	}
	switch combatMod {
	case Flanking:
		if kind == Melee {
			return +2
		}
		return 0
	case HigherGround:
		if kind == Melee {
			return +1
		}
		return 0
	case SqueezingThroughASpace:
		return -4
	}
	return 0
}

func IsDefender(cond types.State, kind AttackTypes, combatMod CombatModifiers) int {
	switch cond {
	case types.Blinded:
		return -2
	case types.Cowering:
		return -2
	case types.Entangled:
		return 0
	case types.Paralyzed, types.Sleeping:
		if kind == Melee {
			return -4
		}
		return 0
	case types.Pinned:
		if kind == Melee {
			return -4
		}
		return 0
	case types.Prone:
		if kind == Melee {
			return -4
		}
		return +4
	case types.Stunned:
		return -2
	}
	switch combatMod {
	case Cover:
		return +4
	case KneelingSitting:
		if kind == Melee {
			return -2
		}
		return +2
	case SqueezingThroughASpace:
		return -4
	}
	return 0
}
