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
  title VARCHAR(50) UNIQUE NOT NULL,
  encounters_id int NOT NULL REFERENCES encounters(id),
  subject VARCHAR(50) NOT NULL, //could be notes or announcement
  kind int NOT NULL DEFAULT 0,  //kind is related about what can be that task
  target int NOT NULL DEFAULT 0, // rpg system target in dices
  options JSONB, // map [display_text(encounters)] id next encounter
  finished BOOLEAN NOT NULL DEFAULT FALSE // if this is finished or not
*/

type Task struct {
	EncounterID int            `json:"encounter_id"`
	Title       string         `json:"title"`
	DisplayText string         `json:"display_text"`
	Kind        TaskKind       `json:"kind"`
	Checks      string         `json:"checks"`
	Target      int            `json:"target"`
	Options     map[string]int `json:"options"`
	Finished    bool           `json:"finished"`
}
