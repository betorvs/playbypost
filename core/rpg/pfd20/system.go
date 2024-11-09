package pfd20

import (
	"fmt"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
)

// System is the Pathfinder 2.0 system
// https://pf2.d20pfsrd.com/

type PathfinderCharacter struct {
	base.Creature
	PFExtended
}

func New(n string, r *rpg.RPGSystem) *PathfinderCharacter {
	return &PathfinderCharacter{
		Creature:   *base.NewCreature(n, r),
		PFExtended: *newPFExtended(),
	}
}

func (c *PathfinderCharacter) Name() string {
	return c.Creature.Name
}

func (c *PathfinderCharacter) SetName(n string) error {
	if n == "" {
		return fmt.Errorf("name is empty")
	}
	c.Creature.Name = n
	return nil
}

func (c *PathfinderCharacter) RPGSystem() *rpg.RPGSystem {
	return c.Creature.RPG
}

func (c *PathfinderCharacter) Damage(v int) error {
	c.HitPoints = c.HitPoints - v
	return nil
}

func (c *PathfinderCharacter) HealthStatus() int {
	return c.HitPoints
}

func (c *PathfinderCharacter) InitiativeBonus() (int, error) {
	p := proficiencyRank(c.Proficiency[Perception].Level, c.Level)
	s := c.calcAbilityModifier(Wisdom)
	return p + s, nil
}

func (c *PathfinderCharacter) SetWeapon(name, kind string, value int, description string) {
	c.PFExtended.Weapon.SetWeapon(name, kind, value, description)
}

func (c *PathfinderCharacter) SetArmor(v int) {
	c.PFExtended.ArmorClassBonus = v
}

type PFExtended struct {
	Ancestry        string
	Background      string
	Class           string
	HitPoints       int
	ArmorClassBonus int
	Level           int
	Proficiency     map[string]Proficiency
	Weapon          base.Weapons
}

type Proficiency struct {
	On     string `json:"on"`
	Level  string `json:"level"`
	Source string `json:"source"`
}

func newPFExtended() *PFExtended {
	return &PFExtended{}
}

func (p PFExtended) getValues() map[string]interface{} {
	return map[string]interface{}{
		"ancestry":    p.Ancestry,
		"background":  p.Background,
		"class":       p.Class,
		"hit_points":  p.HitPoints,
		"armor_class": p.ArmorClassBonus,
		"level":       p.Level,
		"proficiency": p.Proficiency,
		"weapon":      p.Weapon,
	}
}

func (p PFExtended) DefenseBonus(s string) (int, error) {
	switch s {
	case Fortitude:
		p := proficiencyRank(p.Proficiency[Fortitude].Level, p.Level)
		return p, nil
	case Reflex:
		p := proficiencyRank(p.Proficiency[Reflex].Level, p.Level)
		return p, nil
	case Will:
		p := proficiencyRank(p.Proficiency[Will].Level, p.Level)
		return p, nil
	case ArmorClass:
		p := proficiencyRank(p.Proficiency[ArmorClass].Level, p.Level)
		return p, nil
	}
	return 0, nil
}

func (p PFExtended) WeaponBonus(s string) (int, string, error) {
	return 0, "", nil
}

// equipment

const (
	ProficiencyUntrained = "untrained"
	ProficiencyTrained   = "trained"
	ProficiencyExpert    = "expert"
	ProficiencyMaster    = "master"
	ProficiencyLegendary = "legendary"
)

// Proficiency is the level of proficiency + your level
func proficiencyRank(rank string, level int) int {
	switch rank {
	case ProficiencyUntrained:
		return 0
	case ProficiencyTrained:
		return level + 2
	case ProficiencyExpert:
		return level + 4
	case ProficiencyMaster:
		return level + 6
	case ProficiencyLegendary:
		return level + 8
	}
	return 0
}
