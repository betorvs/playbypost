package pg

import (
	"context"
	"fmt"

	"github.com/betorvs/playbypost/core/rpg/base"
	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) SavePlayerTx(ctx context.Context, id, storyID int, creature *base.Creature, extension map[string]interface{}) (int, error) {
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
	ext := types.NewExtension()
	ext.ConvertMap(extension)

	query := "INSERT INTO players(character_name, player_id, stage_id, destroyed, abilities, skills, extensions, rpg) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on players failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			db.Logger.Error("error closing stmt", "error", err)
		}
	}()
	err = tx.StmtContext(ctx, stmt).QueryRow(creature.Name, id, storyID, false, creature.Abilities, creature.Skills, ext, creature.RPG.Name).Scan(&playerID)
	if err != nil {
		db.Logger.Error("tx statement on players failed", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on players failed", "error", err.Error())
		return -1, err
	}
	return playerID, nil
}

func (db *DBX) UpdatePlayer(ctx context.Context, id int, creature *base.Creature, extension map[string]interface{}, destroyed bool) error {
	ext := types.NewExtension()
	ext.ConvertMap(extension)
	query := "UPDATE players SET abilities = $1, skills = $2, extensions = $3, destroyed = $4 WHERE id = $5" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("update players prepare failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			db.Logger.Error("error closing stmt", "error", err)
		}
	}()
	db.Logger.Debug("update player", "creature", creature)
	_, err = stmt.ExecContext(ctx, creature.Abilities, creature.Skills, ext, destroyed, id)
	if err != nil {
		db.Logger.Error("update players exec failed", "error", err.Error())
		return err
	}
	return nil
}

func (db *DBX) UpdatePlayerDetails(ctx context.Context, id int, name, rpg string) error {
	query := "UPDATE players SET character_name = $1, rpg = $2 WHERE id = $3" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("update player details prepare failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			db.Logger.Error("error closing stmt", "error", err)
		}
	}()
	_, err = stmt.ExecContext(ctx, name, rpg, id)
	if err != nil {
		db.Logger.Error("update player details exec failed", "error", err.Error())
		return err
	}
	return nil
}
