package types

import (
	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
)

type GeneratePlayer struct {
	PlayerID int    `json:"player_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	StageID  int    `json:"stage_id"`
	Name     string `json:"name"`
}

type Players struct {
	ID        int            `json:"id,omitempty"`
	Name      string         `json:"name"`
	StageID   int            `json:"stage_id"`
	PlayerID  int            `json:"player_id"`
	Abilities map[string]int `json:"abilities"`
	Skills    map[string]int `json:"skills"`
	RPG       string         `json:"rpg"`
	Extension map[string]int `json:"extension"`
	Destroyed bool           `json:"destroyed"`
}

func NewPlayer() *Players {
	return &Players{
		Abilities: make(map[string]int),
		Skills:    make(map[string]int),
		Extension: map[string]int{},
	}
}

func CreatureToPlayer(p *Players, c *rules.Creature) {
	for k, v := range c.Abilities {
		// db.logger.Info("abilities", "k", k, "v", v)
		key := k
		if v.DisplayName != "" && v.DisplayName != v.Name {
			key = v.DisplayName
		}
		p.Abilities[key] = v.Value
	}
	for k, v := range c.Skills {
		// tiltdb.logger.Info("skills", "k", k, "v", v)
		key := k
		if v.DisplayName != "" && v.DisplayName != v.Name {
			key = v.DisplayName
		}
		p.Skills[key] = v.Value
	}
}

func PlayerToCreature(p *Players, c *rules.Creature) {
	for k, v := range p.Abilities {
		c.Abilities[k] = rules.Ability{
			Name:  k,
			Value: v,
		}
	}
	for k, v := range p.Skills {
		c.Skills[k] = rules.Skill{
			Name:  k,
			Value: v,
		}
		if c.RPG.BaseSystem == rpg.D20 {
			c.Skills[k] = rules.Skill{
				Name:  k,
				Value: v,
				Base:  c.RPG.GetSkillBase(k),
			}
		}
	}
}
