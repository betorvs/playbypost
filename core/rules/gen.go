package rules

import (
	"math/rand"
	"slices"
	"time"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/d10hm"
)

func GenD10Random(n string, rpgSystem *rpg.RPGSystem) (*Creature, error) {
	physical := "physical"
	mental := "mental"
	social := "social"
	person1 := NewCreature(n, rpgSystem)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	{
		a := []int{2, 3, 3}

		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		for k, v := range rpgSystem.Ability.Grouped[physical] {
			ability := Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		a := []int{2, 2, 3}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		for k, v := range rpgSystem.Ability.Grouped[mental] {
			ability := Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		a := []int{2, 2, 2}
		rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })

		for k, v := range rpgSystem.Ability.Grouped[social] {
			ability := Ability{
				Name:  v,
				Value: a[k],
			}
			err := person1.AddAbility(ability)
			if err != nil {
				return person1, err
			}
		}
	}
	{
		b := []int{3, 3, 2, 2, 1}
		rand.Shuffle(len(b), func(i, j int) { b[i], b[j] = b[j], b[i] })
		count := 0
		for _, v := range rpgSystem.Skill.Grouped[physical] {
			if slices.Contains(rpgSystem.Skill.Tags[v], "auto-gen") {
				s := Skill{
					Name:  v,
					Value: b[count],
					Base:  rpgSystem.GetSkillBase(v),
				}
				err := person1.AddSkill(s)
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
		for _, v := range rpgSystem.Skill.Grouped[social] {
			if slices.Contains(rpgSystem.Skill.Tags[v], "auto-gen") {
				s := Skill{
					Name:  v,
					Value: b[count],
					Base:  rpgSystem.GetSkillBase(v),
				}
				err := person1.AddSkill(s)
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
		for _, v := range rpgSystem.Skill.Grouped[mental] {
			if slices.Contains(rpgSystem.Skill.Tags[v], "auto-gen") {
				s := Skill{
					Name:  v,
					Value: b[count],
					Base:  rpgSystem.GetSkillBase(v),
				}
				err := person1.AddSkill(s)
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
		d10hmExtended := d10hm.NewWithValuesExtended(resolve, composture, dexterity, wits)
		// longsword := weapons[0]
		d10hmExtended.Weapon["longsword"] = d10hm.Weapon{
			Name:        "longsword",
			Value:       3,
			Description: "a longsword",
		}
		person1.Extension = d10hmExtended
	}

	return person1, nil
}
