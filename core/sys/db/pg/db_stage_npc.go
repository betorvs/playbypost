package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetNPCByStageID(ctx context.Context, id int) ([]types.Players, error) {
	players := []types.Players{}
	query := "SELECT id, npc_name, stage_id, storyteller_id, destroyed, abilities, skills, rpg FROM non_players WHERE stage_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on non_players by stage_id failed", "error", err.Error())
		return players, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		p := types.NewPlayer()
		c := base.RestoreCreature()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &p.RPG); err != nil {
			db.Logger.Error("scan error on non_players by stage_id", "error", err.Error())
		}
		types.CreatureToPlayer(p, c)
		// if p.ID > 0 {
		players = append(players, *p)
		// }
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on non_players by stage_id", "error", err.Error())
	}
	return players, nil
}

func (db *DBX) GenerateNPC(ctx context.Context, stageID, storytellerID int, creature *base.Creature, extension map[string]interface{}) (int, error) {
	npcID := 0
	ext := types.NewExtension()
	ext.ConvertMap(extension)
	query := "INSERT INTO non_players (npc_name, stage_id, storyteller_id, abilities, skills, extensions, rpg) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id" // dev:finder+query
	err := db.Conn.QueryRowContext(ctx, query, creature.Name, stageID, storytellerID, creature.Abilities, creature.Skills, ext, creature.RPG.Name).Scan(&npcID)
	if err != nil {
		db.Logger.Error("insert on non_players failed", "error", err.Error())
		return npcID, db.parsePostgresError(err)
	}
	return npcID, nil
}

func (db *DBX) UpdateNPC(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error {
	ext := types.NewExtension()
	ext.ConvertMap(extension)
	query := "UPDATE non_players SET abilities = $1, skills = $2, extensions = $3, destroyed = $4 WHERE id = $5" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("update non_players prepare failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			db.Logger.Error("error closing stmt", "error", err)
		}
	}()
	db.Logger.Debug("update NPC", "creature", creature)
	_, err = stmt.ExecContext(ctx, creature.Abilities, creature.Skills, ext, destroyed, id)
	if err != nil {
		db.Logger.Error("update non_players exec failed", "error", err.Error())
		return err
	}
	return nil
}
