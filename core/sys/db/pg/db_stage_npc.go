package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetNPCByStageID(ctx context.Context, id int) ([]types.Players, error) {
	players := []types.Players{}
	query := "SELECT id, npc_name, stage_id, storyteller_id, destroyed, abilities, skills, rpg FROM non_players WHERE stage_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on non_players by stage_id failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		p := types.NewPlayer()
		c := rules.RestoreCreature()
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

func (db *DBX) GenerateNPC(ctx context.Context, name string, stageID, storytellerID int, creature *rules.Creature) (int, error) {
	npcID := 0
	query := "INSERT INTO non_players (npc_name, stage_id, storyteller_id, abilities, skills, extensions, rpg) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := db.Conn.QueryRowContext(ctx, query, name, stageID, storytellerID, creature.Abilities, creature.Skills, creature.Extension, creature.RPG.Name).Scan(&npcID)
	if err != nil {
		db.Logger.Error("insert on non_players failed", "error", err.Error())
		return npcID, err
	}
	return npcID, nil
}

func (db *DBX) UpdateNPC(ctx context.Context, id int, creature *rules.Creature, destroyed bool) error {
	query := "UPDATE non_players SET abilities = $1, skills = $2, extensions = $3, destroyed = $4 WHERE id = $5"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("update non_players prepare failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	db.Logger.Debug("update NPC", "creature", creature)
	_, err = stmt.ExecContext(ctx, creature.Abilities, creature.Skills, creature.Extension, destroyed, id)
	if err != nil {
		db.Logger.Error("update non_players exec failed", "error", err.Error())
		return err
	}
	return nil
}
