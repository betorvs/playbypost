package pg

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/lib/pq"
)

// create auto play
func (db *DBX) CreateAutoPlayTx(ctx context.Context, text string, storyID int, solo bool) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CreateAutoPlay failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	queryStoryKeys := "SELECT k.encoding_key FROM story AS s JOIN story_keys AS k ON s.id = k.story_id WHERE s.id = $1"
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryStoryKeys)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtStoryKeys.Close()
	var encodingKey string
	err = tx.StmtContext(ctx, stmtStoryKeys).QueryRow(storyID).Scan(&encodingKey)
	if err != nil {
		db.Logger.Error("query row SELECT story_keys and story failed", "error", err.Error())
		return -1, err
	}
	queryAutoPlay := "INSERT INTO auto_play(display_text, encoding_key, story_id, solo) VALUES($1, $2, $3, $4) RETURNING id"
	stmtInsert, err := db.Conn.PrepareContext(ctx, queryAutoPlay)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsert failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsert.Close()
	var autoPlayID int
	err = tx.StmtContext(ctx, stmtInsert).QueryRow(text, encodingKey, storyID, solo).Scan(&autoPlayID)
	if err != nil {
		db.Logger.Error("query row insert into auto_play failed", "error", err.Error())
		return -1, err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateAutoPlay failed", "error", err.Error())
		return -1, err
	}
	return autoPlayID, nil

}

// AddAutoPlayNextEncounter
// text string, autoPlayID, encounterID, nextEncounterID int,
func (db *DBX) AddAutoPlayNext(ctx context.Context, next []types.Next) error {

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
	// get encounters list
	queryEncounters := "SELECT e.id FROM auto_play AS ap JOIN encounters AS e ON ap.story_id = e.story_id WHERE ap.id = $1"
	stmtEncounters, err := db.Conn.PrepareContext(ctx, queryEncounters)
	if err != nil {
		db.Logger.Error("tx prepare on encounters failed", "error", err.Error())
		return err
	}
	defer stmtEncounters.Close()
	encountersID := []int{}
	rows, err := tx.StmtContext(ctx, stmtEncounters).Query(next[0].UpstreamID)
	if err != nil {
		db.Logger.Error("query on encounters failed", "error", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			db.Logger.Error("scan error on encounters ", "error", err.Error())
		}
		encountersID = append(encountersID, id)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on encounters", "error", err.Error())
	}
	// check if encounterID and nextEncounterID are in encountersID
	if !slices.Contains(encountersID, next[0].EncounterID) {
		return fmt.Errorf("encounterID not found")
	}
	if !slices.Contains(encountersID, next[0].NextEncounterID) {
		return fmt.Errorf("nextEncounterID not found")
	}

	for _, n := range next {
		query := "INSERT INTO auto_play_next_encounter(display_text, upstream_id, current_encounter_id, next_encounter_id) VALUES ($1, $2, $3, $4) RETURNING id"
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare on auto_play_next_encounter failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		var nextEncounterIDDB int
		err = tx.StmtContext(ctx, stmt).QueryRow(n.Text, n.UpstreamID, n.EncounterID, n.NextEncounterID).Scan(&nextEncounterIDDB)
		if err != nil {
			db.Logger.Error("error on insert into auto_play_next_encounter", "error", err.Error())
			return err
		}
		db.Logger.Debug("adding auto play objectives", "nextEncounterIDDB", nextEncounterIDDB, "values", n.Objective.Values)
		// insert into auto_play_next_objectives
		query = "INSERT INTO auto_play_next_objectives (upstream_id, kind, values) VALUES ($1, $2, $3) RETURNING id"
		stmt, err = db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare on auto_play_next_objectives failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		var objectiveID int
		err = tx.StmtContext(ctx, stmt).QueryRow(nextEncounterIDDB, n.Objective.Kind, pq.Array(n.Objective.Values)).Scan(&objectiveID)
		if err != nil {
			db.Logger.Error("error on insert into auto_play_next_objectives", "error", err.Error())
			return err
		}
		db.Logger.Debug("auto play next objective added", "next_id", nextEncounterIDDB, "objective_id", objectiveID)
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit auto_play_next_encounter", "error", err.Error())
		return err
	}
	return nil
}

