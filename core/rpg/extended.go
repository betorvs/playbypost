package rpg

type ExtendedSystem interface {
	// AbilityBonus(s string) (int, error)
	SkillBonus(s string) (int, error)
	InitiativeBonus() (int, error)
	AttackBonus(s string) (int, error)
	DefenseBonus(s string) (int, error)
	WeaponBonus(s string) (int, string, error)
	Damage(v int) error
	HealthStatus() int
	String() string
}
