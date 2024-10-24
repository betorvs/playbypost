package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetEncounters(ctx context.Context) ([]types.Encounter, error) {
	list := []types.Encounter{}
	query := "SELECT id, title, notes, announcement, story_id, writer_id, first_encounter, last_encounter FROM encounters"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on encounters failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Encounter
		if err := rows.Scan(&s.ID, &s.Title, &s.Notes, &s.Announcement, &s.StoryID, &s.WriterID, &s.FirstEncounter, &s.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters", "error", err.Error())
		}
		list = append(list, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on encounters", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error) {
	encounters := []types.Encounter{}
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, notes, announcement, story_id, writer_id, first_encounter, last_encounter FROM encounters WHERE story_id = $1", storyID)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return encounters, err
	}
	defer rows.Close()
	for rows.Next() {
		var enc types.Encounter
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Notes, &enc.Announcement, &enc.StoryID, &enc.WriterID, &enc.FirstEncounter, &enc.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
		encounters = append(encounters, enc)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return encounters, nil
}

func (db *DBX) GetEncounterByID(ctx context.Context, id int) (types.Encounter, error) {
	enc := types.Encounter{}
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, announcement, notes, story_id, writer_id, first_encounter, last_encounter FROM encounters WHERE id = $1", id)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return enc, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Announcement, &enc.Notes, &enc.StoryID, &enc.WriterID, &enc.FirstEncounter, &enc.LastEncounter); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
	}
	// Check for errors from iterating over rows.
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
		queryCheck := "SELECT COUNT(*) FROM encounters WHERE story_id = $1 AND first_encounter = true"
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
	query := "INSERT INTO encounters(title, announcement, notes, story_id, writer_id, first_encounter, last_encounter) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	stmtInsert, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsert.Close()
	var encounterID int
	err = tx.StmtContext(ctx, stmtInsert).QueryRow(title, announcement, notes, storyID, storytellerID, first, last).Scan(&encounterID)
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
		queryCheck := "SELECT id FROM encounters WHERE story_id = $1 AND first_encounter = true"
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
	query := "Update encounters SET title = $1, announcement = $2, notes = $3, first_encounter = $4, last_encounter = $5 WHERE id = $6 RETURNING id"
	stmtInsert, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsert.Close()
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
