package definitions

import (
	"encoding/json"
	"io"
	"os"
)

const (
	AbilityKind string = "Ability"
	SkillKind   string = "Skill"
)

type Definitions struct {
	Name        string   `json:"name"`         // it should not change as RPG system mechanics like combat require some abilities names to use.
	DisplayName string   `json:"display_name"` // used in DnD as short
	Kind        string   `json:"kind"`         // Ability, Skill
	System      string   `json:"system"`       // d20-3.5, D10HomeMade, D10OldSchool
	Base        string   `json:"base"`         // used for skills to easy predict a skill check
	Group       string   `json:"group"`        // used in d10 to group abilities and skills
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

func LoadFromFile(f string) ([]Definitions, error) {
	var def []Definitions

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

func loadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}
