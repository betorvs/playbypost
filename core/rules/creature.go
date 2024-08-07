package rules

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/betorvs/playbypost/core/rpg"
)

const (
	DestroyedError string = "cannot destroy a creature that has been destroyed"
)

type Creature struct {
	Name      string             `json:"name"`
	Abilities Abilities          `json:"abilities"`
	Skills    Skills             `json:"skills"`
	RPG       *rpg.RPGSystem     `json:"rpg"`
	Extension rpg.ExtendedSystem `json:"extension"`
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

// func (c *Creature) UnmarshalJSON(b []byte) error {
// 	fmt.Printf("Creature UnmarshalJSON: %s\n", b)
// 	var v struct {
// 		Name      string          `json:"name"`
// 		Abilities Abilities       `json:"abilities"`
// 		Skills    Skills          `json:"skills"`
// 		RPG       *rpg.RPGSystem  `json:"rpg"`
// 		Extension json.RawMessage `json:"extension"`
// 		destroyed bool
// 	}
// 	err := json.Unmarshal(b, &v)
// 	if err != nil {
// 		return err
// 	}
// 	c.Name = v.Name
// 	c.Abilities = v.Abilities
// 	c.Skills = v.Skills
// 	c.RPG = v.RPG
// 	c.destroyed = v.destroyed
// 	switch v.RPG.Name {
// 	case rpg.D2035:
// 		var f d20e35.D20Extended

// 		if err := json.Unmarshal(v.Extension, &f); err != nil {
// 			return err
// 		}

// 		c.Extension = &f
// 	case rpg.D10HM:
// 		var f d10hm.D10Extented

// 		if err := json.Unmarshal(v.Extension, &f); err != nil {
// 			return err
// 		}

// 		c.Extension = &f

// 	case rpg.D10OS:
// 		var f d10os.D10Extented

// 		if err := json.Unmarshal(v.Extension, &f); err != nil {
// 			return err
// 		}

// 		c.Extension = &f
// 	}

// 	return nil
// }
