package d10os

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type D10Extented struct {
	Health     int
	WillPower  int
	Initiative int
	// Merits     []Advantages
	// Flaws      []Advantages
	Virtues Virtues
	Size    int
	Armor   int
	Weapon  Weapons
}

type Advantages struct {
	Name        string
	Value       int
	Description string
}

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
	Name        string
	Value       int
	Description string
}

type Virtues struct {
	ConscienceConviction int
	SelfControlInstinct  int
	Courage              int
}

func NewWithValuesExtended(conscience, selfControl, courage, dexterity, wits int) *D10Extented {
	size := 5
	weapon := make(map[string]Weapon)
	return &D10Extented{
		Health:     7,
		WillPower:  5 + courage,
		Initiative: dexterity + wits,
		Virtues: Virtues{
			ConscienceConviction: conscience,
			SelfControlInstinct:  selfControl,
			Courage:              courage,
		},
		Size:   size,
		Weapon: weapon,
	}
}

func NewExtended() *D10Extented {
	return &D10Extented{
		Weapon: make(map[string]Weapon),
	}
}

func (d *D10Extented) SkillBonus(s string) (int, error) {
	return 0, nil
}
func (d *D10Extented) InitiativeBonus() (int, error) {
	return d.Initiative, nil
}
func (d *D10Extented) AttackBonus(s string) (int, error) {
	return 0, nil
}
func (d *D10Extented) DefenseBonus(s string) (int, error) {
	return 0, nil
}
func (d *D10Extented) WeaponBonus(s string) (int, string, error) {
	return 0, "", nil
}
func (d *D10Extented) Damage(v int) error {
	d.Health = d.Health - v
	return nil
}
func (d *D10Extented) HealthStatus() int {
	return d.Health
}
func (d *D10Extented) String() string {
	return fmt.Sprintf("Creature Extended: Health %d, WillPower %d, Initiative %d, Armor %d, Conscience/Conviction %d, SelfControl/Instinct %d, Courage %d, Weapons %v", d.Health, d.WillPower, d.Initiative, d.Armor, d.Virtues.ConscienceConviction, d.Virtues.SelfControlInstinct, d.Virtues.Courage, d.Weapon)
}
