package definitions

import "encoding/json"

const (
	ClassKind string = "Class"
	RaceKind  string = "Race"
)

type D20Definitions struct {
	Name         string         `json:"name"`
	Kind         string         `json:"kind"`   // Race
	System       string         `json:"system"` // d20-3.5
	Size         string         `json:"size"`   // small, medium
	Feats        []string       `json:"race_feats"`
	Description  string         `json:"description"`
	AbilityBonus map[string]int `json:"ability_bonus"`
}

func LoadD20DefinitionsFromFile(f string) ([]D20Definitions, error) {
	var def []D20Definitions

	content, err := loadFile(f)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(content, &def)
	if err != nil {
		return def, err
	}
	return def, nil
}

type ExtensionD20 struct {
	Name          string           `json:"name"`
	Kind          string           `json:"kind"`   // class
	System        string           `json:"system"` // d20-3.5
	Description   string           `json:"description"`
	BaseAttack    string           `json:"base_attack"`          // "good", "average", "poor"
	Fortitude     string           `json:"fortitude"`            // "good", "poor"
	Reflex        string           `json:"reflex"`               // "good", "poor"
	Will          string           `json:"will"`                 // "good", "poor"
	Dice          string           `json:"dice"`                 // d10, d12
	FirstLevel    []string         `json:"first_level_features"` // proficiency_ simple|martial weapon, all armors, shield
	ClassFeatures map[int][]string `json:"class_fetures"`
	SpellPerDay   map[int][]int    `json:"spell_per_day"` // [class level] C L : [SL 0, 1, 2] [spell level 0, 1, 2...]
}

func LoadD20ExtensionFromFile(f string) ([]ExtensionD20, error) {
	var def []ExtensionD20

	content, err := loadFile(f)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(content, &def)
	if err != nil {
		return def, err
	}
	return def, nil
}

type ExtensionWeaponD20 struct {
	Name        string `json:"name"`
	AttackBonus int    `json:"attack_bonus"`
	Description string `json:"description"`
	DamageDice  string `json:"damage"`
}

func LoadWeaponD20FromFile(f string) ([]ExtensionWeaponD20, error) {
	var def []ExtensionWeaponD20

	content, err := loadFile(f)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(content, &def)
	if err != nil {
		return def, err
	}
	return def, nil
}

type ExtensionWeaponD10 struct {
	Name        string `json:"name"`
	Damage      int    `json:"damage"`
	Description string `json:"description"`
}

func LoadWeaponD10FromFile(f string) ([]ExtensionWeaponD10, error) {
	var def []ExtensionWeaponD10

	content, err := loadFile(f)
	if err != nil {
		return def, err
	}

	err = json.Unmarshal(content, &def)
	if err != nil {
		return def, err
	}
	return def, nil
}
