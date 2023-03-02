package rules

import (
	"errors"

	"github.com/betorvs/playbypost/core/types"
)

type CoinsType int

const (
	Copper CoinsType = iota
	Silver
	Gold
	Platinum
)

var ErrExchange = errors.New("cannot exchange for a more valued currency")

func (c CoinsType) String() string {
	switch c {
	case Copper:
		return "copper"
	case Silver:
		return "silver"
	case Gold:
		return "gold"
	case Platinum:
		return "platinum"
	}
	return types.Unknown
}

func (c CoinsType) ShortName() string {
	switch c {
	case Copper:
		return "cp"
	case Silver:
		return "sp"
	case Gold:
		return "gp"
	case Platinum:
		return "pp"
	}
	return types.Unknown
}

func ExchangeRates(value int, sourceType CoinsType, toType CoinsType) (int, error) {
	switch sourceType {
	case Copper:
		switch toType {
		case Copper:
			return value, nil
		case Silver:
			return 0, ErrExchange
		case Gold:
			return 0, ErrExchange
		case Platinum:
			return 0, ErrExchange
		}
	case Silver:
		switch toType {
		case Copper:
			return value * 10, nil
		case Silver:
			return value, nil
		case Gold:
			return 0, ErrExchange
		case Platinum:
			return 0, ErrExchange
		}
	case Gold:
		switch toType {
		case Copper:
			return value * 100, nil
		case Silver:
			return value * 10, nil
		case Gold:
			return value, nil
		case Platinum:
			return 0, ErrExchange
		}
	case Platinum:
		switch toType {
		case Copper:
			return value * 1000, nil
		case Silver:
			return value * 100, nil
		case Gold:
			return value * 10, nil
		case Platinum:
			return value, nil
		}
	}
	return 0, nil
}
