package rules

import "github.com/betorvs/playbypost/core/types"

type Aligment int

const (
	LawfulGood Aligment = iota
	NeutralGood
	ChaoticGood
	LawfulNeutral
	Neutral
	ChaoticNeutral
	LawfulEvil
	NeutralEvil
	ChaoticEvil
)

func (a Aligment) String() string {
	switch a {
	case LawfulGood:
		return "lawful good"
	case NeutralGood:
		return "neutral good"
	case ChaoticGood:
		return "chaotic good"
	case LawfulNeutral:
		return "lawful neutral"
	case Neutral:
		return "neutral"
	case ChaoticNeutral:
		return "chaotic neutral"
	case LawfulEvil:
		return "lawful evil"
	case NeutralEvil:
		return "neutral evil"
	case ChaoticEvil:
		return "chaotic evil"
	}
	return types.Unknown
}

func AlignmentList() [9]Aligment {
	return [9]Aligment{LawfulGood, NeutralGood, ChaoticGood, LawfulNeutral, Neutral, ChaoticNeutral, LawfulEvil, NeutralEvil, ChaoticEvil}
}
