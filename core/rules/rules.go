package rules

import (
	"context"
	"log/slog"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/sys/library"
)

type RolePlaying interface {
	Attack(kind, weapon string) int
	DefenseBonus(kind string) int
	InitiativeBonus() (int, error)
	Destroy() error
	IsDead() bool
	Save(ctx context.Context, id, sID int, save func(ctx context.Context, id, sID int, creature *base.Creature, extension map[string]interface{}) (int, error)) (int, error)
	Update(ctx context.Context, id int, update func(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error) error
	SkillCheck(d rpg.RollInterface, check Check, logger *slog.Logger, lib *library.Library) (Result, error)
	AbilityCheck(d rpg.RollInterface, check Check, logger *slog.Logger, lib *library.Library) (Result, error)
	HealthStatus() int
	Damage(v int) error
	Name() string
	RPGSystem() *rpg.RPGSystem
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
