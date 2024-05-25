package db

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetEncounters(ctx context.Context) ([]types.Encounters, error) {
	list := []types.Encounters{}
	query := "SELECT id, title, notes, announcement, story_id, storyteller_id FROM encounters"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.logger.Error("query on encounters failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Encounters
		// var i int
		// var reward sql.NullString
		if err := rows.Scan(&s.ID, &s.Title, &s.Notes, &s.Announcement, &s.StoryID, &s.StorytellerID); err != nil {
			db.logger.Error("scan error on encounters", "error", err.Error())
		}
		// if reward.Valid {
		// 	s.Reward = reward.String
		// }
		// s.Phase = types.PhaseAtoi(i)
		list = append(list, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows err on encounters", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounters, error) {
	var encounters []types.Encounters
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, notes, announcement, storyteller_id FROM encounters WHERE story_id = $1", storyID)
	if err != nil {
		db.logger.Error("query on encounters by id failed", "error", err.Error())
		return encounters, err
	}
	defer rows.Close()
	for rows.Next() {
		// var i int
		var enc types.Encounters
		// var reward sql.NullString
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Notes, &enc.Announcement, &enc.StorytellerID); err != nil {
			db.logger.Error("scan error on encounters by id", "error", err.Error())
		}
		// if reward.Valid {
		// 	enc.Reward = reward.String
		// }
		// enc.Phase = types.PhaseAtoi(i)
		encounters = append(encounters, enc)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return encounters, nil
}

func (db *DBX) GetEncounterByID(ctx context.Context, id int) (types.Encounters, error) {
	var enc types.Encounters
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, announcement, notes, storyteller_id FROM encounters WHERE id = $1", id)
	if err != nil {
		db.logger.Error("query on encounters by id failed", "error", err.Error())
		return enc, err
	}
	defer rows.Close()
	for rows.Next() {
		// var i int
		// var reward sql.NullString
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Announcement, &enc.Notes, &enc.StorytellerID); err != nil {
			db.logger.Error("scan error on encounters by id", "error", err.Error())
		}
		// if reward.Valid {
		// 	enc.Reward = reward.String
		// }
		// enc.Phase = types.PhaseAtoi(i)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return enc, nil
}

func (db *DBX) CreateEncounter(ctx context.Context, title, announcement, notes string, storyID, storytellerID int) (int, error) {
	query := "INSERT INTO encounters(title, announcement, notes, story_id, storyteller_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into encounters failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(title, announcement, notes, storyID, storytellerID).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into encounters failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

// func (db *DBX) UpdatePhase(ctx context.Context, id, phase int) error {
// 	finished := false
// 	if phase == 3 {
// 		finished = true
// 	}
// 	a, err := db.Conn.ExecContext(ctx, "UPDATE encounters SET phase = $1, finished = $2 WHERE id = $3", phase, finished, id)
// 	if err != nil {
// 		db.logger.Error("update encounters failed", "error", err.Error())
// 		return err
// 	}
// 	r, err := a.RowsAffected()
// 	db.logger.Info("sql Rows Affected", "result", r)
// 	if r > 0 {
// 		return nil
// 	}
// 	return err
// }

func (db *DBX) AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error {

	query := "INSERT INTO encounters_participants_players (players_id, encounters_id) VALUES ($1, $2)"
	if npc {
		query = "INSERT INTO encounters_participants_non_players (non_players_id, encounters_id) VALUES ($1, $2)"
	}
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	for _, id := range players {
		_, err = tx.ExecContext(ctx, query, id, encounterID)
		if err != nil {
			return err
		}
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	// a, err := db.Conn.ExecContext(ctx, query, id, encounterID)
	// if err != nil {
	// 	db.logger.Error("add participants to encounter failed", "error", err.Error())
	// 	return err
	// }
	// r, err := a.RowsAffected()
	// db.logger.Info("add participant sql rows affected", "result", r)
	// if r > 0 {
	// 	return nil
	// }
	return err
}
