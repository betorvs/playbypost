package pfd20

import (
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/sys/library"
)

func GenPFD20Random(name string, r *rpg.RPGSystem, lib *library.Library) (*PathfinderCharacter, error) {
	p := New(name, r)
	// static values for now
	// human ancestry and background warrior and class fighter
	// 18, 12, 16, 10, 12, 10
	_ = p.AddAbility(base.Ability{Name: "Strength", Value: 18}, lib)
	_ = p.AddAbility(base.Ability{Name: "Dexterity", Value: 12}, lib)
	_ = p.AddAbility(base.Ability{Name: "Constitution", Value: 16}, lib)
	_ = p.AddAbility(base.Ability{Name: "Intelligence", Value: 10}, lib)
	_ = p.AddAbility(base.Ability{Name: "Wisdom", Value: 12}, lib)
	_ = p.AddAbility(base.Ability{Name: "Charisma", Value: 10}, lib)
	p.Ancestry = "Human"
	p.Background = "Warrior"
	p.Class = "Fighter"
	return p, nil
}
