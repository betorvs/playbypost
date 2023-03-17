package rules

import "math"

const (
	Strength     string = "strength"
	Dexterity    string = "dexterity"
	Constitution string = "constitution"
	Intelligence string = "intelligence"
	Wisdom       string = "wisdom"
	Charisma     string = "charisma"
)

func CalcAbilityModifier(attr int) int {
	result := math.Floor((float64(attr) - 10) / 2)
	return int(result)
}

type Ability struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}
