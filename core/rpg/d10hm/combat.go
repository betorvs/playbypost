package d10hm

func (c *StorytellingCharacter) Attack(kind, weapon string) int {
	// kind == skill
	switch kind {
	case "weaponry":
		strenght := c.Abilities["strenght"].Value
		weaponry := c.Skills["weaponry"].Value
		weaponValue, _, _ := c.WeaponBonus(weapon)
		return strenght + weaponry + weaponValue
	}
	return 0
}

func (c *StorytellingCharacter) DefenseBonus(kind string) int {
	switch kind {
	case "melee":
		return c.Armor + c.Defense
	case "ranged":
		return c.Armor
	}
	return 0
}

func (c *StorytellingCharacter) InitiativeBonus() (int, error) {
	return c.Initiative, nil
}
