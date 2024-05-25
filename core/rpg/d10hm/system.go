package d10hm

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type D10Extented struct {
	Health     int
	Defense    int
	WillPower  int
	Initiative int
	// Merits      []Advantages
	// Flaws       []Advantages
	// MoralSystem MoralSystem
	Size   int
	Armor  int
	Weapon Weapons
}

func (d *D10Extented) String() string {
	if d == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Creature Extended: Health %d, Defense %d, WillPower %d, Initiative %d, Size %d, Weapons %v", d.Health, d.Defense, d.WillPower, d.Initiative, d.Size, d.Weapon)
}

type Advantages struct {
	Name        string
	Value       int
	Description string
}

// type MoralSystem struct {
// 	Name  string
// 	Value int
// }

type Weapons map[string]Weapon

func (a Weapons) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Weapons) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type Weapon struct {
	Name        string `json:"name"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

func NewExtended(resolve, composture, dexterity, wits int) *D10Extented {
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

func RestoreExtended() *D10Extented {
	return &D10Extented{}
}

func lowerValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *D10Extented) SkillBonus(s string) (int, error) {
	return 0, nil
}

func (d *D10Extented) InitiativeBonus() (int, error) {
	return d.Initiative, nil
}

func (d *D10Extented) AttackBonus(s string) (int, error) {
	if value, ok := d.Weapon[s]; ok {
		return value.Value, nil
	}
	return 0, nil
}

func (d *D10Extented) DefenseBonus(a string) (int, error) {
	if a == "ranged" {
		return d.Armor, nil
	}
	return d.Armor + d.Defense, nil
}

func (d *D10Extented) WeaponBonus(s string) (int, string, error) {
	if value, ok := d.Weapon[s]; ok {
		return value.Value, "", nil
	}
	return 0, "", nil
}

func (d *D10Extented) Damage(v int) error {
	d.Health = d.Health - v
	return nil
}

func (d *D10Extented) HealthStatus() int {
	return d.Health
}