// create auto play
func (db *DBX) CreateAutoPlayChannelTx(ctx context.Context, channelID, userID string, autoPlayID int) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CreateAutoPlayChannel failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check if autoPlayID exists
	queryAutoPlay := "SELECT id FROM auto_play WHERE id = $1"
	stmt, err := db.Conn.PrepareContext(ctx, queryAutoPlay)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var apID int
	err = tx.StmtContext(ctx, stmt).QueryRow(autoPlayID).Scan(&apID)
	if err != nil {
		return -1, fmt.Errorf("auto play id not found")
	}
	if apID == 0 {
		db.Logger.Error("auto play id not found", "error", err.Error())
		return -1, fmt.Errorf("auto play id not found")
	}

	// check if channel already in use
	query := "SELECT id FROM auto_play_channel WHERE active = 'true' AND channel = $1 AND upstream_id = $2"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare SELECT on auto_play_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var autoPlayChannelid int
	err = tx.StmtContext(ctx, stmt).QueryRow(channelID, apID).Scan(&autoPlayChannelid)
	if err == nil {
		return -1, fmt.Errorf("channel already in use")
	}
	// query for encounter_id
	queryEncounter := "SELECT id FROM encounters WHERE story_id = (SELECT story_id FROM auto_play WHERE id = $1) AND first_encounter = true"
	stmtEncounter, err := db.Conn.PrepareContext(ctx, queryEncounter)
	if err != nil {
		db.Logger.Error("tx prepare on encounters failed", "error", err.Error())
		return -1, err
	}
	defer stmtEncounter.Close()
	var encounterID int
	err = tx.StmtContext(ctx, stmtEncounter).QueryRow(autoPlayID).Scan(&encounterID)
	if err != nil {
		db.Logger.Error("query row SELECT encounters failed", "error", err.Error())
		return -1, err
	}
	// encounter id should be greater than 0
	if encounterID == 0 {
		db.Logger.Error("first encounter id not found", "error", err.Error())
		return -1, fmt.Errorf("first encounter id not found")
	}

	// insert into auto_play_channel
	query = "INSERT INTO auto_play_channel(channel, upstream_id, active) VALUES($1, $2, $3) RETURNING id"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare insert on auto_play_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	err = tx.StmtContext(ctx, stmt).QueryRow(channelID, apID, true).Scan(&autoPlayChannelid)
	if err != nil {
		db.Logger.Error("tx insert on auto_play_channel failed", "error", err.Error())
		return -1, err
	}
	// insert into users
	// check if user exists
	// check user exist
	queryUser := "SELECT id FROM users WHERE userid = $1"
	stmtQueryUser, err := db.Conn.PrepareContext(ctx, queryUser)
	if err != nil {
		db.Logger.Error("tx prepare on queryUser failed", "error", err.Error())
		return -1, err
	}
	defer stmtQueryUser.Close()
	var validUserID int
	err = tx.StmtContext(ctx, stmtQueryUser).QueryRow(userID).Scan(&validUserID)
	if err != nil {
		db.Logger.Debug("user not found", "return", err.Error())
		// just log this error
		// return -1, err

	}
	if validUserID == 0 {
		id, err := db.createUser(ctx, userID, tx)
		if err != nil {
			db.Logger.Error("insert into users failed", "error", err.Error())
			return -1, err
		}
		validUserID = id
	}

	// add user id to auto_play_group
	query = "INSERT INTO auto_play_group(user_id, upstream_id, last_update_at, active) VALUES($1, $2, $3, $4)"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare insert into auto_play_group failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, validUserID, autoPlayChannelid, time.Now(), true)
	if err != nil {
		db.Logger.Error("tx insert into auto_play_group failed", "error", err.Error())
		return -1, err
	}

	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateAutoPlayChannel failed", "error", err.Error())
		return -1, err
	}
	return encounterID, nil
}

// add registry to auto_play_encounter_activities
func (db *DBX) RegisterActivitiesAutoPlay(ctx context.Context, autoPlayID, encounterID int, actions types.Actions) error {
	query := "INSERT INTO auto_play_encounter_activities(upstream_id, encounter_id, actions) VALUES($1, $2, $3)"
	_, err := db.Conn.ExecContext(ctx, query, autoPlayID, encounterID, actions)
	if err != nil {
		db.Logger.Error("error on insert into auto_play_encounter_activities", "error", err.Error(), "upstream_id", autoPlayID, "encounter_id", encounterID, "actions", actions)
		return err
	}
	return nil
}

// update processed activity
func (db *DBX) UpdateProcessedAutoPlay(ctx context.Context, id int, processed bool, actions types.Actions) error {
	query := "UPDATE auto_play_encounter_activities SET processed = $1, actions = $2 WHERE id = $3"
	_, err := db.Conn.ExecContext(ctx, query, processed, actions, id)
	if err != nil {
		db.Logger.Error("error on update auto_play_encounter_activities", "error", err.Error())
		return err
	}
	return nil
}

func (db *DBX) UpdateAutoPlayGroup(ctx context.Context, id, count int, date time.Time) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on update auto_play_group failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	query := "UPDATE auto_play_group SET last_update_at=$1, interactions=$2 WHERE id=$3"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare update auto_play_group failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).Exec(date, count, id)
	if err != nil {
		db.Logger.Error("tx update auto_play_group failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on update auto_play_group failed", "error", err.Error())
		return err
	}
	return nil
}

