package db

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg/d20e35"
)

func (db *DBX) saveExtensionD20(ctx context.Context, playerId int, npc bool, extension *d20e35.D20Extended) (int, error) {
	query := "INSERT INTO extension_d20_e35_pc(player_id, level_total, hit_points, armor_class, class, multiclass, race, size, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	if npc {
		query = "INSERT INTO extension_d20_e35_npc(player_id, level_total, hit_points, armor_class, class, multiclass, race, size, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	}
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("save extenstion d20 prepare failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(playerId, extension.Level, extension.HitPoints, extension.ArmorClass, extension.Class, extension.Multiclass, extension.Race, extension.Size, extension.Weapon).Scan(&res)
	if err != nil {
		db.logger.Error("save extenstion d20 queryRow failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}
