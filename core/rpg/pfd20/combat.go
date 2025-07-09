package pfd20

func (c *PathfinderCharacter) Attack(kind, weapon string) int {
	// kind == melee or ranged physical attack
	switch kind {
	case Melee:
		strenght := c.calcAbilityModifier(Strength)
		w, ok := c.Weapon.GetWeapon(weapon)
		if ok {
			p := proficiencyRank(c.Proficiency[w.Kind].Level, c.Level)
			return strenght + p
		}
		return strenght
	case Ranged:
		dexterity := c.calcAbilityModifier(Dexterity)
		w, ok := c.Weapon.GetWeapon(weapon)
		if ok {
			p := proficiencyRank(c.Proficiency[w.Kind].Level, c.Level)
			return dexterity + p
		}
		return dexterity

	}
	return 0
}

func (c *PathfinderCharacter) DefenseBonus(kind string) int {
	dex := c.calcAbilityModifier(Dexterity)
	switch kind {
	case Melee:
		return c.ArmorClassBonus + dex + proficiencyRank(c.Proficiency[ArmorClass].Level, c.Level)
	case Ranged:
		return c.ArmorClassBonus + dex + proficiencyRank(c.Proficiency[ArmorClass].Level, c.Level)
	}
	return 0
}
