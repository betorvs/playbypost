package pg

import (
	"context"
	"database/sql"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/betorvs/playbypost/core/utils"
	"github.com/lib/pq"
)

func (db *DBX) GetAutoPlay(ctx context.Context) ([]types.AutoPlay, error) {
	autoPlay := []types.AutoPlay{}
	query := "SELECT id, display_text, story_id, solo FROM auto_play"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on auto_play failed", "error", err.Error())
		return autoPlay, err
	}
	defer rows.Close()
	for rows.Next() {
		var auto types.AutoPlay

		if err := rows.Scan(&auto.ID, &auto.Text, &auto.StoryID, &auto.Solo); err != nil {
			db.Logger.Error("scan error on auto_play ", "error", err.Error())
		}
		autoPlay = append(autoPlay, auto)

	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play", "error", err.Error())
	}
	return autoPlay, nil
}

func (db *DBX) GetAutoPlayByID(ctx context.Context, autoPlayID int) (types.AutoPlay, error) {
	autoPlay := types.AutoPlay{}
	query := "SELECT id, display_text, story_id, solo FROM	auto_play WHERE id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, autoPlayID)
	if err != nil {
		db.Logger.Error("query on auto_play by id failed", "error", err.Error())
		return autoPlay, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&autoPlay.ID, &autoPlay.Text, &autoPlay.StoryID, &autoPlay.Solo); err != nil {
			db.Logger.Error("scan error on auto_play by id ", "error", err.Error())
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play by id", "error", err.Error())
	}
	return autoPlay, nil
}

func (db *DBX) GetAutoPlayEncounterListByStoryID(ctx context.Context, storyID int) (types.EncounterList, error) {
	list := types.EncounterList{}
	query := "SELECT a.id, e.title AS encounter, e.id AS encounter_id, n.title AS next_encounter, n.id AS next_id FROM auto_play_next_encounter AS a JOIN encounters AS e ON e.id = a.current_encounter_id JOIN encounters AS n ON n.id = a.next_encounter_id WHERE e.story_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, storyID)
	if err != nil {
		db.Logger.Error("query on auto_play_next_encounter by story_id failed", "error", err.Error())
		return list, err
	}
	defer rows.Close()
	for rows.Next() {
		var next types.EncounterWithNext
		if err := rows.Scan(&next.ID, &next.Encounter, &next.EncounterID, &next.NextEncounter, &next.NextID); err != nil {
			db.Logger.Error("scan error on auto_play_next_encounter by story_id ", "error", err.Error())
		}
		list.Link = append(list.Link, next)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play_next_encounter by story_id", "error", err.Error())
	}
	queryEncounter := "SELECT id, title AS name FROM encounters WHERE story_id = $1"
	rowsEncounter, err := db.Conn.QueryContext(ctx, queryEncounter, storyID)
	if err != nil {
		db.Logger.Error("query on encounters by story_id failed", "error", err.Error())
		return list, err
	}
	defer rowsEncounter.Close()
	for rowsEncounter.Next() {
		var generic types.Options
		if err := rowsEncounter.Scan(&generic.ID, &generic.Name); err != nil {
			db.Logger.Error("scan error on encounters by story_id ", "error", err.Error())
		}
		list.EncounterList = append(list.EncounterList, generic)
	}
	return list, nil
}

