package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/rpg/d10hm"
	"github.com/betorvs/playbypost/core/rules"
)

type GeneratePlayer struct {
	PlayerID int    `json:"player_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	StageID  int    `json:"stage_id"`
	Name     string `json:"name"`
}

type Players struct {
	ID         int            `json:"id,omitempty"`
	Name       string         `json:"name"`
	StageID    int            `json:"stage_id"`
	PlayerID   int            `json:"player_id"`
	Abilities  map[string]int `json:"abilities"`
	Skills     map[string]int `json:"skills"`
	RPG        string         `json:"rpg"`
	Extensions Extensions     `json:"extensions"`
	Destroyed  bool           `json:"destroyed"`
}

type Extensions map[string]interface{}

func NewExtension() Extensions {
	return make(map[string]interface{})
}

func (a Extensions) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Extensions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed on extended")
	}

	return json.Unmarshal(b, &a)
}

func (a *Extensions) ConvertMap(m map[string]interface{}) {
	for k, v := range m {
		(*a)[k] = v
	}
}

func NewPlayer() *Players {
	return &Players{
		Abilities:  make(map[string]int),
		Skills:     make(map[string]int),
		Extensions: NewExtension(),
	}
}

func CreatureToPlayer(p *Players, c *base.Creature) {
	for k, v := range c.Abilities {
		key := k
		if v.DisplayName != "" && v.DisplayName != v.Name {
			key = v.DisplayName
		}
		p.Abilities[key] = v.Value
	}
	for k, v := range c.Skills {
		key := k
		if v.DisplayName != "" && v.DisplayName != v.Name {
			key = v.DisplayName
		}
		p.Skills[key] = v.Value
	}
}

func PlayerToCreature(p *Players, c *base.Creature) rules.RolePlaying {
	for k, v := range p.Abilities {
		c.Abilities[k] = base.Ability{
			Name:  k,
			Value: v,
		}
	}
	for k, v := range p.Skills {
		c.Skills[k] = base.Skill{
			Name:  k,
			Value: v,
		}
		if c.RPG.BaseSystem == rpg.D20 {
			c.Skills[k] = base.Skill{
				Name:  k,
				Value: v,
				Base:  c.RPG.GetSkillBase(k),
			}
		}
	}
	switch c.RPG.Name {
	case rpg.D10HM:
		character := d10hm.New(c.Name, c.RPG)
		character.Creature = *c
		return character
	}
	return nil
}

func GenerateRandomPlayer(name string, rpgSystem *rpg.RPGSystem) (rules.RolePlaying, error) {
	switch rpgSystem.Name {
	case rpg.D10HM:
		character, err := d10hm.GenD10Random(name, rpgSystem)
		return character, err
	}
	return nil, nil
}
