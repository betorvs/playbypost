package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetEncounters(ctx context.Context) ([]types.Encounter, error) {
	list := []types.Encounter{}
	query := "SELECT id, title, notes, announcement, story_id, writer_id, first_encounter, last_encounter FROM encounters" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on encounters failed", "error", err.Error())
		return list, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		var s types.Encounter
		if err := rows.Scan(&s.ID, &s.Title, &s.Notes, &s.Announcement, &s.StoryID, &s.WriterID, &s.FirstEncounter, &s.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on encounters", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error) {
	encounters := []types.Encounter{}
	query := "SELECT id, title, notes, announcement, story_id, writer_id, first_encounter, last_encounter FROM encounters WHERE story_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, storyID)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return encounters, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		var enc types.Encounter
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Notes, &enc.Announcement, &enc.StoryID, &enc.WriterID, &enc.FirstEncounter, &enc.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
		encounters = append(encounters, enc)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return encounters, nil
}

func (db *DBX) GetEncounterByStoryIDWithPagination(ctx context.Context, storyID, limit, cursor int) ([]types.Encounter, int, int, error) {
	encounters := []types.Encounter{}
	total := 0
	lastID := -1
	{
		query := "SELECT COUNT(*) FROM encounters WHERE story_id = $1" // dev:finder+query
		if err := db.Conn.QueryRowContext(ctx, query, storyID).Scan(&total); err != nil {
			db.Logger.Error("query on encounters by id failed", "error", err.Error())
			return encounters, lastID, total, err
		}

	}

	query := "SELECT id, title, notes, announcement, story_id, writer_id, first_encounter, last_encounter FROM encounters WHERE story_id = $1 AND id > $2 LIMIT $3" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, storyID, cursor, limit)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return encounters, lastID, total, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	count := 1
	for rows.Next() {
		var enc types.Encounter
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Notes, &enc.Announcement, &enc.StoryID, &enc.WriterID, &enc.FirstEncounter, &enc.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
		if count == limit {
			lastID = enc.ID
			db.Logger.Debug("limit reached", "limit", limit, "count", count, "encounter_id", enc.ID)
		}
		encounters = append(encounters, enc)
		count++
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return encounters, lastID, total, nil
}

func (db *DBX) GetEncounterByID(ctx context.Context, id int) (types.Encounter, error) {
	enc := types.Encounter{}
	query := "SELECT id, title, announcement, notes, story_id, writer_id, first_encounter, last_encounter FROM encounters WHERE id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return enc, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Announcement, &enc.Notes, &enc.StoryID, &enc.WriterID, &enc.FirstEncounter, &enc.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return enc, nil
}

func (db *DBX) CreateEncounterTx(ctx context.Context, title, announcement, notes string, storyID, storytellerID int, first, last bool) (int, error) {

	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// first can be only one
	// last can have multiple
	if first {
		var count int
		queryCheck := "SELECT COUNT(*) FROM encounters WHERE story_id = $1 AND first_encounter = true" // dev:finder+query
		if err = tx.QueryRowContext(ctx, queryCheck, storyID).Scan(&count); err != nil {
			if err != sql.ErrNoRows {
				db.Logger.Error("no rows passed", "err", err.Error())
				return -1, err
			}
		}
		if count > 0 {
			db.Logger.Error("first encounter already exists")
			return -1, fmt.Errorf("first encounter already exists")
		}
	}

	// Prepare the statement.
	query := "INSERT INTO encounters(title, announcement, notes, story_id, writer_id, first_encounter, last_encounter) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id" // dev:finder+query
	stmtInsert, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtInsert.Close()
		if err != nil {
			db.Logger.Error("error closing stmtInsert", "error", err)
		}
	}()
	var encounterID int
	err = tx.StmtContext(ctx, stmtInsert).QueryRow(title, announcement, notes, storyID, storytellerID, first, last).Scan(&encounterID)
	if err != nil {
		db.Logger.Error("error on insert into encounters", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit encounters", "error", err.Error())
		return -1, err
	}
	return encounterID, nil
}

func (db *DBX) UpdateEncounterTx(ctx context.Context, title, announcement, notes string, id, storyID int, first, last bool) (int, error) {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// first can be only one
	// last can have multiple
	if first {
		var firstID int
		queryCheck := "SELECT id FROM encounters WHERE story_id = $1 AND first_encounter = true" // dev:finder+query
		if err = tx.QueryRowContext(ctx, queryCheck, storyID).Scan(&firstID); err != nil {
			if err != sql.ErrNoRows {
				db.Logger.Error("no rows passed", "err", err.Error())
				return -1, err
			}
		}
		if firstID != id {
			db.Logger.Error("first encounter already exists")
			return -1, fmt.Errorf("first encounter already exists")
		}
	}

	// Prepare the statement.
	query := "UPDATE encounters SET title = $1, announcement = $2, notes = $3, first_encounter = $4, last_encounter = $5 WHERE id = $6 RETURNING id" // dev:finder+query
	stmtInsert, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtInsert.Close()
		if err != nil {
			db.Logger.Error("error closing stmtInsert", "error", err)
		}
	}()
	var encounterID int
	err = tx.StmtContext(ctx, stmtInsert).QueryRow(title, announcement, notes, first, last, id).Scan(&encounterID)
	if err != nil {
		db.Logger.Error("error on insert into encounters", "error", err.Error())
		return -1, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit encounters", "error", err.Error())
		return -1, err
	}
	return encounterID, nil
}

// deleteEncounterByID deletes an encounter by its ID.
func (db *DBX) DeleteEncounterByID(ctx context.Context, id int) error {
	// tx
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check if this encounter is associated with a stage
	var countStage int
	queryCheckStage := "SELECT COUNT(*) FROM stage_encounters WHERE encounter_id = $1" // dev:finder+query
	if err = tx.QueryRowContext(ctx, queryCheckStage, id).Scan(&countStage); err != nil {
		if err != sql.ErrNoRows {
			db.Logger.Error("no rows passed", "err", err.Error())
			return err
		}
	}
	if countStage > 0 {
		db.Logger.Error("encounter is associated with a stage")
		return fmt.Errorf("encounter is associated with a stage")
	}
	// check if this encounter is associated with a next encounter in auto play
	var countNext int
	queryCheckNext := "SELECT COUNT(*) FROM auto_play_next_encounter WHERE current_encounter_id = $1 OR next_encounter_id = $1" // dev:finder+query
	if err = tx.QueryRowContext(ctx, queryCheckNext, id).Scan(&countNext); err != nil {
		if err != sql.ErrNoRows {
			db.Logger.Error("no rows passed", "err", err.Error())
			return err
		}
	}
	if countNext > 0 {
		db.Logger.Error("encounter is associated with a next encounter in auto play")
		return fmt.Errorf("encounter is associated with a next encounter in auto play")
	}
	// Prepare the statement.
	query := "DELETE FROM encounters WHERE id = $1" // dev:finder+query
	stmtInsert, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmtInsert.Close()
		if err != nil {
			db.Logger.Error("error closing stmtInsert", "error", err)
		}
	}()
	_, err = tx.StmtContext(ctx, stmtInsert).Exec(id)
	if err != nil {
		db.Logger.Error("error on delete FROM encounters", "error", err.Error())
		return err
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit encounters", "error", err.Error())
		return err
	}
	return nil
}
