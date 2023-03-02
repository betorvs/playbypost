package rules

// Dwarves, Hill
// Elves, High
// Gnomes, Rock
// Half-Elves
// Half-Orcs
// Halflings, Lightfoot
// Humans

/* Race
CreatureType
Size
Basic speed
benefits: extra feat, extra skill
languages
weapon familiarity
favored class
special abilities
bonus
 skills
 savings
 attack
 armor class
abilities scores

*/

type Race struct {
	Type  Creatures
	Size  CreaturesSizes
	Speed int
}
