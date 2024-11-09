package types

import "fmt"

const Unknown string = "unknown"

type Actions int

const (
	NoAction Actions = iota
	CheckAbility
	CheckSkill
	DoSingleAttack
	DoTotalAttack
	DoSpecialAttack
	CastSpell
	WearMagicItem
	KeepItem
	UseMagicItem
	UseSpecialAbility
	DoTotalDefense
)

func (a Actions) String() string {
	switch a {
	case CheckAbility:
		return "ability check"
	case CheckSkill:
		return "skill check"
	case DoSingleAttack:
		return "single attack"
	case DoTotalAttack:
		return "total attack"
	case DoSpecialAttack:
		return "special attack"
	case CastSpell:
		return "cast spell"
	case WearMagicItem:
		return "wear magic item"
	case KeepItem:
		return "keep item"
	case UseMagicItem:
		return "use magic item"
	case UseSpecialAbility:
		return "use special ability"
	case DoTotalDefense:
		return "total defense"
	}
	return "no action"
}

func ActAtoi(a string) (Actions, error) {
	switch a {
	case "ability check":
		return CheckAbility, nil
	case "skill check":
		return CheckSkill, nil
	case "attack", "single attack":
		return DoSingleAttack, nil
	case "total attack":
		return DoTotalAttack, nil
	case "special attack":
		return DoSpecialAttack, nil
	case "cast spell":
		return CastSpell, nil
	case "wear magic item":
		return WearMagicItem, nil
	case "keep item":
		return KeepItem, nil
	case "use magic item":
		return UseMagicItem, nil
	case "use special ability":
		return UseSpecialAbility, nil
	case "total defense":
		return DoTotalDefense, nil
	}
	return NoAction, fmt.Errorf("invalid")
}

type Effect int

const (
	NoneEffect Effect = iota
	DamageEffect
	HealEffect
	ChangeConditionEffect
)

func (e Effect) String() string {
	switch e {
	case NoneEffect:
		return "none"
	case DamageEffect:
		return "damage"
	case HealEffect:
		return "heal"
	case ChangeConditionEffect:
		return "change condition"
	}
	return Unknown
}
