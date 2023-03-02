package types

const Unknown string = "unknown"

type Actions int

const (
	DoAttack Actions = iota
	DoSpecialAttack
	MonsterSpecialAttack
	CastSpell
	ActivateMagicItem
	UseSpecialAbility
	DoTotalDefense
)

func (a Actions) String() string {
	switch a {
	case DoAttack:
		return "attack"
	case MonsterSpecialAttack:
		return "special attack"
	case DoSpecialAttack:
		return "monster special attack"
	case CastSpell:
		return "cast spell"
	case ActivateMagicItem:
		return "activate magic item"
	case UseSpecialAbility:
		return "use special ability"
	case DoTotalDefense:
		return "total defense"
	}
	return Unknown
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
