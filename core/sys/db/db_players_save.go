package db

import (
	"context"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rpg/d10hm"
	"github.com/betorvs/playbypost/core/rules"
)

func (db *DBX) SavePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpg *rpg.RPGSystem) (int, error) {

	var playerID int
	var err error
	{
		r, er := db.savePlayer(ctx, id, storyID, npc, creature)
		playerID = r
		err = er
	}
	{
		switch {
		case rpg.Name == "d20-3.5":
		case rpg.Name == "D10HomeMade":
			_, er := db.saveExtensionD10HM(ctx, playerID, npc, creature.Extension.(*d10hm.D10Extented))
			err = er

		case rpg.Name == "D10OldSchool":

		}

	}
	return playerID, err
}

func (db *DBX) savePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature) (int, error) {
	query := "INSERT INTO players(character_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	if npc {
		query = "INSERT INTO non_players(npc_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	}
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("save players prepare failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(creature.Name, id, storyID, false, creature.Abilities, creature.Skills, creature.RPG.Name).Scan(&res)
	if err != nil {
		db.logger.Error("save players queryRow failed", "error", err.Error())
		return -1, err
	}

	return res, nil
}

func (db *DBX) SavePlayerTx(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpgSystem *rpg.RPGSystem) (int, error) {
	var playerID int
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.logger.Error("tx begin on players failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	query := "INSERT INTO players(character_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	if npc {
		query = "INSERT INTO non_players(npc_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	}
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("tx prepare on players failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	err = tx.StmtContext(ctx, stmt).QueryRow(creature.Name, id, storyID, false, creature.Abilities, creature.Skills, creature.RPG.Name).Scan(&playerID)
	if err != nil {
		db.logger.Error("tx statement on players failed", "error", err.Error())
		return -1, err
	}

	// switch rpgSystem.Name {
	// case rpg.D2035:
	// case rpg.D10HM:
	// 	query := "INSERT INTO extension_d10_homemade_pc(player_id, health, defense, willpower, initiative, size, armor, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	// 	if npc {
	// 		query = "INSERT INTO extension_d10_homemade_npc(player_id, health, defense, willpower, initiative, size, armor, weapon) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	// 	}
	// 	extension := creature.Extension.(*d10hm.D10Extented)
	// 	_, err = tx.ExecContext(ctx, query, playerID, extension.Health, extension.Defense, extension.WillPower, extension.Initiative, extension.Size, extension.Armor, extension.Weapon)
	// 	if err != nil {
	// 		db.logger.Error("tx extension on players failed", "error", err.Error())
	// 		return -1, err
	// 	}

	// case rpg.D10OS:

	// }

	if err = tx.Commit(); err != nil {
		db.logger.Error("tx commit on initiative failed", "error", err.Error())
		return -1, err
	}
	return playerID, nil
}
