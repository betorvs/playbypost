package d10hm

import (
	"math/rand"
	"slices"
	"time"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/sys/library"
)

// func (c *StorytellingCharacter) GenRandom(name string, rpgSystem *rpg.RPGSystem) (rules.RolePlaying, error) {
// 	return GenD10Random(name, rpgSystem)
// }

func GenD10Random(n string, rpgSystem *rpg.RPGSystem, lib *library.Library) (*StorytellingCharacter, error) {
	physical := "physical"
	mental := "mental"
	social := "social"
	person1 := newCreature(n, rpgSystem)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	{
		a := []int{2, 3, 3}

		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for k, v := range lib.Ability.Grouped[physical] {
			ability := base.Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability, lib)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		a := []int{2, 2, 3}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		for k, v := range lib.Ability.Grouped[mental] {
			ability := base.Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability, lib)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		a := []int{2, 2, 2}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		for k, v := range lib.Ability.Grouped[social] {
			ability := base.Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability, lib)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		b := []int{3, 3, 2, 2, 1}
		rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
		count := 0
		for _, v := range lib.Skill.Grouped[physical] {
			if slices.Contains(lib.Skill.Tags[v], "auto-gen") {
				s := base.Skill{
					Name:  v,
					Value: b[count],
					Base:  lib.GetSkillBase(v),
				}
				err := person1.AddSkill(s, lib)
				if err != nil {
					return person1, err
				}
				count++
			}

		}
	}
	{
		b := []int{3, 2, 2}
		rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
		count := 0
		for _, v := range lib.Skill.Grouped[social] {
			if slices.Contains(lib.Skill.Tags[v], "auto-gen") {
				s := base.Skill{
					Name:  v,
					Value: b[count],
					Base:  lib.GetSkillBase(v),
				}
				err := person1.AddSkill(s, lib)
				if err != nil {
					return person1, err
				}
				count++
			}

		}
	}

	{
		b := []int{2, 2}
		count := 0
		for _, v := range lib.Skill.Grouped[mental] {
			if slices.Contains(lib.Skill.Tags[v], "auto-gen") {
				s := base.Skill{
					Name:  v,
					Value: b[count],
					Base:  lib.GetSkillBase(v),
				}
				err := person1.AddSkill(s, lib)
				if err != nil {
					return person1, err
				}
				count++
			}
		}
	}

	switch rpgSystem.Name {
	case rpg.D10HM:
		resolve := person1.Abilities["resolve"].Value
		composture := person1.Abilities["composture"].Value
		dexterity := person1.Abilities["dexterity"].Value
		wits := person1.Abilities["wits"].Value
		d10hmExtended := newWithValuesExtended(resolve, composture, dexterity, wits)
		// longsword := weapons[0]
		d10hmExtended.Weapon.SetWeapon("longsword", "melee", 3, "a longsword")
		// d10hmExtended.Weapon["longsword"] = Weapon{
		// 	Name:        "longsword",
		// 	Value:       3,
		// 	Description: "a longsword",
		// }
		person1.D10Extented = d10hmExtended
	}

	return person1, nil
}

func newCreature(n string, r *rpg.RPGSystem) *StorytellingCharacter {
	return &StorytellingCharacter{
		Creature: base.Creature{
			Name:      n,
			Abilities: make(map[string]base.Ability),
			Skills:    make(map[string]base.Skill),
			RPG:       r,
		}}
}
