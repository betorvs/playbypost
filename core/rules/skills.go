package rules

import "github.com/betorvs/playbypost/core/types"

/*
Based on https://www.dandwiki.com/wiki/UA:Alternative_Skill_Systems
Because the interaction between players and system will be done via chat systems and it can became very difficult to manage all new ranks and half ranks.

Class Skills: 1d20 + character's level + modifiers

Cross-Class Skills: 1d20 + modifiers

It will add https://www.dandwiki.com/wiki/UA:Skill_Knowledge_(Feat)

*/

type Skills int

const (
	Appraise Skills = iota
	Autohypnosis
	Balance
	Bluff
	Climb
	Concentration
	ControlShape
	Craft
	DecipherScript
	Diplomacy
	DisableDevice
	Disguise
	EscapeArtist
	Forgery
	GatherInformation
	HandleAnimal
	HealSkill
	Hide
	Intimidate
	Jump
	Knowledge
	Listen
	MoveSilently
	OpenLock
	Perform
	Psicraft
	Profession
	Ride
	Search
	SenseMotive
	SleightOfHand
	SpeakLanguage
	Spellcraft
	Spot
	Survival
	Swim
	Tumble
	UseMagicDevice
	UsePsionicDevice
	UseRope
)

func (s Skills) String() string {
	switch s {
	case Appraise:
		return "appraise"
	case Autohypnosis:
		return "autohypnosis"
	case Balance:
		return "balance"
	case Bluff:
		return "bluff"
	case Climb:
		return "climb"
	case Concentration:
		return "concentration"
	case ControlShape:
		return "control shape"
	case Craft:
		return "craft"
	case DecipherScript:
		return "decipher script"
	case Diplomacy:
		return "diplomacy"
	case DisableDevice:
		return "disable device"
	case Disguise:
		return "disguise"
	case EscapeArtist:
		return "escape artist"
	case Forgery:
		return "forgery"
	case GatherInformation:
		return "gather information"
	case HandleAnimal:
		return "handle animal"
	case HealSkill:
		return "heal"
	case Hide:
		return "hide"
	case Intimidate:
		return "intimidate"
	case Jump:
		return "jump"
	case Knowledge:
		return "knowledge"
	case Listen:
		return "listen"
	case MoveSilently:
		return "move silently"
	case OpenLock:
		return "open lock"
	case Perform:
		return "perform"
	case Psicraft:
		return "psicraft"
	case Profession:
		return "profession"
	case Ride:
		return "ride"
	case Search:
		return "search"
	case SenseMotive:
		return "sense motive"
	case SleightOfHand:
		return "sleight of hand"
	case SpeakLanguage:
		return "speak language"
	case Spellcraft:
		return "spellcraft"
	case Spot:
		return "spot"
	case Survival:
		return "survival"
	case Swim:
		return "swim"
	case Tumble:
		return "tumble"
	case UseMagicDevice:
		return "use magic device"
	case UsePsionicDevice:
		return "use psionic device"
	case UseRope:
		return "use rope"
	}
	return types.Unknown
}

func (s Skills) AbilityKey() string {
	switch s {
	case Appraise:
		return Intelligence
	case Autohypnosis:
		return Wisdom
	case Balance:
		return Dexterity
	case Bluff:
		return Charisma
	case Climb:
		return Strength
	case Concentration:
		return Constitution
	case ControlShape:
		return Wisdom
	case Craft:
		return Intelligence
	case DecipherScript:
		return Intelligence
	case Diplomacy:
		return Charisma
	case DisableDevice:
		return Intelligence
	case Disguise:
		return Charisma
	case EscapeArtist:
		return Dexterity
	case Forgery:
		return Intelligence
	case GatherInformation:
		return Charisma
	case HandleAnimal:
		return Charisma
	case HealSkill:
		return Will
	case Hide:
		return Dexterity
	case Intimidate:
		return Charisma
	case Jump:
		return Strength
	case Knowledge:
		return Intelligence
	case Listen:
		return Wisdom
	case MoveSilently:
		return Dexterity
	case OpenLock:
		return Dexterity
	case Perform:
		return Charisma
	case Psicraft:
		return Intelligence
	case Profession:
		return Wisdom
	case Ride:
		return Dexterity
	case Search:
		return Intelligence
	case SenseMotive:
		return Wisdom
	case SleightOfHand:
		return Dexterity
	case SpeakLanguage:
		return types.Unknown
	case Spellcraft:
		return Intelligence
	case Spot:
		return Wisdom
	case Survival:
		return Wisdom
	case Swim:
		return Strength
	case Tumble:
		return Dexterity
	case UseMagicDevice:
		return Charisma
	case UsePsionicDevice:
		return Charisma
	case UseRope:
		return Dexterity
	}
	return types.Unknown
}
