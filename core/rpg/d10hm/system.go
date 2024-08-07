package d10hm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type D10Extented struct {
	Health     int `json:"health"`
	Defense    int `json:"defense"`
	WillPower  int `json:"willpower"`
	Initiative int `json:"initiative"`
	// Merits      []Advantages
	// Flaws       []Advantages
	// MoralSystem MoralSystem
	Size   int     `json:"size"`
	Armor  int     `json:"armor"`
	Weapon Weapons `json:"weapon"`
}

func (a D10Extented) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// func (a *D10Extented) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}

// 	return json.Unmarshal(b, &a)
// }

func (a *D10Extented) SetValues(values map[string]interface{}, convertInterfaceInt func(interface{}) int) {
	if values == nil {
		return
	}
	a.Health = convertInterfaceInt(values["health"])
	a.Defense = convertInterfaceInt(values["defense"])
	a.WillPower = convertInterfaceInt(values["willpower"])
	a.Initiative = convertInterfaceInt(values["initiative"])
	a.Size = convertInterfaceInt(values["size"])
	a.Armor = convertInterfaceInt(values["armor"])
	a.Weapon.Scan(values["weapon"])
}

func (d D10Extented) String() string {
	// if d == nil {
	// 	return "<nil>"
	// }
	return fmt.Sprintf("Creature Extended: Health %d, Defense %d, WillPower %d, Initiative %d, Size %d, Weapons %v", d.Health, d.Defense, d.WillPower, d.Initiative, d.Size, d.Weapon)
}

// type Advantages struct {
// 	Name        string
// 	Value       int
// 	Description string
// }

// type MoralSystem struct {
// 	Name  string
// 	Value int
// }

type Weapons map[string]Weapon

func (a Weapons) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Weapons) Scan(value interface{}) error {
	// b, ok := value.([]byte)
	// if !ok {
	// 	return errors.New("type assertion to []byte failed")
	// }

	// return json.Unmarshal(b, &a)
	var data []byte
	if b, ok := value.([]byte); ok {
		data = b
	} else if s, ok := value.(string); ok {
		data = []byte(s)
	} else if value == nil {
		return nil
	} else if v, ok := value.(map[string]interface{}); ok {
		data, _ = json.Marshal(v)
	} else {
		return fmt.Errorf("unable to convert %v to []byte", value)
	}

	return json.Unmarshal(data, &a)
}

type Weapon struct {
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

func NewWithValuesExtended(resolve, composture, dexterity, wits int) *D10Extented {
	size := 5
	weapon := make(map[string]Weapon)
	return &D10Extented{
		Health:     size + resolve,
		WillPower:  resolve + composture,
		Initiative: dexterity + composture,
		Defense:    lowerValue(dexterity, wits),
		Size:       size,
		Weapon:     weapon,
	}
}

func NewExtended() *D10Extented {
	return &D10Extented{
		Weapon: make(map[string]Weapon),
	}
}

func lowerValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d D10Extented) SkillBonus(s string) (int, error) {
	return 0, nil
}

func (d D10Extented) InitiativeBonus() (int, error) {
	return d.Initiative, nil
}

func (d D10Extented) AttackBonus(s string) (int, error) {
	if value, ok := d.Weapon[s]; ok {
		return value.Value, nil
	}
	return 0, nil
}

func (d D10Extented) DefenseBonus(a string) (int, error) {
	if a == "ranged" {
		return d.Armor, nil
	}
	return d.Armor + d.Defense, nil
}

func (d D10Extented) WeaponBonus(s string) (int, string, error) {
	if value, ok := d.Weapon[s]; ok {
		return value.Value, "", nil
	}
	return 0, "", nil
}

func (d *D10Extented) Damage(v int) error {
	d.Health = d.Health - v
	return nil
}

func (d D10Extented) HealthStatus() int {
	return d.Health
}
