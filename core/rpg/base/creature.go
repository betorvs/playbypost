package base

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/sys/library"
)

const (
	DestroyedError string = "cannot destroy a creature that has been destroyed"
	AbilityInvalid string = "ability not valid"
	SkillInvalid   string = "skill not valid"
)

type Creature struct {
	Name      string         `json:"name"`
	Abilities Abilities      `json:"abilities"`
	Skills    Skills         `json:"skills"`
	RPG       *rpg.RPGSystem `json:"rpg"`
	destroyed bool
}

func (a Abilities) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Abilities) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Abilities map[string]Ability

type Ability struct {
	Name        string
	DisplayName string
	Value       int
}

func (a Skills) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Skills) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Skills map[string]Skill

type Skill struct {
	Name        string
	DisplayName string
	Value       int
	Base        string
}

func NewCreature(n string, r *rpg.RPGSystem) *Creature {
	return &Creature{
		Name:      n,
		Abilities: make(map[string]Ability),
		Skills:    make(map[string]Skill),
		RPG:       r,
		destroyed: false,
	}
}

func (c *Creature) Destroy() error {
	if !c.destroyed {
		c.destroyed = true
		return nil
	}
	return errors.New(DestroyedError)
}

func (c *Creature) IsDead() bool {
	return c.destroyed
}

func RestoreCreature() *Creature {
	return &Creature{
		Abilities: make(map[string]Ability),
		Skills:    make(map[string]Skill),
		destroyed: false,
	}
}

func (c *Creature) AddAbility(a Ability, lib *library.Library) error {
	if slices.Contains(lib.Ability.List, a.Name) {
		c.Abilities[a.Name] = a
		return nil
	}
	return errors.New(AbilityInvalid)
}

func (c *Creature) AddSkill(s Skill, lib *library.Library) error {
	if slices.Contains(lib.Skill.List, s.Name) {
		c.Skills[s.Name] = s
		return nil
	}
	return errors.New(SkillInvalid)
}
