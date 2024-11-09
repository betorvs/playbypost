package library

import (
	"log/slog"
	"os"
	"slices"
)

func (l *Library) initDescriptions(f string, logger *slog.Logger) {
	content, err := loadDescriptionsFromFile(f)
	if err != nil {
		logger.Error("error from definitions", "error", err.Error())
		os.Exit(2)
	}
	// fmt.Printf("content: %#v\n", content)
	for _, v := range content {
		// fmt.Printf("%s: %v \n", v.Name, v.Kind)
		if v.Kind == AbilityKind {
			l.AppendAbilities(v.Name)
			if v.Group != "" {
				l.AppendAbilityPerGroup(v.Group, v.Name)
			}
			if len(v.Tags) != 0 {
				l.AppendAbilityTags(v.Name, v.Tags)
			}
		}
		if v.Kind == SkillKind {
			l.AppendSkills(v.Name)
			if slices.Contains(l.Ability.List, v.Base) {
				l.AppendSkillBase(v.Name, v.Base)
			}
			if v.Group != "" {
				l.AppendSkillPerGroup(v.Group, v.Name)
			}
			if len(v.Tags) != 0 {
				l.AppendSkillTags(v.Name, v.Tags)
			}
		}

	}
}
