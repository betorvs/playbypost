package rules

import (
	"errors"

	"github.com/betorvs/playbypost/core/types"
)

var ErrSpecialAttacksTypes = errors.New("cannot convert string to SpecialAttacks")

type SpecialAttacks int

const (
	AidAnother SpecialAttacks = iota
	BullRush
	Charge
	Disarm
	Feint
	Grapple
	ManufacturedandNaturalWeaponFighting
	MountedCombat
	NaturalWeapons
	Overrun
	Sunder
	ThrowSplashWeapon
	Trip
	TurnOrRebukeUndead
	TwoWeaponFighting
)

func (s SpecialAttacks) String() string {
	switch s {
	case AidAnother:
		return "aid another"
	case BullRush:
		return "bull rush"
	case Charge:
		return "charge"
	case Disarm:
		return "disarm"
	case Feint:
		return "feint"
	case Grapple:
		return "grapple"
	case ManufacturedandNaturalWeaponFighting:
		return "manufactured and natural weapon fighting"
	case MountedCombat:
		return "mounted combat"
	case NaturalWeapons:
		return "natural weapons"
	case Overrun:
		return "overrun"
	case Sunder:
		return "sunder"
	case ThrowSplashWeapon:
		return "throw splash weapon"
	case Trip:
		return "trip"
	case TurnOrRebukeUndead:
		return "turn or rebuke undead"
	case TwoWeaponFighting:
		return "two weapon fighting"
	}
	return types.Unknown
}

func (s SpecialAttacks) Atoi(a string) (SpecialAttacks, error) {
	switch a {
	case "aid another", "aidanother":
		return AidAnother, nil
	case "bull rush", "bullrush":
		return BullRush, nil
	case "charge":
		return Charge, nil
	case "disarm":
		return Disarm, nil
	case "feint":
		return Feint, nil
	case "grapple":
		return Grapple, nil
	case "manufactured and natural weapon fighting", "manufacturednaturalweaponfighting", "manufacturedandnaturalweaponfighting":
		return ManufacturedandNaturalWeaponFighting, nil
	case "mounted combat", "mountedcombat":
		return MountedCombat, nil
	case "natural weapons", "naturalweapons":
		return NaturalWeapons, nil
	case "overrun":
		return Overrun, nil
	case "sunder":
		return Sunder, nil
	case "throw splash weapon", "throwsplashweapon":
		return ThrowSplashWeapon, nil
	case "trip":
		return Trip, nil
	case "turn or rebuke undead", "turnorrebukeundead", "turnrebukeundead", "turn undead", "turnundead", "rebuke undead", "rebukeundead":
		return TurnOrRebukeUndead, nil
	case "two weapon fighting", "twoweaponfighting":
		return TwoWeaponFighting, nil
	}
	return -1, ErrSpecialAttacksTypes
}

/*
SpecialAttack method returns attack modifier and defense modifier
*/
func (c *Combat) SpecialAttack(attack SpecialAttacks, valueFromDice int, parameter string) (int, int) {
	switch attack {
	case AidAnother:
		// https://www.dandwiki.com/wiki/SRD:Aid_Another
		// requires attack test against CA 10
		// choice one, or attack or defense

		// PCtoPC or NPCtoNPC usually
		var succeeded bool
		switch c.TargetOption {
		case PCtoNPCMonster, PCtoPC:
			kind := c.Attacker.ActorPC.AttackOption[c.Attacker.ActorPC.PreferenceAttackIndex].Type
			succeeded = c.Attacker.ActorPC.AttackOption[c.Attacker.ActorPC.PreferenceAttackIndex].CheckAttack(valueFromDice, c.Attacker.ActorPC.CalcAttackTotal(kind), 10)
		case NPCtoPC, NPCtoNPC:
			succeeded = c.Attacker.ActorNPC.AttackOption[c.Attacker.ActorNPC.PreferenceAttackIndex].CheckAttack(valueFromDice, c.Attacker.ActorNPC.BaseAttack, 10)
		}
		if succeeded && parameter == "defense" {
			return 0, 2
		} else if succeeded && parameter == "attack" {
			return 2, 0
		}

		return 0, 0

	case BullRush:
		// https://www.dandwiki.com/wiki/SRD:Bull_Rush
		// movement
		// strenght check
		// provoke an opportunity attack

	case Charge:
		// https://www.dandwiki.com/wiki/SRD:Charge
		// require move more than 2 squares

		// attack extended

		return 2, -2

	case Disarm:
		// https://www.dandwiki.com/wiki/SRD:Disarm
		// movement
		// provoke an opportunity attack

	case Feint:
		// https://www.dandwiki.com/wiki/SRD:Feint
		// movement
		// skill check

	case Grapple:
		// https://www.dandwiki.com/wiki/SRD:Grapple
		// movement
		// provoke an opportunity attack

	case ManufacturedandNaturalWeaponFighting:
		// https://www.dandwiki.com/wiki/SRD:Manufactured_and_Natural_Weapon_Fighting
		// an monster attacking with weapon and natural attack
		// consider it as two weapon fighting
		return -5, 0

	case MountedCombat:
		// https://www.dandwiki.com/wiki/SRD:Mounted_Combat_(Rule)
		// fighting condition
		// skill check

		// ranged -4, 0 or -8, 0 if running
		// ranged can do a full attack
		return 1, 0

	case NaturalWeapons:
		// https://www.dandwiki.com/wiki/SRD:Natural_Weapons
		// just monster natural attacks

	case Overrun:
		// https://www.dandwiki.com/wiki/SRD:Overrun
		// https://www.dandwiki.com/wiki/SRD:Overrun
		// movement
		// provoke an opportunity attack

	case Sunder:
		// https://www.dandwiki.com/wiki/SRD:Sunder
		// movement
		// provoke an opportunity attack

	case ThrowSplashWeapon:
		// https://www.dandwiki.com/wiki/SRD:Throw_Splash_Weapon
		// ranged touch attack

	case Trip:
		// https://www.dandwiki.com/wiki/SRD:Trip
		// movement
		// imobilizacao

	case TurnOrRebukeUndead:
		// https://www.dandwiki.com/wiki/SRD:Turn_or_Rebuke_Undead
		// action against monsters

	case TwoWeaponFighting:
		// https://www.dandwiki.com/wiki/SRD:Two-Weapon_Fighting
		/*
			Circumstances									Primary Hand	Off Hand
			Normal penalties										–6			–10
			Off-hand weapon is light								–4			–8
			Two-Weapon Fighting feat								–4			–4
			Off-hand weapon is light and Two-Weapon Fighting feat	–2			–2
		*/

	}
	return 0, 0
}
