package pfd20

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg/base"
)

func (c *PathfinderCharacter) Save(ctx context.Context, id, sID int, save func(ctx context.Context, id, sID int, creature *base.Creature, extension map[string]interface{}) (int, error)) (int, error) {
	return save(ctx, id, sID, &c.Creature, c.getValues())
	// return 0, nil
}

func (c *PathfinderCharacter) Update(ctx context.Context, id int, update func(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error) error {
	err := update(ctx, id, &c.Creature, c.getValues(), c.IsDead())
	return err
}
