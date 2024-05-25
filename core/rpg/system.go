package rpg

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/betorvs/dice"
	"github.com/betorvs/playbypost/core/sys/definitions"
)

const (
	Solo     string = "solo"     // solo adventures
	Didactic string = "didactic" // didactic adventures using 1d6
	D10      string = "d10base"
	D20      string = "d20Base"
	D2035    string = "d20-3.5"      // Based on d20 3.5 SRD
	D10HM    string = "D10HomeMade"  // D10 based on World of Darkness
	D10OS    string = "D10OldSchool" // D10 Based on old Vampire the masquerade 2nd edition system

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
*/

const (
	OnePerOne RankingSystem = iota
	AbilityD20
	SkillRanks
	SkillModifiers
	LevelBasedSkill
)

type RPGSystem struct {
	Name              string
	BaseSystem        string
	BaseDice          string
	SuccessRule       Measurement
	AbilityRank       RankingSystem
	SkillRank         RankingSystem
	DamageCalculation Measurement
	Ability           Ability
	Skill             Skill
	RestrictiveTasks  bool
	Logger            *slog.Logger
}

type Ability struct {
	Grouped map[string][]string
	List    []string
	Tags    map[string][]string
}

type Skill struct {
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

func LoadRPGSystemsDefault(k string, l *slog.Logger) *RPGSystem {
	a := Ability{
		Grouped: make(map[string][]string),
		List:    []string{},
		Tags:    make(map[string][]string),
	}
	s := Skill{
		List:      []string{},
		SkillBase: make(map[string]string),
		Grouped:   make(map[string][]string),
		Tags:      make(map[string][]string),
	}
	d20 := &RPGSystem{
		Name:              D2035,
		BaseSystem:        D20,
		BaseDice:          "1d20",
		SuccessRule:       GreaterThan,
		SkillRank:         SkillRanks,
		DamageCalculation: Sum,
		Ability:           a,
		Skill:             s,
		Logger:            l,
	}
	switch k {
	case Solo:
		return &RPGSystem{
			Name:             Solo,
			BaseDice:         "1d6",
			RestrictiveTasks: true,
			Logger:           l,
		}
	case Didactic:
		return &RPGSystem{
			Name:             Didactic,
			BaseDice:         "1d6",
			RestrictiveTasks: true,
			Logger:           l,
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
			Ability:           a,
			Skill:             s,
			Logger:            l,
		}
	case D10OS:
		return &RPGSystem{
			Name:              D10OS,
			BaseDice:          "d10",
			BaseSystem:        D10,
			SuccessRule:       DifficultAndCount,
			AbilityRank:       OnePerOne,
			SkillRank:         OnePerOne,
			DamageCalculation: DifficultAndCount,
			Ability:           a,
			Skill:             s,
			Logger:            l,
		}
	case D2035:
		return d20
	default:
		return d20
	}
}

func (r *RPGSystem) AppendSkills(s string) {
	r.Skill.List = append(r.Skill.List, strings.ToLower(s))
}

func (r *RPGSystem) AppendAbilities(s string) {
	r.Ability.List = append(r.Ability.List, strings.ToLower(s))
}

func (r *RPGSystem) AppendSkillBase(skillName, abilityName string) {
	r.Skill.SkillBase[strings.ToLower(skillName)] = strings.ToLower(abilityName)
}

func (r *RPGSystem) AppendAbilityPerGroup(group, ability string) {
	r.Ability.Grouped[strings.ToLower(group)] = append(r.Ability.Grouped[strings.ToLower(group)], strings.ToLower(ability))
}

func (r *RPGSystem) AppendAbilityTags(ability string, tags []string) {
	for _, v := range tags {
		r.Ability.Tags[strings.ToLower(ability)] = append(r.Ability.Tags[strings.ToLower(ability)], strings.ToLower(v))
	}
}

func (r *RPGSystem) AppendSkillPerGroup(group, skill string) {
	r.Skill.Grouped[strings.ToLower(group)] = append(r.Skill.Grouped[strings.ToLower(group)], strings.ToLower(skill))
}

func (r *RPGSystem) AppendSkillTags(skill string, tags []string) {
	for _, v := range tags {
		r.Skill.Tags[strings.ToLower(skill)] = append(r.Skill.Tags[strings.ToLower(skill)], strings.ToLower(v))
	}
}

func (r *RPGSystem) GetSkillBase(skillName string) string {
	if s, ok := r.Skill.SkillBase[skillName]; ok {
		return s
	}
	return ""
}

func (r *RPGSystem) InitiativeDice() string {
	switch r.Name {
	case D2035:
		return r.BaseDice
	case D10HM:
		return fmt.Sprintf("%d%s", 1, r.BaseDice)
	case D10OS:
		return fmt.Sprintf("%d%s", 1, r.BaseDice)
	}
	return r.BaseDice
}

func (r *RPGSystem) InitDefinitions(f string, logger *slog.Logger) {
	content, err := definitions.LoadFromFile(f)
	if err != nil {
		logger.Error("error from definitions", "error", err.Error())
		os.Exit(2)
	}
	// fmt.Printf("content: %#v\n", content)
	for _, v := range content {
		// fmt.Printf("%s: %v \n", v.Name, v.Kind)
		if v.Kind == definitions.AbilityKind {
			r.AppendAbilities(v.Name)
			if v.Group != "" {
				r.AppendAbilityPerGroup(v.Group, v.Name)
			}
			if len(v.Tags) != 0 {
				r.AppendAbilityTags(v.Name, v.Tags)
			}
		}
		if v.Kind == definitions.SkillKind {
			r.AppendSkills(v.Name)
			if slices.Contains(r.Ability.List, v.Base) {
				r.AppendSkillBase(v.Name, v.Base)
			}
			if v.Group != "" {
				r.AppendSkillPerGroup(v.Group, v.Name)
			}
			if len(v.Tags) != 0 {
				r.AppendSkillTags(v.Name, v.Tags)
			}
		}

	}
}

// Roll struct
type Roll struct {
	RPGSystem *RPGSystem
}

// FreeRoll func returns a int value, a string text and error
func (r Roll) FreeRoll(name, text string) (int, string, error) {
	diceRolled, _, err := dice.Roll(text)
	if err != nil {
		return 0, "No dices to roll", err
	} else {
		message := fmt.Sprintf("%s rolled %s and result %v with rolls %s", name, diceRolled.Description(), diceRolled.Int(), diceRolled.String())
		return diceRolled.Int(), message, nil
	}
}

func (r Roll) Check(name string) (int, string, error) {
	diceRolled, _, err := dice.Roll(r.RPGSystem.BaseDice)
	if err != nil {
		return 0, "No dices to roll", err
	} else {
		message := fmt.Sprintf("%s rolled %s and result %v with rolls %s", name, diceRolled.Description(), diceRolled.Int(), diceRolled.String())
		return diceRolled.Int(), message, nil
	}
}

func (r Roll) Dice(m, target int) string {
	dices := 1
	if m > 0 {
		dices = m
	}
	var dice string
	switch r.RPGSystem.Name {
	case D10HM:
		dice = fmt.Sprintf("%d%srv8", dices, r.RPGSystem.BaseDice)
	case D10OS:
		dice = fmt.Sprintf("%d%srv%d", dices, r.RPGSystem.BaseDice, target)
	}
	return dice
}
