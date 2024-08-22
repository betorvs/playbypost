package pg

import (
	"context"
	"fmt"

	"github.com/betorvs/playbypost/core/rpg"
	"github.com/betorvs/playbypost/core/rules"
)

// func (db *DBX) SavePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature, rpg *rpg.RPGSystem) (int, error) {

// 	var playerID int
// 	var err error
// 	{
// 		r, er := db.savePlayer(ctx, id, storyID, npc, creature)
// 		playerID = r
// 		err = er
// 	}
// 	{
// 		switch {
// 		case rpg.Name == "d20-3.5":
// 		case rpg.Name == "D10HomeMade":
// 			_, er := db.saveExtensionD10HM(ctx, playerID, npc, creature.Extension.(*d10hm.D10Extented))
// 			err = er

// 		case rpg.Name == "D10OldSchool":

// 		}

// 	}
// 	return playerID, err
// }

// func (db *DBX) savePlayer(ctx context.Context, id, storyID int, npc bool, creature *rules.Creature) (int, error) {
// 	query := "INSERT INTO players(character_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
// 	if npc {
// 		query = "INSERT INTO non_players(npc_name, player_id, stage_id, destroyed, abilities, skills, rpg) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
// 	}
// 	stmt, err := db.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		db.Logger.Error("save players prepare failed", "error", err.Error())
// 		return -1, err
// 	}
// 	defer stmt.Close()
// 	var res int
// 	err = stmt.QueryRow(creature.Name, id, storyID, false, creature.Abilities, creature.Skills, creature.RPG.Name).Scan(&res)
// 	if err != nil {
// 		db.Logger.Error("save players queryRow failed", "error", err.Error())
// 		return -1, err
// 	}

// 	return res, nil
// }

func (db *DBX) SavePlayerTx(ctx context.Context, id, storyID int, creature *rules.Creature, rpgSystem *rpg.RPGSystem) (int, error) {
	var playerID int
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on players failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()

	query := "INSERT INTO players(character_name, player_id, stage_id, destroyed, abilities, skills, extensions, rpg) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on players failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	err = tx.StmtContext(ctx, stmt).QueryRow(creature.Name, id, storyID, false, creature.Abilities, creature.Skills, creature.Extension, rpgSystem.Name).Scan(&playerID)
	if err != nil {
		db.Logger.Error("tx statement on players failed", "error", err.Error())
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on players failed", "error", err.Error())
		return -1, err
	}
	return playerID, nil
}

func (db *DBX) UpdatePlayer(ctx context.Context, id int, creature *rules.Creature, destroyed bool) error {
	query := "UPDATE players SET abilities = $1, skills = $2, extensions = $3, destroyed = $4 WHERE id = $5"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("update players prepare failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	db.Logger.Info("update player", "creature", creature)
	_, err = stmt.ExecContext(ctx, creature.Abilities, creature.Skills, creature.Extension, destroyed, id)
	if err != nil {
		db.Logger.Error("update players exec failed", "error", err.Error())
		return err
	}
	return nil
}
