//go:build unit

package initiative

import (
	"fmt"
	"sort"
	"testing"
)

func TestInitiative(t *testing.T) {
	party := Participants{
		{Name: "test1", Bonus: 5, Result: 10},
		{Name: "test2", Bonus: 5, Result: 5},
		{Name: "test3", Bonus: 5, Result: 15},
	}
	sort.Sort(sort.Reverse(party))
	i := Initiative{
		Name:         "test",
		Position:     -1,
		Participants: party,
	}
	fmt.Printf("party %v \n", party)
	if i.Current() != "test3" {
		t.Errorf("expected test3 got %s", i.Current())
	}
	// first next will return 0
	if i.Next() != 0 {
		t.Errorf("expected 0 got %d", i.Position)
	}
	// move to next participant
	i.Next()
	if i.Current() != "test1" {
		t.Errorf("expected test1 got %s", i.Current())
	}
	// move to next participant
	i.Next()
	if i.Position != 2 {
		t.Errorf("expected 2 got %d", i.Position)
	}
	// move to next participant
	i.Next()
	if i.Current() != "test3" {
		t.Errorf("expected test3 got %s", i.Current())
	}

	if i.Position != 0 {
		t.Errorf("expected 0 got %d", i.Position)
	}
	// move to next participant
	i.Next()
	if i.Current() != "test1" {
		t.Errorf("expected test1 got %s", i.Current())
	}
}
