package rules

import (
	"fmt"

	"github.com/betorvs/playbypost/core/sys/diceroll"
	"github.com/betorvs/playbypost/core/types"
)

type TargetOption int

const (
	PCtoNPCMonster TargetOption = iota
	NPCtoPC
	PCtoPC
	NPCtoNPC
)

type Combater interface {
	Act(action types.Actions, dice diceroll.DiceRoller, parameter ...any) (int, types.Effect, types.State, error)
	Effect(effect types.Effect, value int, state types.State, dice diceroll.DiceRoller) (bool, types.State, error)
	CopySingleTargetState() types.State
	CopySingleTargetName() string
}

/*
Combat struct defines a turn of a combat and run a action between to CombatParticipants
*/
type Combat struct {
	Attacker        *CombatParticipant
	TargetOption    TargetOption
	SingleTarget    *CombatParticipant
	MultipleTargets []CombatParticipant
}

func NewCombatSingle(attacker *CombatParticipant, single *CombatParticipant) *Combat {
	var targetOption TargetOption
	switch {
	case attacker.ActorPC != nil && single.ActorPC != nil:
		targetOption = PCtoPC
	case attacker.ActorPC != nil:
		targetOption = PCtoNPCMonster
	case attacker.ActorNPC != nil:
		targetOption = NPCtoPC
	}
	c := Combat{
		Attacker:     attacker,
		SingleTarget: single,
		TargetOption: targetOption,
	}

	return &c
}

func NewCombatMultiple(Attacker *CombatParticipant, multiple ...CombatParticipant) *Combat {
	var targetOption TargetOption
	switch {
	case Attacker.ActorPC != nil && multiple[0].ActorPC != nil:
		targetOption = PCtoPC
	case Attacker.ActorPC != nil:
		targetOption = PCtoNPCMonster
	case Attacker.ActorNPC != nil:
		targetOption = NPCtoPC
	}

	c := Combat{
		Attacker:     Attacker,
		TargetOption: targetOption,
	}
	// need to validate it
	if len(multiple) == 1 {
		c.SingleTarget = &multiple[0]
	} else {
		c.MultipleTargets = multiple
	}

	return &c
}

