package mechanism

import (
	"fmt"

	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/diceroll"
	"github.com/betorvs/playbypost/core/types"
)

type Actors int

const (
	PC Actors = iota
	NPC
)

func (a Actors) String() string {
	switch a {
	case PC:
		return "pc"
	case NPC:
		return "npc"
	}
	return types.Unknown
}

type Command interface {
	Call() error
	Undo() error
}

type CombatCommand struct {
	combat        rules.Combater
	action        types.Actions
	parameter     any
	desiredEffect types.Effect
	previousState types.State
	state         types.State
	succeeded     bool
	value         int
	dice          diceroll.DiceRoller
}

func (p *CombatCommand) Call() {
	value, effect, resultState, err := p.combat.Act(p.action, p.dice, p.parameter)
	if err != nil {
		fmt.Println("Combat Command Act Call error", err.Error())
		p.succeeded = false
	}
	// registering data
	p.previousState = p.combat.CopySingleTargetState()
	p.value = value
	p.desiredEffect = effect

	if value != 0 {
		// calling effect to apply it
		status, state, err := p.combat.Effect(effect, value, resultState, p.dice)
		if err != nil {
			fmt.Println("Combat Command Effect Call error", err.Error())
			p.succeeded = false
		}
		p.succeeded = status
		p.state = state
		if p.succeeded {
			fmt.Printf("You Hit! Result: your target %s is %s \n", p.combat.CopySingleTargetName(), p.state)
			return
		}
	}

	fmt.Println("Your act failed")
}

func (p *CombatCommand) Undo() {

	// damage
	if p.desiredEffect == types.DamageEffect {
		status, _, err := p.combat.Effect(types.HealEffect, p.value, p.state, p.dice)
		if err != nil {
			fmt.Println("Combat Command Effect damage Undo error", err.Error())
			p.succeeded = false
		}
		p.succeeded = status
	}

	// change condition
	if p.previousState != p.state {
		_, _, err := p.combat.Effect(types.ChangeConditionEffect, 0, p.previousState, p.dice)
		if err != nil {
			fmt.Println("Combat Command Effect change condition Undo error", err.Error())
			p.succeeded = false
		}
		p.succeeded = true
	}

	p.state = p.previousState
	if p.succeeded {
		fmt.Println("Restoring previous action", p.action, "and state", p.state)
		return
	}
	fmt.Println("Undo failed")
}

func NewCombatCommand(combat rules.Combater, action types.Actions, parameter any, dice diceroll.DiceRoller) *CombatCommand {
	return &CombatCommand{
		combat:    combat,
		action:    action,
		parameter: parameter,
		dice:      dice,
	}
}
