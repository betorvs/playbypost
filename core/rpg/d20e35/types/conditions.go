package types

/*
State represents state of character: alive, dead, blinded, confused
*/
type State int

const (
	Alive State = iota
	AbilityBurn
	AbilityDamaged
	AbilityDrained
	Blinded
	BlownAway
	Checked
	Confused
	Cowering
	Dazed
	Dazzled
	Dead
	Deafened
	Disabled
	Dying
	EnergyDrained
	Entangled
	Exhausted
	Fascinated
	Fatigued
	FlatFooted
	Frightened
	Grappling
	Incorporeal
	Invisible
	KnockedDown
	Nauseated
	Panicked
	Paralyzed
	Petrified
	Pinned
	Prone
	Shaken
	Sickened
	Stable
	Staggered
	Stunned
	Sleeping
	Turned
	Unconscious
)

func (c State) String() string {
	switch c {
	case Alive:
		return "alive"
	case AbilityBurn:
		return "ability burn"
	case AbilityDamaged:
		return "ability damaged"
	case AbilityDrained:
		return "ability drained"
	case Blinded:
		return "blinded"
	case BlownAway:
		return "blownAway"
	case Checked:
		return "checked"
	case Confused:
		return "confused"
	case Cowering:
		return "cowering"
	case Dazed:
		return "dazed"
	case Dazzled:
		return "dazzled"
	case Dead:
		return "dead"
	case Deafened:
		return "deafened"
	case Disabled:
		return "disabled"
	case Dying:
		return "dying"
	case EnergyDrained:
		return "energy drained"
	case Entangled:
		return "entangled"
	case Exhausted:
		return "exhausted"
	case Fascinated:
		return "fascinated"
	case Fatigued:
		return "fatigued"
	case FlatFooted:
		return "flatFooted"
	case Frightened:
		return "frightened"
	case Grappling:
		return "grappling"
	case Incorporeal:
		return "incorporeal"
	case Invisible:
		return "invisible"
	case KnockedDown:
		return "knocked down"
	case Nauseated:
		return "nauseated"
	case Panicked:
		return "panicked"
	case Paralyzed:
		return "paralyzed"
	case Petrified:
		return "petrified"
	case Pinned:
		return "pinned"
	case Prone:
		return "prone"
	case Shaken:
		return "shaken"
	case Sickened:
		return "sickened"
	case Sleeping:
		return "sleeping"
	case Stable:
		return "stable"
	case Staggered:
		return "staggered"
	case Stunned:
		return "stunned"
	case Turned:
		return "turned"
	case Unconscious:
		return "unconscious"
	}
	return Unknown
}
