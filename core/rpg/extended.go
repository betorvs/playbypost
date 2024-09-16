package rpg

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/betorvs/playbypost/core/rpg/d10hm"
	"github.com/betorvs/playbypost/core/rpg/d20e35"
)

type ExtendedSystem interface {
	// AbilityBonus(s string) (int, error)
	SkillBonus(s string) (int, error)
	InitiativeBonus() (int, error)
	AttackBonus(s string) (int, error)
	DefenseBonus(s string) (int, error)
	WeaponBonus(s string) (int, string, error)
	Damage(v int) error
	HealthStatus() int
	SetWeapon(name string, value int, description string)
	SetArmor(v int)
	GetValues() map[string]interface{}
	String() string
}

func NewExtendedSystem(r *RPGSystem, values map[string]interface{}) ExtendedSystem {
	switch r.Name {
	case D10HM:
		extended := d10hm.NewExtended()
		extended.SetValues(values, convertInterfaceInt)
		return extended
	case D2035:
		return d20e35.NewExtended()
	default:
		return nil
	}
}

type Extended map[string]interface{}

func NewExtended() Extended {
	return make(map[string]interface{})
}

func (a Extended) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Extended) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed on extended")
	}

	return json.Unmarshal(b, &a)
}

func convertInterfaceInt(x interface{}) int {
	switch v := x.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		return 0
	default:
		return 0
	}
}
