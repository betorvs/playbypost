package rpg

import (
	"fmt"
)

const (
	AutoPlay string = "autoplay" // solo and didatic adventures
	D10      string = "d10base"
	D20      string = "d20Base"
	D2035    string = "d20-3.5"     // Based on d20 3.5 SRD
	D10HM    string = "D10HomeMade" // D10 based on World of Darkness
	PFD20    string = "Pathfinder"  // Pathfinder d20
)

/*
Equal and LowerThan does not have implementation at the moment. LowerThan exist in OldDragon System
Sum returns sum of all values from all dices, ex 3d6.
GreaterThan means you roll a single dice, ex 1d20, and if your result is greather than, you succeed
CountResults means you roll several dice, for example 5d10, with the same result, for example greater than 7, you count how many good results you have.
DifficultAndCount means rolling several dice, say 5d10, with a target of 5, and counting how many good rolls you get
FromResult, used only for Damage Calculation, means it will use the same as result from attack
*/
type Measurement int

const (
	Equal Measurement = iota
	Sum
	GreaterThan
	LowerThan
	CountResults
	DifficultAndCount
	FromResult
)

type RankingSystem int

/*
OnePerOne is used in d10 system which means each point represents a new dice to roll
https://www.dandwiki.com/wiki/UA:Alternative_Skill_Systems
SkillRanks = default skill system from D20
SkillModifiers = 8 class skil and 4 not class skill + complex multiclass system
LevelBasedSkill = add level to check if class skill
ProficiencyRankAndLevel = add proficiency rank and level to check if class/background skill
*/

const (
	OnePerOne RankingSystem = iota
	AbilityD20
	SkillRanks
	SkillModifiers
	LevelBasedSkill
	ProficiencyRankAndLevel
)

type RPGSystem struct {
	Name              string
	BaseSystem        string
	BaseDice          string
	SuccessRule       Measurement
	AbilityRank       RankingSystem
	SkillRank         RankingSystem
	DamageCalculation Measurement
	// Ability           AbilityDescription
	// Skill             SkillDescription
	RestrictiveTasks bool
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

func (r *RPGSystem) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Name %s, BaseSystem %s, Base Dice %s", r.Name, r.BaseSystem, r.BaseDice)
}

func LoadRPGSystemsDefault(k string) *RPGSystem {
	// a := AbilityDescription{
	// 	Grouped: make(map[string][]string),
	// 	List:    []string{},
	// 	Tags:    make(map[string][]string),
	// }
	// s := SkillDescription{
	// 	List:      []string{},
	// 	SkillBase: make(map[string]string),
	// 	Grouped:   make(map[string][]string),
	// 	Tags:      make(map[string][]string),
	// }
	d20 := &RPGSystem{
		Name:              D2035,
		BaseSystem:        D20,
		BaseDice:          "1d20",
		SuccessRule:       GreaterThan,
		SkillRank:         SkillRanks,
		DamageCalculation: Sum,
		// Ability:           a,
		// Skill:             s,
	}
	switch k {
	case AutoPlay:
		return &RPGSystem{
			Name:             AutoPlay,
			BaseDice:         "1d6",
			RestrictiveTasks: true,
		}
	case D10HM:
		return &RPGSystem{
			Name:              D10HM,
			BaseDice:          "d10",
			BaseSystem:        D10,
			SuccessRule:       CountResults,
			AbilityRank:       OnePerOne,
			SkillRank:         OnePerOne,
			DamageCalculation: FromResult,
			// Ability:           a,
			// Skill:             s,
		}
	case PFD20:
		return &RPGSystem{
			Name:              PFD20,
			BaseSystem:        D20,
			BaseDice:          "1d20",
			SuccessRule:       GreaterThan,
			SkillRank:         ProficiencyRankAndLevel,
			DamageCalculation: Sum,
			// Ability:           a,
			// Skill:             s,
		}
	case D2035:
		return d20
	default:
		return d20
	}
}

// func (r *RPGSystem) AppendSkills(s string) {
// 	r.Skill.List = append(r.Skill.List, strings.ToLower(s))
// }

// func (r *RPGSystem) AppendAbilities(s string) {
// 	r.Ability.List = append(r.Ability.List, strings.ToLower(s))
// }

// func (r *RPGSystem) AppendSkillBase(skillName, abilityName string) {
// 	r.Skill.SkillBase[strings.ToLower(skillName)] = strings.ToLower(abilityName)
// }

// func (r *RPGSystem) AppendAbilityPerGroup(group, ability string) {
// 	r.Ability.Grouped[strings.ToLower(group)] = append(r.Ability.Grouped[strings.ToLower(group)], strings.ToLower(ability))
// }

// func (r *RPGSystem) AppendAbilityTags(ability string, tags []string) {
// 	for _, v := range tags {
// 		r.Ability.Tags[strings.ToLower(ability)] = append(r.Ability.Tags[strings.ToLower(ability)], strings.ToLower(v))
// 	}
// }

// func (r *RPGSystem) AppendSkillPerGroup(group, skill string) {
// 	r.Skill.Grouped[strings.ToLower(group)] = append(r.Skill.Grouped[strings.ToLower(group)], strings.ToLower(skill))
// }

// func (r *RPGSystem) AppendSkillTags(skill string, tags []string) {
// 	for _, v := range tags {
// 		r.Skill.Tags[strings.ToLower(skill)] = append(r.Skill.Tags[strings.ToLower(skill)], strings.ToLower(v))
// 	}
// }

// func (r *RPGSystem) GetSkillBase(skillName string) string {
// 	if s, ok := r.Skill.SkillBase[skillName]; ok {
// 		return s
// 	}
// 	return ""
// }

func (r *RPGSystem) InitiativeDice() string {
	switch r.Name {
	case D2035, PFD20:
		return r.BaseDice
	case D10HM:
		return fmt.Sprintf("%d%s", 1, r.BaseDice)
	}
	return r.BaseDice
}

// func (r *RPGSystem) InitDefinitions(f string, logger *slog.Logger) {
// 	content, err := definitions.LoadFromFile(f)
// 	if err != nil {
// 		logger.Error("error from definitions", "error", err.Error())
// 		os.Exit(2)
// 	}
// 	// fmt.Printf("content: %#v\n", content)
// 	for _, v := range content {
// 		// fmt.Printf("%s: %v \n", v.Name, v.Kind)
// 		if v.Kind == definitions.AbilityKind {
// 			r.AppendAbilities(v.Name)
// 			if v.Group != "" {
// 				r.AppendAbilityPerGroup(v.Group, v.Name)
// 			}
// 			if len(v.Tags) != 0 {
// 				r.AppendAbilityTags(v.Name, v.Tags)
// 			}
// 		}
// 		if v.Kind == definitions.SkillKind {
// 			r.AppendSkills(v.Name)
// 			if slices.Contains(r.Ability.List, v.Base) {
// 				r.AppendSkillBase(v.Name, v.Base)
// 			}
// 			if v.Group != "" {
// 				r.AppendSkillPerGroup(v.Group, v.Name)
// 			}
// 			if len(v.Tags) != 0 {
// 				r.AppendSkillTags(v.Name, v.Tags)
// 			}
// 		}

// 	}
// }
