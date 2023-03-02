package mechanism

import (
	"fmt"
	"sort"
	"strings"
)

const tiebreaker = 500

type RoundControl struct {
	pointer      int
	rollsOrder   []float64
	participants map[string]float64
	order        map[int]string
}

// NewRoundControl factory function creates a RoundControl with the first participant that rolled his initiative
func NewRoundControl(s string, i int) *RoundControl {
	p := make(map[string]float64)
	sl := []float64{tiebreakerCall(i, 0)}
	// recover value from slice of floats
	p[s] = sl[0]
	enc := RoundControl{
		pointer:      0,
		rollsOrder:   sl,
		participants: p,
		order:        calculateOrder(sl, p),
	}
	return &enc
}

// Append function will keep RoundControl ordered after each participant being added
func (e *RoundControl) Append(s string, value int) error {
	// cannot add participants with the same "name"
	if _, ok := e.participants[s]; ok {
		return fmt.Errorf("participant already in %s", s)
	}
	// add validation for 3th participant with same rolled initiative
	size := len(e.rollsOrder)
	// if size >= 2 {
	// 	fmt.Println("checking")
	// }
	//
	ns := make([]float64, size+1)
	copy(ns, e.rollsOrder)
	i := tiebreakerCall(value, size)
	for k, v := range ns {
		// if i > v => it should take place and move v to next
		// else assign v to the latest value
		// fmt.Println("value", i)
		if i > v {
			ns[k] = i
			i = v
		} else {
			ns[k] = v
		}
	}
	e.rollsOrder = ns
	// add new participant to map
	e.participants[s] = tiebreakerCall(value, size)
	// recalculate ordered map
	e.order = calculateOrder(e.rollsOrder, e.participants)
	// fmt.Println("rollsOrder", e.rollsOrder, "participants", e.participants, "order", e.order)
	return nil
}

func (e *RoundControl) Next() string {
	next := e.order[e.pointer]
	// increase pointer for next turn
	e.pointer++
	if e.pointer >= len(e.rollsOrder) {
		// if pointer reachs the end of slice
		// re set it to 0
		e.pointer = 0
	}
	fmt.Println("Hey, It's your turn", next)
	return next
}

func (e RoundControl) String() string {
	size := len(e.rollsOrder)
	s := make([]string, size)
	for k, v := range e.order {
		s[k] = fmt.Sprintf("%d - %s", k, v)
	}
	return strings.Join(s, ", ")
}

func (e *RoundControl) Remove(s string) {
	v := e.participants[s]
	// search in rollsOrder
	i := sort.Search(len(e.rollsOrder), func(i int) bool { return e.rollsOrder[i] <= v })
	// remove it from slice of ints
	e.rollsOrder = append(e.rollsOrder[:i], e.rollsOrder[i+1:]...)
	// delete(e.initiativeRolls, v)
	delete(e.participants, s)
	e.order = calculateOrder(e.rollsOrder, e.participants)
}

func calculateOrder(s []float64, m map[string]float64) map[int]string {
	ordered := make(map[int]string, len(s))
	for k, v := range m {
		i := sort.Search(len(s), func(i int) bool { return s[i] <= v })
		// fmt.Println("1 - index", i, "participant", k, "value", v)
		// if _, ok := ordered[i]; ok {
		// 	fmt.Println("2 - index", i, "participant", k, "value", v)
		// 	// i++
		// }
		ordered[i] = k
	}
	// fmt.Println("len", len(ordered))
	return ordered
}

func tiebreakerCall(value, subtractor int) float64 {
	v := (tiebreaker - float64(subtractor)) / 1000
	return float64(value) + v
}
