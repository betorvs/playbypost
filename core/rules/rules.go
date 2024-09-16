package rules

const (
	AbilityInvalid string = "ability not valid"
	SkillInvalid   string = "skill not valid"
)

type RolePlayingGame struct {
}

type Check struct {
	Ability   string
	Override  string
	Skill     string
	Target    int
	Difficult int
}

type Result struct {
	Success     bool
	Description string
	Result      int
	Rolled      string
}
