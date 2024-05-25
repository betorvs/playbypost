package db

import (
	"context"
	"errors"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/d10hm"
	"github.com/betorvs/playbypost/core/rpg/d10os"
	"github.com/betorvs/playbypost/core/rpg/d20e35"
)

func (db *DBX) SaveExtension(ctx context.Context, playerId int, npc bool, rpg *rpg.RPGSystem, extension interface{}) (int, error) {
	switch {
	case rpg.Name == "d20-3.5":
		res, err := db.saveExtensionD20(ctx, playerId, npc, extension.(*d20e35.D20Extended))
		return res, err
	case rpg.Name == "D10HomeMade":
		res, err := db.saveExtensionD10HM(ctx, playerId, npc, extension.(*d10hm.D10Extented))
		return res, err
	case rpg.Name == "D10OldSchool":
		res, err := db.saveExtensionD10OS(ctx, playerId, npc, extension.(*d10os.D10Extented))
		return res, err
	}
	return -1, errors.New("not implemented")
}
