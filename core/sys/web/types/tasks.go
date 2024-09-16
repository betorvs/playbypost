package types

type TaskKind int

const (
	EmptyTask        TaskKind = iota // empty
	CombatTask                       // combat
	SkillCheck                       // check skill against target
	AbilityCheck                     // check ability against target
	SocialTask                       // no actions required, wait master to close it
	ChooseTask                       // announce subject and ask for choice, use options as description and id to next encounter
	RandomChoiceTask                 // announce subject and ask to continue, roll a dice and open new task (use options like above)
)

func TaskKindAtoi(i int) TaskKind {
	return TaskKind(i)
}

func (t TaskKind) String() string {
	switch t {
	case CombatTask:
		return "CombatTask"
	case SkillCheck:
		return "SkillCheck"
	case AbilityCheck:
		return "AbilityCheck"
	case SocialTask:
		return "SocialTask"
	case ChooseTask:
		return "ChooseTask"
	case RandomChoiceTask:
		return "RandomChoiceTask"
	case EmptyTask:
		return "EmptyTask"
	}
	return "EmptyTask"
}

/*
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  description VARCHAR(50) UNIQUE NOT NULL,
  kind int NOT NULL DEFAULT 0,
  ability VARCHAR(50) NOT NULL,
  skill VARCHAR(50) NOT NULL,
  target int NOT NULL DEFAULT 0
*/

type Task struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Kind        TaskKind `json:"kind"`
	Ability     string   `json:"ability"`
	Skill       string   `json:"skill"`
	Target      int      `json:"target"`
}
