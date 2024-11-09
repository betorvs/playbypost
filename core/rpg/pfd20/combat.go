package pfd20

func (c *PathfinderCharacter) Attack(kind, weapon string) int {
	// kind == melee or ranged physical attack
	switch kind {
	case Melee:
		strenght := c.calcAbilityModifier(Strength)
		w, ok := c.PFExtended.Weapon.GetWeapon(weapon)
		if ok {
			p := proficiencyRank(c.PFExtended.Proficiency[w.Kind].Level, c.PFExtended.Level)
			return strenght + p
		}
		return strenght
	case Ranged:
		dexterity := c.calcAbilityModifier(Dexterity)
		w, ok := c.PFExtended.Weapon.GetWeapon(weapon)
		if ok {
			p := proficiencyRank(c.PFExtended.Proficiency[w.Kind].Level, c.PFExtended.Level)
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
		return c.PFExtended.ArmorClassBonus + dex + proficiencyRank(c.PFExtended.Proficiency[ArmorClass].Level, c.PFExtended.Level)
	case Ranged:
		return c.PFExtended.ArmorClassBonus + dex + proficiencyRank(c.PFExtended.Proficiency[ArmorClass].Level, c.PFExtended.Level)
	}
	return 0
}
