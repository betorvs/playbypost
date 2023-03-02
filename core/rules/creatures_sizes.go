package rules

import "github.com/betorvs/playbypost/core/types"

type CreaturesSizes int

const (
	Fine CreaturesSizes = iota
	Diminutive
	Tiny
	Small
	Medium
	Large
	Huge
	Gargantuan
	Colossal
)

func (c CreaturesSizes) String() string {
	switch c {
	case Fine:
		return "fine"
	case Diminutive:
		return "diminutive"
	case Tiny:
		return "tiny"
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Large:
		return "large"
	case Huge:
		return "huge"
	case Gargantuan:
		return "gargantuan"
	case Colossal:
		return "colossal"
	}
	return types.Unknown
}

func (c CreaturesSizes) AttackModifier() int {
	switch c {
	case Fine:
		return +8
	case Diminutive:
		return +4
	case Tiny:
		return +2
	case Small:
		return +1
	case Medium:
		return 0
	case Large:
		return -1
	case Huge:
		return -2
	case Gargantuan:
		return -4
	case Colossal:
		return -8
	}
	return 0
}
