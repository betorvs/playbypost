package d20e35

import "math"

const (
	Fortitude string = "fortitude"
	Reflex    string = "reflex"
	Will      string = "will"
)

type KindBonus int

const (
	Poor KindBonus = iota
	Average
	Good
)

func (k KindBonus) String() string {
	switch k {
	case Poor:
		return "poor"
	case Average:
		return "average"
	case Good:
		return "good"
	}
	return "poor"
}

type HitPointsMethod int

const (
	Half HitPointsMethod = iota
	Full
)

type HitDices int

const (
	D4 HitDices = iota
	D6
	D8
	D10
	D12
)

func (h HitDices) String() string {
	switch h {
	case D4:
		return "d4"
	case D6:
		return "d6"
	case D8:
		return "d8"
	case D10:
		return "d10"
	case D12:
		return "d12"
	}
	return "d4"
}

func DiceAtoi(s string) HitDices {
	switch s {
	case "d4":
		return D4

	case "d6":
		return D6

	case "d8":
		return D8

	case "d10":
		return D10

	case "d12":
		return D12
	}
	return D4
}

func (h HitDices) Value(m HitPointsMethod) int {
	switch h {
	case D4:
		if m == Full {
			return 4
		}
		return 2
	case D6:
		if m == Full {
			return 6
		}
		return 3
	case D8:
		if m == Full {
			return 8
		}
		return 4
	case D10:
		if m == Full {
			return 10
		}
		return 5
	case D12:
		if m == Full {
			return 12
		}
		return 6
	}
	return 2
}

type SavingThrows struct {
	Fortitude KindBonus
	Reflex    KindBonus
	Will      KindBonus
}

/*
baseAttackCalc func

https://nwn2.fandom.com/wiki/Base_attack_bonus
*/
func (c D20Extended) baseAttackCalc(kind KindBonus) int {
	// fmt.Println("base calc", c.Level, kind)
	var result float64
	switch kind {
	case Good:
		result = float64(c.Level)
	case Average:
		result = math.Floor(float64(c.Level) * 0.75)
	case Poor:
		result = math.Floor(float64(c.Level) * 0.5)
	}
	return int(result)
}

/*
SavingThrowsCalcBonus func
*/
func (c D20Extended) savingThrowsCalcBonus(s string) int {
	switch s {
	case Fortitude:
		return c.baseSaveCalc(c.SavingThrows.Fortitude)
	case Reflex:
		return c.baseSaveCalc(c.SavingThrows.Reflex)
	case Will:
		return c.baseSaveCalc(c.SavingThrows.Will)
	}
	return 0
}

/*
BaseSaveCalc func

https://nwn2.fandom.com/wiki/Base_save_bonus
*/
func (c D20Extended) baseSaveCalc(s KindBonus) int {
	var result float64
	switch s {
	case Good:
		result = math.Floor((float64(c.Level) / 2) + 2)
	case Poor:
		result = math.Floor((float64(c.Level) / 3))
	}

	return int(result)
}
