package base

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Weapon struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}

type Weapons map[string]Weapon

func NewWeaponEmpty() Weapons {
	return make(Weapons)
}

func NewWeapons(name, kind string, value int, description string) Weapons {
	weapons := make(Weapons)
	weapons.SetWeapon(name, kind, value, description)
	return weapons
}

func (a Weapons) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Weapons) Scan(value interface{}) error {
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

func (w Weapons) GetWeapon(name string) (Weapon, bool) {
	v, ok := w[name]
	return v, ok
}

func (w Weapons) SetWeapon(name, kind string, value int, description string) {
	w[name] = Weapon{name, kind, value, description}
}
