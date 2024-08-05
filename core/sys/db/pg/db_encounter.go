package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetEncounters(ctx context.Context) ([]types.Encounter, error) {
	list := []types.Encounter{}
	query := "SELECT id, title, notes, announcement, story_id, writer_id FROM encounters"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on encounters failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Encounter
		// var i int
		// var reward sql.NullString
		if err := rows.Scan(&s.ID, &s.Title, &s.Notes, &s.Announcement, &s.StoryID, &s.WriterID); err != nil {
			db.Logger.Error("scan error on encounters", "error", err.Error())
		}
		// if reward.Valid {
		// 	s.Reward = reward.String
		// }
		// s.Phase = types.PhaseAtoi(i)
		list = append(list, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on encounters", "error", err.Error())
	}
	return list, nil
}

func (db *DBX) GetEncounterByStoryID(ctx context.Context, storyID int) ([]types.Encounter, error) {
	var encounters []types.Encounter
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, notes, announcement, story_id, writer_id FROM encounters WHERE story_id = $1", storyID)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return encounters, err
	}
	defer rows.Close()
	for rows.Next() {
		// var i int
		var enc types.Encounter
		// var reward sql.NullString
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Notes, &enc.Announcement, &enc.StoryID, &enc.WriterID); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
		// if reward.Valid {
		// 	enc.Reward = reward.String
		// }
		// enc.Phase = types.PhaseAtoi(i)
		encounters = append(encounters, enc)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return encounters, nil
}

func (db *DBX) GetEncounterByID(ctx context.Context, id int) (types.Encounter, error) {
	var enc types.Encounter
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, announcement, notes, writer_id FROM encounters WHERE id = $1", id)
	if err != nil {
		db.Logger.Error("query on encounters by id failed", "error", err.Error())
		return enc, err
	}
	defer rows.Close()
	for rows.Next() {
		// var i int
		// var reward sql.NullString
		if err := rows.Scan(&enc.ID, &enc.Title, &enc.Announcement, &enc.Notes, &enc.WriterID); err != nil {
			db.Logger.Error("scan error on encounters by id", "error", err.Error())
		}
		// if reward.Valid {
		// 	enc.Reward = reward.String
		// }
		// enc.Phase = types.PhaseAtoi(i)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on encounters by id", "error", err.Error())
	}
	return enc, nil
}

func (db *DBX) CreateEncounter(ctx context.Context, title, announcement, notes string, storyID, storytellerID int) (int, error) {
	query := "INSERT INTO encounters(title, announcement, notes, story_id, writer_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into encounters failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(title, announcement, notes, storyID, storytellerID).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into encounters failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}
