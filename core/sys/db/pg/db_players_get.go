package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetPlayers(ctx context.Context, rpgSystem *rpg.RPGSystem) ([]types.Players, error) {
	players := []types.Players{}
	query := "SELECT id, character_name, stage_id, player_id, destroyed, abilities, skills, extensions, rpg FROM players"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on players failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		// var p types.Players
		p := types.NewPlayer()
		c := rules.RestoreCreature()
		extended := rpg.NewExtended()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &extended, &p.RPG); err != nil {
			db.Logger.Error("scan error on players", "error", err.Error())
		}
		c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
		types.CreatureToPlayer(p, c)
		// for k, v := range c.Abilities {
		// 	db.Logger.Debug("abilities", "k", k, "v", v)
		// 	key := k
		// 	if v.DisplayName != "" {
		// 		key = v.DisplayName
		// 	}
		// 	p.Abilities[key] = v.Value
		// }
		players = append(players, *p)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on players", "error", err.Error())
	}
	return players, nil
}

func (db *DBX) GetPlayerByID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error) {
	players := types.Players{}
	query := "SELECT id, character_name, stage_id, player_id, destroyed, abilities, skills, extensions, rpg  FROM players WHERE id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on players by id failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		p := types.NewPlayer()
		c := rules.RestoreCreature()
		extended := rpg.NewExtended()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &extended, &p.RPG); err != nil {
			db.Logger.Error("scan error on players by id ", "error", err.Error())
		}
		c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
		types.CreatureToPlayer(p, c)
		if p.ID > 0 {
			players = *p
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on players by id ", "error", err.Error())
	}
	return players, nil
}

func (db *DBX) GetPlayerByPlayerID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) (types.Players, error) {
	players := types.Players{}
	query := "SELECT id, character_name, stage_id, player_id, destroyed, abilities, skills, extensions, rpg FROM players WHERE player_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on players by player_id failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		p := types.NewPlayer()
		c := rules.RestoreCreature()
		c.RPG = rpgSystem
		extended := rpg.NewExtended()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &extended, &p.RPG); err != nil {
			db.Logger.Error("scan error on players by player_id ", "error", err.Error())
		}
		if p.ID > 0 {
			c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
			types.CreatureToPlayer(p, c)
			players = *p
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on players by player_id ", "error", err.Error())
	}
	return players, nil
}

func (db *DBX) GetPlayerByStageID(ctx context.Context, id int, rpgSystem *rpg.RPGSystem) ([]types.Players, error) {
	players := []types.Players{}
	query := "SELECT id, character_name, stage_id, player_id, destroyed, abilities, skills, extensions, rpg FROM players WHERE stage_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on players by stage_id failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		p := types.NewPlayer()
		c := rules.RestoreCreature()
		extended := rpg.NewExtended()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &extended, &p.RPG); err != nil {
			db.Logger.Error("scan error on players by stage_id", "error", err.Error())
		}
		c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
		types.CreatureToPlayer(p, c)
		if p.ID > 0 {
			players = append(players, *p)
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on players by stage_id", "error", err.Error())
	}
	return players, nil
}

func (db *DBX) GetPlayerByUserID(ctx context.Context, id, channel string, rpgSystem *rpg.RPGSystem) (types.Players, error) {
	players := types.Players{}
	query := "SELECT p.id, p.character_name, p.stage_id, p.player_id, p.destroyed, p.abilities, p.skills, p.extensions, p.rpg, u.userid, c.channel FROM players AS p JOIN users AS u ON p.player_id = u.id JOIN stage_channel AS c ON c.stage_id = p.stage_id WHERE u.userid = $1"
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on players by userid failed", "error", err.Error())
		return players, err
	}
	defer rows.Close()
	for rows.Next() {
		p := types.NewPlayer()
		c := rules.RestoreCreature()
		var userID, dbchannel string
		extended := rpg.NewExtended()
		if err := rows.Scan(&p.ID, &p.Name, &p.StageID, &p.PlayerID, &p.Destroyed, &c.Abilities, &c.Skills, &extended, &p.RPG, &userID, &dbchannel); err != nil {
			db.Logger.Error("scan error on players by userid ", "error", err.Error())
		}
		if p.ID > 0 {
			if channel != "" && channel == dbchannel {
				c.Extension = rpg.NewExtendedSystem(rpgSystem, extended)
				types.CreatureToPlayer(p, c)
				players = *p
			}
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on players by userid ", "error", err.Error())
	}
	return players, nil
}
