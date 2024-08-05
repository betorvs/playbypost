package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg/d10hm"
)

func (db *DBX) saveExtensionD10HM(ctx context.Context, playerId int, npc bool, extension *d10hm.D10Extented) (int, error) {
	query := "INSERT INTO extension_d10_homemade_pc(player_id, health, defense, willpower, initiative, size, armor, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	if npc {
		query = "INSERT INTO extension_d10_homemade_npc(player_id, health, defense, willpower, initiative, size, armor, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	}
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("save extenstion d10 homemade prepare failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(playerId, extension.Health, extension.Defense, extension.WillPower, extension.Initiative, extension.Size, extension.Armor, extension.Weapon).Scan(&res)
	if err != nil {
		db.Logger.Error("save extenstion d10 homemade queryRow failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}
