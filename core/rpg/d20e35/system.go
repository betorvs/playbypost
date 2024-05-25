package d20e35

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type D20Extended struct {
	Level        int
	HitPoints    int
	HitDice      HitDices
	BaseAttack   KindBonus
	SavingThrows SavingThrows
	ArmorClass   int
	Class        MultiClass
	Multiclass   bool
	Race         string
	Size         string
	// ClassAbilities []OtherAbilities
	// RaceAbilities  []OtherAbilities
	// Feat           []OtherAbilities
	Weapon Weapons
}

type MultiClass map[string]int

func (a MultiClass) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *MultiClass) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
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
	Name        string `json:"name"`
	AttackBonus int    `json:"attack_bonus"`
	Description string `json:"description"`
	DamageDice  string `json:"damage"`
}

type OtherAbilities struct {
	Name        string   `json:"name"`
	Bonus       int      `json:"bonus"`
	Description string   `json:"description"`
	ApplyOn     []string `json:"apply_on"`
	Limitation  []string `json:"limitation"`
}

func (a OtherAbilities) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *OtherAbilities) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (d D20Extended) String() string {
	return fmt.Sprintf("Creature Extented: level %d, Hit Points %d, Armor Class %d, Class: %v, Race: %s, Size: %s, Weapons: %v", d.Level, d.HitPoints, d.ArmorClass, d.Class, d.Race, d.Size, d.Weapon)
}

func (d *D20Extended) AbilityBonus(s string) (int, error) {
	return 0, nil
}
func (d *D20Extended) SkillBonus(s string) (int, error) {
	return 0, nil
}

func (d *D20Extended) InitiativeBonus() (int, error) {
	return 0, nil
}

func (d *D20Extended) AttackBonus(s string) (int, error) {
	value := d.baseAttackCalc(d.BaseAttack)
	return value, nil
}

// Defense can be:
// AC, Fortitude, Reflex, Will
func (d *D20Extended) DefenseBonus(s string) (int, error) {
	switch strings.ToLower(s) {
	case "AC", "ArmorCless":
		return 10, nil
	case Fortitude, "for":
		value := d.savingThrowsCalcBonus(Fortitude)
		return value, nil
	case Reflex, "ref":
		value := d.savingThrowsCalcBonus(Reflex)
		return value, nil
	case Will, "wil":
		value := d.savingThrowsCalcBonus(Will)
		return value, nil
	}
	return 0, nil
}

func (d *D20Extended) WeaponBonus(s string) (int, string, error) {
	if value, ok := d.Weapon[s]; ok {
		return value.AttackBonus, value.DamageDice, nil
	}
	return 0, "", nil
}

func (d *D20Extended) Damage(v int) error {
	d.HitPoints = d.HitPoints - v
	return nil
}

func (d *D20Extended) HealthStatus() int {
	return d.HitPoints
}

func (d *D20Extended) calcHitPoints(level, bonus int, method HitPointsMethod) int {
	value := d.HitDice.Value(method) + bonus
	if level > 1 {
		multiplier := (level - 1) * bonus
		value += multiplier
	}
	return value
}

func NewExtended(bonus, level int, hitDices HitDices, hitDicesMethod HitPointsMethod, class string) *D20Extended {
	weapon := make(map[string]Weapon)
	classes := make(map[string]int)
	classes[class] = level
	d := D20Extended{
		Level:   level,
		Class:   classes,
		HitDice: hitDices,
		Weapon:  weapon,
	}
	d.calcHitPoints(level, bonus, hitDicesMethod)

	return &d
}

func RestoreExtended() *D20Extended {
	return &D20Extended{}
}
