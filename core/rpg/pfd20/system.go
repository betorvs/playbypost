package pfd20

// System is the Pathfinder 2.0 system
// https://pf2.d20pfsrd.com/

type PFExtended struct {
	Ancestry    Ancestry
	Background  Background
	Class       Class
	HitPoints   int
	Proficiency map[string]Proficiency
}

func NewPFExtended() *PFExtended {
	return &PFExtended{}
}

func (p PFExtended) SkillBonus(s string) (int, error) {
	return 0, nil
}

func (p PFExtended) InitiativeBonus() (int, error) {
	return 0, nil
}

func (p PFExtended) AttackBonus(s string) (int, error) {
	return 0, nil
}

func (p PFExtended) DefenseBonus(s string) (int, error) {
	switch s {
	case Fortitude:
		p := ProficiencyRank(p.Proficiency["fortitude"].Level, p.Class.Level)
		return p, nil
	case Reflex:
		p := ProficiencyRank(p.Proficiency["reflex"].Level, p.Class.Level)
		return p, nil
	case Will:
		p := ProficiencyRank(p.Proficiency["will"].Level, p.Class.Level)
		return p, nil
	case ArmorClass:
		p := ProficiencyRank(p.Proficiency["armor_class"].Level, p.Class.Level)
		return p, nil
	}
	return 0, nil
}

func (p PFExtended) WeaponBonus(s string) (int, string, error) {
	return 0, "", nil
}

func (p PFExtended) Damage(v int) error {
	return nil
}

func (p PFExtended) HealthStatus() int {
	return 0
}

func (p PFExtended) SetWeapon(name string, value int, description string) {
}

func (p PFExtended) SetArmor(v int) {
}

func (p PFExtended) GetValues() map[string]interface{} {
	return nil
}

func (p PFExtended) String() string {
	return ""
}

// ancestries
// hit points size speed languages abilities boosts flaws traits
// example: human 8 ; medium ; 25 ; common ; 2 abilities
type Ancestry struct {
	Name            string   `json:"name"`
	StartHitPoints  int      `json:"start_hit_points"`
	Size            string   `json:"size"`
	Speed           int      `json:"speed"`
	Languages       []string `json:"languages"`
	LanguagesBonus  int      `json:"languages_bonus"` // number of languages you can learn + Intelligence modifier
	AbilitiesBoosts []Boosts `json:"abilities_boosts"`
	Flaws           []Boosts `json:"flaws"`
	Traits          []string `json:"traits"`
}

type Boosts struct {
	Ability string `json:"ability"`
	OR      string `json:"or,omitempty"`
	AND     string `json:"and,omitempty"`
	Source  string `json:"source"`
}

// backgrounds
// ability boosts, skill training, skill feat
// example: warrior ; 2 abilities ; 2 skills ; 1 feat
type Background struct {
	Name            string   `json:"name"`
	AbilitiesBoosts []Boosts `json:"abilities_boosts"`
	SkillTraining   []string `json:"skill_training"`
	SkillFeat       []string `json:"skill_feat"`
}

// classes
// ability boosts, hit points, key ability, initial proficiencies, proficiencies, skills, feats, equipment, spells
// - skills
// - feats
// - spells
type Class struct {
	Name              string        `json:"name"`
	Level             int           `json:"level"`
	AbilitiesBoosts   []Boosts      `json:"abilities_boosts"`
	HitPointsPerLevel int           `json:"hit_points_per_level"`
	Proficiencies     []Proficiency `json:"proficiencies"`
	FreeSkills        int           `json:"free_skills"` // number of free skills you can choose plus Intelligence modifier
	Features          []Features    `json:"features"`
	HasSpell          bool          `json:"has_spell"`
	SpellsPerDay      map[int]int   `json:"spells_per_day"` // level and number of spells
}

type Features struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Level       int    `json:"level"`
	Description string `json:"description"`
}

type Proficiency struct {
	On     string `json:"on"`
	Or     string `json:"or,omitempty"`
	Level  string `json:"level"`
	Source string `json:"source"`
}

// equipment

const (
	ProficiencyUntrained = "untrained"
	ProficiencyTrained   = "trained"
	ProficiencyExpert    = "expert"
	ProficiencyMaster    = "master"
	ProficiencyLegendary = "legendary"
)

// Proficiency is the level of proficiency + your level
func ProficiencyRank(rank string, level int) int {
	switch rank {
	case ProficiencyUntrained:
		return 0
	case ProficiencyTrained:
		return level + 2
	case ProficiencyExpert:
		return level + 4
	case ProficiencyMaster:
		return level + 6
	case ProficiencyLegendary:
		return level + 8
	}
	return 0
}