func (c *Combat) Act(action types.Actions, dice diceroll.DiceRoller, parameter ...any) (int, types.Effect, types.State, error) {
	switch action {
	case types.DoAttack:
		switch c.TargetOption {
		case PCtoNPCMonster:
			if c.Attacker.ActorPC.totalDefense {
				// removing total defense
				c.Attacker.ActorPC.RemoveTotalDefense()
			}
			armorClass := c.SingleTarget.ActorNPC.ArmorClass
			if c.SingleTarget.ActorNPC.totalDefense {
				// https://www.dandwiki.com/wiki/SRD:Standard_Actions#Total_Defense
				armorClass += 4
			}
			diceResult, _, _ := dice.DiceRoll(c.Attacker.ActorPC.Name, "1d20")
			// fmt.Println("pc", c.Attacker)
			attackindex := c.Attacker.ActorPC.PreferenceAttackIndex
			if parameter[0].(int) != c.Attacker.ActorPC.PreferenceAttackIndex {
				attackindex = parameter[0].(int)
			}
			kind := c.Attacker.ActorPC.AttackOption[attackindex].Type
			damage, _ := c.Attacker.ActorPC.AttackOption[attackindex].Attack(diceResult, c.Attacker.ActorPC.CalcAttackTotal(kind), armorClass, dice)
			if damage.Successful {
				return damage.Damage, types.DamageEffect, c.SingleTarget.State, nil
			}
			return 0, types.NoneEffect, c.SingleTarget.State, nil
		case NPCtoPC:
			if c.Attacker.ActorNPC.totalDefense {
				// removing total defense
				c.Attacker.ActorNPC.RemoveTotalDefense()
			}
			armorClass := c.SingleTarget.ActorPC.ArmorClass
			if c.SingleTarget.ActorPC.totalDefense {
				// https://www.dandwiki.com/wiki/SRD:Standard_Actions#Total_Defense
				armorClass += 4
			}
			attackindex := c.Attacker.ActorNPC.PreferenceAttackIndex
			if parameter[0].(int) != c.Attacker.ActorNPC.PreferenceAttackIndex {
				attackindex = parameter[0].(int)
			}
			diceResult, _, _ := dice.DiceRoll(c.Attacker.ActorNPC.Name, "1d20")
			damage, _ := c.Attacker.ActorNPC.AttackOption[attackindex].Attack(diceResult, c.Attacker.ActorNPC.BaseAttack, armorClass, dice)
			if damage.Successful {
				return damage.Damage, types.DamageEffect, c.SingleTarget.State, nil
			}
			return 0, types.NoneEffect, c.SingleTarget.State, nil
		case PCtoPC:
			if c.Attacker.ActorPC.totalDefense {
				// removing total defense
				c.Attacker.ActorPC.RemoveTotalDefense()
			}
			armorClass := c.SingleTarget.ActorPC.ArmorClass
			if c.SingleTarget.ActorPC.totalDefense {
				// https://www.dandwiki.com/wiki/SRD:Standard_Actions#Total_Defense
				armorClass += 4
			}
			attackindex := c.Attacker.ActorPC.PreferenceAttackIndex
			if parameter[0].(int) != c.Attacker.ActorPC.PreferenceAttackIndex {
				attackindex = parameter[0].(int)
			}
			diceResult, _, _ := dice.DiceRoll(c.Attacker.ActorPC.Name, "1d20")
			kind := c.Attacker.ActorPC.AttackOption[attackindex].Type
			damage, _ := c.Attacker.ActorPC.AttackOption[attackindex].Attack(diceResult, c.Attacker.ActorPC.CalcAttackTotal(kind), armorClass, dice)
			if damage.Successful {
				return damage.Damage, types.DamageEffect, c.SingleTarget.State, nil
			}
			return 0, types.NoneEffect, c.SingleTarget.State, nil
		}

	case types.DoSpecialAttack:
		// require 2 parameters
		var specialattack SpecialAttacks
		v, err := specialattack.Atoi(parameter[0].(string))
		if err != nil {
			return 0, types.NoneEffect, types.Alive, nil
		}
		var diceResult int
		switch c.TargetOption {
		case PCtoNPCMonster:
			diceResult, _, _ = dice.DiceRoll(c.Attacker.ActorPC.Name, "1d20")
		case NPCtoPC:
			diceResult, _, _ = dice.DiceRoll(c.Attacker.ActorNPC.Name, "1d20")
		case PCtoPC:
			diceResult, _, _ = dice.DiceRoll(c.Attacker.ActorPC.Name, "1d20")
		}

		a, b := c.SpecialAttack(v, diceResult, parameter[1].(string))
		fmt.Println("special attack", v, "values are", a, b)

		return 0, types.NoneEffect, types.Alive, nil
	case types.MonsterSpecialAttack:
		// only for c.TargetOption == NPCtoPC
		return 0, types.NoneEffect, types.Alive, nil
	case types.CastSpell:
		return 0, types.NoneEffect, types.Alive, nil
	case types.ActivateMagicItem:
		return 0, types.NoneEffect, types.Alive, nil
	case types.UseSpecialAbility:
		return 0, types.NoneEffect, types.Alive, nil
	case types.DoTotalDefense:
		switch c.TargetOption {
		case PCtoNPCMonster, PCtoPC:
			c.Attacker.ActorPC.SetTotalDefense()
		case NPCtoPC:
			c.Attacker.ActorNPC.SetTotalDefense()
		}
		return 0, types.NoneEffect, types.Alive, nil
	}
	return 0, types.NoneEffect, types.Alive, nil
}

/*
Useful code
switch c.TargetOption {
		case PCtoNPCMonster:
		case NPCtoPC:
		case PCtoPC:
		}
*/

