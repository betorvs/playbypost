package d10hm

import (
	"fmt"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
)

type StorytellingCharacter struct {
	base.Creature
	D10Extented
}

func New(n string, r *rpg.RPGSystem) *StorytellingCharacter {
	return &StorytellingCharacter{
		Creature:    *base.NewCreature(n, r),
		D10Extented: *newExtended(),
	}
}

func (c *StorytellingCharacter) Name() string {
	return c.Creature.Name
}

func (c *StorytellingCharacter) SetName(n string) error {
	if n == "" {
		return fmt.Errorf("name is empty")
	}
	c.Creature.Name = n
	return nil
}

func (c *StorytellingCharacter) RPGSystem() *rpg.RPGSystem {
	return c.Creature.RPG
}

func (c *StorytellingCharacter) Damage(v int) error {
	c.Health = c.Health - v
	return nil
}

func (c *StorytellingCharacter) HealthStatus() int {
	return c.Health
}

func (c *StorytellingCharacter) SetWeapon(name, kind string, value int, description string) {
	c.D10Extented.Weapon.SetWeapon(name, kind, value, description)
}

func (c *StorytellingCharacter) SetArmor(v int) {
	c.D10Extented.Armor = v
}

func (d D10Extented) weaponBonus(s string) (int, error) {
	w, ok := d.Weapon.GetWeapon(s)
	if !ok {
		return 0, fmt.Errorf("weapon %s not found", s)
	}
	return w.Value, nil
}

// type Weapon struct {
// 	Name        string `json:"name"`
// 	Value       int    `json:"value"`
// 	Description string `json:"description"`
// }

type D10Extented struct {
	Health     int `json:"health"`
	Defense    int `json:"defense"`
	WillPower  int `json:"willpower"`
	Initiative int `json:"initiative"`
	// Merits      []Advantages
	// Flaws       []Advantages
	// MoralSystem MoralSystem
	Size   int          `json:"size"`
	Armor  int          `json:"armor"`
	Weapon base.Weapons `json:"weapon"`
}

// func (a D10Extented) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// type Weapons map[string]Weapon

// func (a Weapons) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// func (a *Weapons) Scan(value interface{}) error {
// 	var data []byte
// 	if b, ok := value.([]byte); ok {
// 		data = b
// 	} else if s, ok := value.(string); ok {
// 		data = []byte(s)
// 	} else if value == nil {
// 		return nil
// 	} else if v, ok := value.(map[string]interface{}); ok {
// 		data, _ = json.Marshal(v)
// 	} else {
// 		return fmt.Errorf("unable to convert %v to []byte", value)
// 	}

// 	return json.Unmarshal(data, &a)
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
	_ = a.Weapon.Scan(values["weapon"])
}

// func (d D10Extented) String() string {
// 	return fmt.Sprintf("Creature Extended: Health %d, Defense %d, WillPower %d, Initiative %d, Size %d, Weapons %v", d.Health, d.Defense, d.WillPower, d.Initiative, d.Size, d.Weapon)
// }

//

func (a D10Extented) getValues() map[string]interface{} {
	weapon := "weapon"
	var weaponValue int
	for _, v := range a.Weapon {
		weapon = fmt.Sprintf("weapon:%s", v.Name)
		weaponValue = v.Value
	}

	return map[string]interface{}{
		"health":     a.Health,
		"defense":    a.Defense,
		"willpower":  a.WillPower,
		"initiative": a.Initiative,
		"size":       a.Size,
		"armor":      a.Armor,
		weapon:       weaponValue,
	}
}

func newWithValuesExtended(resolve, composture, dexterity, wits int) D10Extented {
	size := 5
	weapon := base.NewWeaponEmpty()
	return D10Extented{
		Health:     size + resolve,
		WillPower:  resolve + composture,
		Initiative: dexterity + composture,
		Defense:    lowerValue(dexterity, wits),
		Size:       size,
		Weapon:     weapon,
	}
}

func newExtended() *D10Extented {
	return &D10Extented{
		Weapon: base.NewWeaponEmpty(),
	}
}

func lowerValue(a, b int) int {
	if a < b {
		return a
	}
	return b
}
