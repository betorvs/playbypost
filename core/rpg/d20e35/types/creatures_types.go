package types

type Creatures int

const (
	Aberration Creatures = iota
	Animal
	Celestial
	Construct
	Dragon
	Elemental
	Fey
	Fiend
	Giant
	Humanoid
	MagicalBeast
	MonstrousHumanoid
	Ooze
	Outsider
	Plant
	Undead
	Vermin
)

func (c Creatures) String() string {
	switch c {
	case Aberration:
		return "aberration"
	case Animal:
		return "animal"
	case Celestial:
		return "celestial"
	case Construct:
		return "construct"
	case Dragon:
		return "dragon"
	case Elemental:
		return "elemental"
	case Fey:
		return "fey"
	case Fiend:
		return "fiend"
	case Giant:
		return "giant"
	case Humanoid:
		return "humanoid"
	case MagicalBeast:
		return "magical beast"
	case MonstrousHumanoid:
		return "monstrous humanoid"
	case Ooze:
		return "ooze"
	case Outsider:
		return "outsider"
	case Plant:
		return "plant"
	case Undead:
		return "undead"
	case Vermin:
		return "vermin"
	}
	return Unknown
}