// update auto_play_state
func (db *DBX) UpdateAutoPlayState(ctx context.Context, autoPlayChannel string, encounterID int) error {
	db.Logger.Debug("update auto play state", "autoPlayChannelID", autoPlayChannel, "encounterID", encounterID)
	// start tx
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on UpdateAutoPlayState failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// find auto_play_channel_id
	queryAutoPlayChannel := "SELECT id FROM auto_play_channel WHERE active = 'true' AND channel = $1"
	stmt, err := db.Conn.PrepareContext(ctx, queryAutoPlayChannel)
	if err != nil {
		db.Logger.Error("tx prepare SELECT on queryAutoPlayChannel failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var autoPlayChannelID int
	err = tx.StmtContext(ctx, stmt).QueryRow(autoPlayChannel).Scan(&autoPlayChannelID)
	if err != nil {
		db.Logger.Error("auto play channel id not found", "error", err.Error())
		return err
	}

	// check if autoPlayID exists
	queryAutoPlayState := "SELECT id FROM auto_play_state WHERE active = 'true' AND upstream_id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, queryAutoPlayState)
	if err != nil {
		db.Logger.Error("tx prepare SELECT on queryAutoPlayState failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var apsID int
	err = tx.StmtContext(ctx, stmt).QueryRow(autoPlayChannelID).Scan(&apsID)
	if err != nil {
		db.Logger.Debug("auto play state id not found")
		// return fmt.Errorf("auto play state id not found")
	}
	if apsID == 0 {
		db.Logger.Debug("auto play state id not found")
		query := "INSERT INTO auto_play_state(upstream_id, encounter_id, active) VALUES($1, $2, $3)"
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare to insert into auto_play_state failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, autoPlayChannelID, encounterID, true)
		if err != nil {
			db.Logger.Error("tx insert into auto_play_state failed", "error", err.Error())
			return err
		}

	} else {
		db.Logger.Debug("auto play state", "apsID", apsID, "encounterID", encounterID, "autoPlayChannelID", autoPlayChannelID)
		query := "UPDATE auto_play_state SET encounter_id = $1 WHERE id = $2"
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare to update on auto_play_state failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, encounterID, apsID)
		if err != nil {
			db.Logger.Error("tx update on auto_play_state failed", "error", err.Error())
			return err
		}
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on UpdateAutoPlayState failed", "error", err.Error())
		return err
	}

	return nil
}

func (db *DBX) CloseAutoPlayChannel(ctx context.Context, channelID string, autoPlayID int) error {
	// tx start
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CloseAutoPlayChannel failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// get channel id
	query := "SELECT id FROM auto_play_channel WHERE active = 'true' AND channel = $1 AND upstream_id = $2"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_channel failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var autoPlayChannelID int
	err = tx.StmtContext(ctx, stmt).QueryRow(channelID, autoPlayID).Scan(&autoPlayChannelID)
	if err != nil {
		db.Logger.Error("tx query on auto_play_channel failed", "error", err.Error())
		return err
	}

	query = "UPDATE auto_play_channel SET active = false WHERE id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_channel failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, autoPlayChannelID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_channel failed", "error", err.Error())
		return err
	}
	// close group
	query = "UPDATE auto_play_group SET active = false WHERE upstream_id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_group failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, autoPlayChannelID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_group failed", "error", err.Error())
		return err
	}
	// close state
	query = "UPDATE auto_play_state SET active = false WHERE upstream_id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_state failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, autoPlayChannelID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_state failed", "error", err.Error())
		return err
	}

	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CloseAutoPlayChannel failed", "error", err.Error())
		return err
	}

	return nil
}

func (db *DBX) DeleteAutoPlayNextEncounter(ctx context.Context, id int) error {
	// start tx
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on DeleteNextEncounter failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// SELECT ids to delete
	query := "SELECT a.id, apno.id FROM auto_play_next_encounter AS a JOIN auto_play_next_objectives AS apno ON apno.upstream_id = a.id WHERE a.id = $1"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_next_encounter failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var nextID, objectiveID int
	err = tx.StmtContext(ctx, stmt).QueryRow(id).Scan(&nextID, &objectiveID)
	if err != nil {
		db.Logger.Error("tx query on auto_play_next_encounter failed", "error", err.Error())
		return err
	}
	// delete FROM auto_play_next_objectives
	query = "DELETE FROM auto_play_next_objectives WHERE id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_next_objectives failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, objectiveID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_next_objectives failed", "error", err.Error())
		return err
	}
	// delete FROM auto_play_next_encounter
	query = "DELETE FROM auto_play_next_encounter WHERE id = $1"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_next_encounter failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, nextID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_next_encounter failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on DeleteNextEncounter failed", "error", err.Error())
		return err
	}

	return nil
}
