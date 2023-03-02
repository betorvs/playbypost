package rules

// Armor struct
type Armor struct {
	Name              string    `json:"name"`
	Title             string    `json:"title"`
	Kind              string    `json:"kind"`
	Cost              int       `json:"cost"`
	CoinType          CoinsType `json:"coin_type"`
	ArmorClass        int       `json:"armor_class"`
	DexterityModifier int       `json:"dexterity_modifier"`
	Stealth           bool      `json:"stealth"`
	Strength          int       `json:"strength"`
	Weight            int       `json:"weight"`
	Measure           string    `json:"measure"`
}