func (c *Combat) Effect(effect types.Effect, value int, state types.State, dice diceroll.DiceRoller) (bool, types.State, error) {
	switch effect {
	case types.DamageEffect:
		switch c.TargetOption {
		case PCtoNPCMonster:
			c.SingleTarget.ActorNPC.HitPoints = c.SingleTarget.ActorNPC.HitPoints - value
			if c.SingleTarget.ActorNPC.HitPoints <= 0 {
				return true, types.Dead, nil
			}
			return true, types.Alive, nil
		case NPCtoPC:
			c.SingleTarget.ActorPC.HitPoints = c.SingleTarget.ActorPC.HitPoints - value
			if c.SingleTarget.ActorPC.HitPoints <= 0 {
				return true, types.Dead, nil
			}
			return true, types.Alive, nil
		}

	case types.HealEffect:
		switch c.TargetOption {
		case PCtoNPCMonster:
			c.SingleTarget.ActorNPC.HitPoints = c.SingleTarget.ActorNPC.HitPoints + value
			switch {
			case c.SingleTarget.ActorNPC.HitPoints >= 0:
				return true, types.Alive, nil
			case c.SingleTarget.ActorNPC.HitPoints <= -1 && c.SingleTarget.ActorNPC.HitPoints >= -9:
				return true, types.Unconscious, nil
			case c.SingleTarget.ActorNPC.HitPoints < -10:
				return true, types.Dead, nil
			}
		case PCtoPC:
			c.SingleTarget.ActorPC.HitPoints = c.SingleTarget.ActorPC.HitPoints + value
			switch {
			case c.SingleTarget.ActorPC.HitPoints >= 0:
				return true, types.Alive, nil
			case c.SingleTarget.ActorPC.HitPoints <= -1 && c.SingleTarget.ActorPC.HitPoints >= -9:
				return true, types.Unconscious, nil
			case c.SingleTarget.ActorPC.HitPoints < -10:
				return true, types.Dead, nil
			}
		}

		return false, types.Unconscious, nil
	case types.ChangeConditionEffect:
		switch c.TargetOption {
		case PCtoNPCMonster:
			err := c.SingleTarget.ActorNPC.ChangeCondition(state)
			if err != nil {
				return false, c.SingleTarget.ActorNPC.state, err
			}
			return true, c.SingleTarget.ActorNPC.state, nil
		case NPCtoPC:
			err := c.SingleTarget.ActorPC.ChangeCondition(state)
			if err != nil {
				return false, c.SingleTarget.ActorPC.state, err
			}
			return true, c.SingleTarget.ActorPC.state, nil
		}
		return false, types.Alive, nil
	}
	return false, types.Alive, nil
}

func (c Combat) CopySingleTargetState() types.State {
	return c.SingleTarget.State
}

func (c Combat) CopySingleTargetName() string {
	switch c.TargetOption {
	case PCtoNPCMonster:
		return c.SingleTarget.ActorNPC.Name
	case NPCtoPC:
		return c.SingleTarget.ActorPC.Name
	case PCtoPC:
		return c.SingleTarget.ActorPC.Name
	}
	return c.SingleTarget.ActorPC.Name
}

/*
CombatParticipant struct returns all possible participants from a combat
*/
type CombatParticipant struct {
	ActorPC  *Player
	ActorNPC *Monster
	State    types.State
}

func NewPCCombatParticipant(PC *Player) *CombatParticipant {
	return &CombatParticipant{
		ActorPC: PC,
		State:   PC.state,
	}
}

func NewNPCCombatParticipant(NPC *Monster) *CombatParticipant {
	return &CombatParticipant{
		ActorNPC: NPC,
		State:    NPC.state,
	}
}

func (c CombatParticipant) String() string {
	switch {
	case c.ActorNPC != nil:
		return fmt.Sprintf("Name: %s, Title: %s, CA: %d\n", c.ActorNPC.Name, c.ActorNPC.Title, c.ActorNPC.ArmorClass)
	case c.ActorPC != nil:
		return fmt.Sprintf("Name: %s, Level: %d, CA: %d, AttackOptions: %+v\n", c.ActorPC.Name, c.ActorPC.Level, c.ActorPC.ArmorClass, c.ActorPC.AttackOption[0])
	}
	return types.Unknown
}
