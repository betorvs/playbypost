package library

import "strings"

// ancestries
// hit points size speed languages abilities boosts flaws traits
// example: human 8 ; medium ; 25 ; common ; 2 abilities
type PFAncestry struct {
	Name            string     `json:"name"`
	StartHitPoints  int        `json:"start_hit_points"`
	Size            string     `json:"size"`
	Speed           int        `json:"speed"`
	Languages       []string   `json:"languages"`
	LanguagesBonus  int        `json:"languages_bonus"` // number of languages you can learn + Intelligence modifier
	AbilitiesBoosts []PFBoosts `json:"abilities_boosts"`
	Flaws           []PFBoosts `json:"flaws"`
	Traits          []string   `json:"traits"`
}

type PFBoosts struct {
	Ability string `json:"ability"`
	OR      string `json:"or,omitempty"`
	AND     string `json:"and,omitempty"`
	Source  string `json:"source"`
}

type PFAncestryDescription struct {
	Ancestries map[string]PFAncestry `json:"ancestries"`
	List       []string              `json:"list"`
}

func (a *PFAncestryDescription) Append(ancestry PFAncestry) {
	a.Ancestries[strings.ToLower(ancestry.Name)] = ancestry
	a.List = append(a.List, strings.ToLower(ancestry.Name))
}

// backgrounds
// ability boosts, skill training, skill feat
// example: warrior ; 2 abilities ; 2 skills ; 1 feat
type PFBackground struct {
	Name            string     `json:"name"`
	AbilitiesBoosts []PFBoosts `json:"abilities_boosts"`
	SkillTraining   []string   `json:"skill_training"`
	SkillFeat       []string   `json:"skill_feat"`
}

type PFBackgroundDescription struct {
	Backgrounds map[string]PFBackground `json:"backgrounds"`
	List        []string                `json:"list"`
}

func (b *PFBackgroundDescription) Append(background PFBackground) {
	b.Backgrounds[strings.ToLower(background.Name)] = background
	b.List = append(b.List, strings.ToLower(background.Name))
}

// classes
// ability boosts, hit points, key ability, initial proficiencies, proficiencies, skills, feats, equipment, spells
// - skills
// - feats
// - spells
type PFClass struct {
	Name              string          `json:"name"`
	Level             int             `json:"level"`
	AbilitiesBoosts   []PFBoosts      `json:"abilities_boosts"`
	HitPointsPerLevel int             `json:"hit_points_per_level"`
	Proficiencies     []PFProficiency `json:"proficiencies"`
	FreeSkills        int             `json:"free_skills"` // number of free skills you can choose plus Intelligence modifier
	Features          []PFFeatures    `json:"features"`
	HasSpell          bool            `json:"has_spell"`
	SpellsPerDay      map[int]int     `json:"spells_per_day"` // level and number of spells
}

type PFFeatures struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type PFProficiency struct {
	On     string `json:"on"`
	Or     string `json:"or,omitempty"`
	Level  string `json:"level"`
	Source string `json:"source"`
}

type PFClassDescription struct {
	Classes map[string]PFClass `json:"classes"`
	List    []string           `json:"list"`
}

func (c *PFClassDescription) Append(class PFClass) {
	c.Classes[strings.ToLower(class.Name)] = class
	c.List = append(c.List, strings.ToLower(class.Name))
}
