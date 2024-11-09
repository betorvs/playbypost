package library

import (
	"fmt"
	"strings"
)

const (
	AbilityKind string = "Ability"
	SkillKind   string = "Skill"
)

type Library struct {
	Ability      AbilityDescription
	Skill        SkillDescription
	PFAncestry   PFAncestryDescription
	PFBackground PFBackgroundDescription
	PFClass      PFClassDescription
}

func New() *Library {
	return &Library{
		Ability: AbilityDescription{
			Grouped: make(map[string][]string),
			List:    []string{},
			Tags:    make(map[string][]string),
		},
		Skill: SkillDescription{
			List:      []string{},
			SkillBase: make(map[string]string),
			Grouped:   make(map[string][]string),
			Tags:      make(map[string][]string),
		},
		PFAncestry: PFAncestryDescription{
			Ancestries: make(map[string]PFAncestry),
			List:       []string{},
		},
		PFBackground: PFBackgroundDescription{
			Backgrounds: make(map[string]PFBackground),
			List:        []string{},
		},
		PFClass: PFClassDescription{
			Classes: make(map[string]PFClass),
			List:    []string{},
		},
	}
}

func (l *Library) String() string {
	common := fmt.Sprintf("Library imported. Abilities: %v, Skills: %v", len(l.Ability.List), len(l.Skill.List))
	if len(l.PFAncestry.List) > 0 {
		common = fmt.Sprintf("%s, Ancestries: %v", common, len(l.PFAncestry.List))
	}
	if len(l.PFBackground.List) > 0 {
		common = fmt.Sprintf("%s, Backgrounds: %v", common, len(l.PFBackground.List))
	}
	if len(l.PFClass.List) > 0 {
		common = fmt.Sprintf("%s, Classes: %v", common, len(l.PFClass.List))
	}
	return common
}

func (l *Library) AppendSkills(s string) {
	l.Skill.List = append(l.Skill.List, strings.ToLower(s))
}

func (l *Library) AppendAbilities(s string) {
	l.Ability.List = append(l.Ability.List, strings.ToLower(s))
}

func (l *Library) AppendSkillBase(skillName, abilityName string) {
	l.Skill.SkillBase[strings.ToLower(skillName)] = strings.ToLower(abilityName)
}

func (l *Library) AppendAbilityPerGroup(group, ability string) {
	l.Ability.Grouped[strings.ToLower(group)] = append(l.Ability.Grouped[strings.ToLower(group)], strings.ToLower(ability))
}

func (l *Library) AppendAbilityTags(ability string, tags []string) {
	for _, v := range tags {
		l.Ability.Tags[strings.ToLower(ability)] = append(l.Ability.Tags[strings.ToLower(ability)], strings.ToLower(v))
	}
}

func (l *Library) AppendSkillPerGroup(group, skill string) {
	l.Skill.Grouped[strings.ToLower(group)] = append(l.Skill.Grouped[strings.ToLower(group)], strings.ToLower(skill))
}

func (l *Library) AppendSkillTags(skill string, tags []string) {
	for _, v := range tags {
		l.Skill.Tags[strings.ToLower(skill)] = append(l.Skill.Tags[strings.ToLower(skill)], strings.ToLower(v))
	}
}

func (l *Library) GetSkillBase(skillName string) string {
	if s, ok := l.Skill.SkillBase[skillName]; ok {
		return s
	}
	return ""
}

type AbilityDescription struct {
	Grouped map[string][]string
	List    []string
	Tags    map[string][]string
}

type SkillDescription struct {
	List      []string
	SkillBase map[string]string
	Grouped   map[string][]string
	Tags      map[string][]string
}

type Descriptions struct {
	Name        string   `json:"name"`         // it should not change as RPG system mechanics like combat require some abilities names to use.
	DisplayName string   `json:"display_name"` // used in DnD as short
	Kind        string   `json:"kind"`         // Ability, Skill
	System      string   `json:"system"`       // d20-3.5, D10HomeMade, D10OldSchool
	Base        string   `json:"base"`         // used for skills to easy predict a skill check
	Group       string   `json:"group"`        // used in d10 to group abilities and skills
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}
