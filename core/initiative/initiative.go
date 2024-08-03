package initiative

import (
	"sort"

	"github.com/betorvs/playbypost/core/rpg"
)

/*
https://www.dandwiki.com/wiki/SRD:Initiative
At the start of a battle, each combatant makes an initiative check. An initiative check is a Dexterity check. Each character applies his or her Dexterity modifier to the roll. Characters act in order, counting down from highest result to lowest. In every round that follows, the characters act in the same order (unless a character takes an action that results in his or her initiative changing; see Special Initiative Actions).

If two or more combatants have the same initiative check result, the combatants who are tied act in order of total initiative modifier (highest first). If there is still a tie, the tied characters should roll again to determine which one of them goes before the other.
*/

type Initiative struct {
	Name         string
	Position     int
	Participants Participants
}

func (i *Initiative) Next() int {
	// fmt.Println("executing next ", i.Position, "len", i.Participants.Len())
	if i.Position >= i.Participants.Len()-1 || i.Position == -1 {
		i.Position = 0
		return 0
	}
	i.Position++

	return i.Position
}

func (i Initiative) NextInfo() int {
	// fmt.Println("executing next ", i.Position, "len", i.Participants.Len())
	if i.Position >= i.Participants.Len()-1 || i.Position == -1 {
		return 0
	}
	i.Position++

	return i.Position
}

func (i *Initiative) Remove(index int) {
	i.Participants = append(i.Participants[:index], i.Participants[index+1:]...)
}

type Participant struct {
	Name   string
	Bonus  int
	Result int
}

// https://medium.com/@briankworld/how-to-implement-custom-sorting-with-custom-structs-in-go-322e9c1d26b8
type Participants []Participant

// Implement the Len method required by sort.Interface
func (p Participants) Len() int {
	return len(p)
}

// Implement the Less method required by sort.Interface
func (p Participants) Less(i, j int) bool {
	if p[i].Result == p[j].Result {
		// fmt.Printf("result i %v j %v \n", p[i].Result, p[j].Result)
		// fmt.Printf("bonus i %v j %v \n", p[i].Bonus, p[j].Bonus)
		return p[i].Bonus < p[j].Bonus
	}
	return p[i].Result < p[j].Result
}

// Implement the Swap method required by sort.Interface
func (p Participants) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// NewInitiative creates a inititiave struct with all participants
// list should have user_id and initiative bonus in a map[string]int
// every initiative should have a name
// dice should come from rpg.InitiativeDice()
func NewInitiative(d rpg.Roll, list map[string]int, name string, dice string) Initiative {
	party := Participants{}
	for k, v := range list {
		if k != "" && v != 0 {
			roll, _ := d.FreeRoll(k, dice)
			p := Participant{
				Name:   k,
				Bonus:  v,
				Result: roll.Result,
			}
			party = append(party, p)
		}
	}
	sort.Sort(sort.Reverse(party))
	i := Initiative{
		Name:         name,
		Position:     -1,
		Participants: party,
	}

	return i
}