// GetAutoPlayByChannelID returns auto play by channel id
// case channel id is not found, return nil, []AutoPlayStartOptions, error
// case channel id is found, return auto play, []AutoPlayStartOptions, error
func (db *DBX) GetAutoPlayOptionsByChannelID(ctx context.Context, channelID, userID string) (types.AutoPlayOptions, error) {
	autoPlay := types.AutoPlayOptions{}

	query := `SELECT ap.id, ap.display_text, ap.story_id, ap.solo, ap.encoding_key, 
	ac.id AS auto_play_channel_id, ac.channel AS channel_id, 
	apg.id, apg.user_id, apg.last_update_at, apg.interactions, apne.id, 
	apne.upstream_id, apne.display_text, apne.current_encounter_id, apne.next_encounter_id,
	apno.kind, apno.values 
	FROM auto_play_channel AS ac 
	JOIN auto_play AS ap ON ap.id = ac.upstream_id 
	JOIN auto_play_state AS aps ON aps.upstream_id = ac.id 
	JOIN auto_play_group AS apg ON apg.upstream_id = ac.id 
	JOIN auto_play_next_encounter AS apne ON apne.upstream_id = ap.id 
	JOIN auto_play_next_objectives AS apno ON apno.upstream_id = apne.id
	WHERE ac.active = 'true' AND apg.active = 'true' AND aps.active = 'true' AND apne.current_encounter_id = aps.encounter_id AND ac.channel = $1`
	rows, err := db.Conn.QueryContext(ctx, query, channelID)
	if err != nil {
		db.Logger.Error("query on auto_play_channel by channel_id failed", "error", err.Error())
		return autoPlay, err
	}
	defer rows.Close()
	for rows.Next() {
		var group types.AutoPlayGroup
		var next types.Next
		var values []sql.NullInt64
		if err := rows.Scan(&autoPlay.AutoPlay.ID, &autoPlay.AutoPlay.Text, &autoPlay.AutoPlay.StoryID, &autoPlay.AutoPlay.Solo, &autoPlay.EncodingKey, &autoPlay.AutoPlayChannelID, &autoPlay.ChannelID, &group.ID, &group.UserID, &group.LastUpdateAt, &group.Interactions, &next.ID, &next.UpstreamID, &next.Text, &next.EncounterID, &next.NextEncounterID, &next.Objective.Kind, pq.Array(&values)); err != nil {
			db.Logger.Error("scan error on auto_play_channel by channel_id ", "error", err.Error())
		}
		if len(values) > 0 {
			for _, v := range values {
				if v.Valid {
					next.Objective.Values = append(next.Objective.Values, int(v.Int64))
				}
			}
		}
		if group.UserID == userID {
			autoPlay.Group = append(autoPlay.Group, group)
		}
		autoPlay.NextEncounters = append(autoPlay.NextEncounters, next)

	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play_channel by channel_id", "error", err.Error())
	}
	return autoPlay, nil
}

func (db *DBX) GetAutoPlayActivities(ctx context.Context) ([]types.Activity, error) {
	autoPlay := []types.Activity{}
	query := "SELECT id, upstream_id, encounter_id, actions, processed FROM auto_play_encounter_activities"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on auto_play_encounter_activities failed", "error", err.Error())
		return autoPlay, err
	}
	defer rows.Close()
	for rows.Next() {
		var auto types.Activity
		if err := rows.Scan(&auto.ID, &auto.UpstreamID, &auto.EncounterID, &auto.Actions, &auto.Processed); err != nil {
			db.Logger.Error("scan error on auto_play_encounter_activities ", "error", err.Error())
		}
		autoPlay = append(autoPlay, auto)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play_encounter_activities", "error", err.Error())
	}
	return autoPlay, nil
}

// get encounter by encounter_id and auto_play_id
func (db *DBX) GetAnnounceByEncounterID(ctx context.Context, encounterID, autoPlayID int) (string, bool, error) {
	query := "SELECT ap.encoding_key, e.announcement, e.last_encounter FROM auto_play AS ap JOIN encounters AS e ON ap.story_id = e.story_id WHERE e.id = $1 AND ap.id = $2"
	var encodingKey, encAnnounce string
	var last bool
	err := db.Conn.QueryRowContext(ctx, query, encounterID, autoPlayID).Scan(&encodingKey, &encAnnounce, &last)
	if err != nil {
		db.Logger.Error("query row SELECT auto_play_next_encounter failed", "error", err.Error())
		return "", false, err
	}
	text, err := utils.DecryptText(encAnnounce, encodingKey)
	if err != nil {
		db.Logger.Error("error on decrypt text", "error", err.Error())
		return "", false, err
	}
	return text, last, nil
}

// func GetNextEncounterByAutoPlayID
func (db *DBX) GetNextEncounterByAutoPlayID(ctx context.Context, autoPlayID int) ([]types.Next, error) {
	next := []types.Next{}
	query := "SELECT a.id, a.upstream_id, a.current_encounter_id, a.next_encounter_id, a.display_text, apno.kind, apno.values FROM auto_play_next_encounter AS a JOIN auto_play_next_objectives AS apno ON apno.upstream_id = a.id WHERE a.upstream_id = $1"
	rows, err := db.Conn.QueryContext(ctx, query, autoPlayID)
	if err != nil {
		db.Logger.Error("query on auto_play_next_encounter failed", "error", err.Error())
		return next, err
	}
	defer rows.Close()
	for rows.Next() {
		var n types.Next
		var o types.Objective
		var values []sql.NullInt64
		if err := rows.Scan(&n.ID, &n.UpstreamID, &n.EncounterID, &n.NextEncounterID, &n.Text, &o.Kind, pq.Array(&values)); err != nil {
			db.Logger.Error("scan error on auto_play_next_encounter ", "error", err.Error())
		}
		n.Objective = o
		if len(values) > 0 {
			for _, v := range values {
				if v.Valid {
					n.Objective.Values = append(n.Objective.Values, int(v.Int64))
				}
			}
		}

		next = append(next, n)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on auto_play_next_encounter", "error", err.Error())
	}
	return next, nil
}
