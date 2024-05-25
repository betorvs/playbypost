package db

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg/d10os"
)

func (db *DBX) saveExtensionD10OS(ctx context.Context, playerId int, npc bool, extension *d10os.D10Extented) (int, error) {
	query := "INSERT INTO extension_d10_oldschool_pc(player_id, health, willpower, initiative, size, armor, conscience_conviction, self_control_instinct, courage, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	if npc {
		query = "INSERT INTO extension_d10_oldschool_npc(player_id, health, willpower, initiative, size, armor, conscience_conviction, self_control_instinct, courage, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	}
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("save extenstion d10 oldschool prepare failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(playerId, extension.Health, extension.WillPower, extension.Initiative, extension.Size, extension.Armor, extension.Virtues.ConscienceConviction, extension.Virtues.SelfControlInstinct, extension.Virtues.Courage, extension.Weapon).Scan(&res)
	if err != nil {
		db.logger.Error("save extenstion d10 oldschool queryRow failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}
