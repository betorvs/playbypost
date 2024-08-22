package pg

import (
	"context"
	"fmt"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
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
	queryStoryKeys := "select k.encoding_key from story AS s JOIN story_keys AS k ON s.id = k.story_id WHERE s.id = $1"
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryStoryKeys)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtStoryKeys.Close()
	var encodingKey string
	err = tx.StmtContext(ctx, stmtStoryKeys).QueryRow(storyID).Scan(&encodingKey)
	if err != nil {
		db.Logger.Error("query row select story_keys and story failed", "error", err.Error())
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
func (db *DBX) AddAutoPlayNext(ctx context.Context, text string, autoPlayID, encounterID, nextEncounterID int) error {

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
	rows, err := tx.StmtContext(ctx, stmtEncounters).Query(autoPlayID)
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
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on encounters", "error", err.Error())
	}
	// check if encounterID and nextEncounterID are in encountersID
	if !slices.Contains(encountersID, encounterID) {
		return fmt.Errorf("encounterID not found")
	}
	if !slices.Contains(encountersID, nextEncounterID) {
		return fmt.Errorf("nextEncounterID not found")
	}

	query := "INSERT INTO auto_play_next_encounter (display_text, auto_play_id, current_encounter_id, next_encounter_id) VALUES ($1, $2, $3, $4)"
	_, err = tx.ExecContext(ctx, query, text, autoPlayID, encounterID, nextEncounterID)
	if err != nil {
		db.Logger.Error("error on insert into auto_play_next_encounter", "error", err.Error())
		return err
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
	query := "SELECT id FROM auto_play_channel WHERE channel = $1 AND auto_play_id = $2"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare select on auto_play_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var id int
	err = tx.StmtContext(ctx, stmt).QueryRow(channelID, apID).Scan(&id)
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
		db.Logger.Error("query row select encounters failed", "error", err.Error())
		return -1, err
	}
	// encounter id should be greater than 0
	if encounterID == 0 {
		db.Logger.Error("first encounter id not found", "error", err.Error())
		return -1, fmt.Errorf("first encounter id not found")
	}

	// insert into auto_play_channel
	query = "INSERT INTO auto_play_channel(channel, auto_play_id, active) VALUES($1, $2, $3)"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare insert on auto_play_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, channelID, apID, true)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_channel failed", "error", err.Error())
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
		db.Logger.Info("user not found", "return", err.Error())
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
	query = "INSERT INTO auto_play_group(user_id, auto_play_id) VALUES($1, $2)"
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play_group failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, validUserID, apID)
	if err != nil {
		db.Logger.Error("tx exec on auto_play_group failed", "error", err.Error())
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
	query := "INSERT INTO auto_play_encounter_activities(auto_play_id, encounter_id, actions) VALUES($1, $2, $3)"
	_, err := db.Conn.ExecContext(ctx, query, autoPlayID, encounterID, actions)
	if err != nil {
		db.Logger.Error("error on insert into auto_play_encounter_activities", "error", err.Error(), "autoPlayID", autoPlayID, "encounter_id", encounterID, "actions", actions)
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

// update auto_play_state
func (db *DBX) UpdateAutoPlayState(ctx context.Context, autoPlayID int, encounterID int) error {
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
	// check if autoPlayID exists
	queryAutoPlay := "SELECT id FROM auto_play WHERE id = $1"
	stmt, err := db.Conn.PrepareContext(ctx, queryAutoPlay)
	if err != nil {
		db.Logger.Error("tx prepare on auto_play failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var apID int
	err = tx.StmtContext(ctx, stmt).QueryRow(autoPlayID).Scan(&apID)
	if err != nil {
		return fmt.Errorf("auto play id not found")
	}
	if apID == 0 {
		db.Logger.Info("auto play id not found")
		query := "INSERT INTO auto_play_state(auto_play_id, encounter_id) VALUES($1, $2)"
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare on auto_play_state failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, autoPlayID, encounterID)
		if err != nil {
			db.Logger.Error("tx exec on auto_play_state failed", "error", err.Error())
			return err
		}

	} else {
		query := "UPDATE auto_play_state SET encounter_id = $1 WHERE auto_play_id = $2"
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("tx prepare on auto_play_state failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, encounterID, autoPlayID)
		if err != nil {
			db.Logger.Error("tx exec on auto_play_state failed", "error", err.Error())
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
	query := "UPDATE auto_play_channel SET active = false WHERE channel = $1 AND auto_play_id = $2"
	_, err := db.Conn.ExecContext(ctx, query, channelID, autoPlayID)
	if err != nil {
		db.Logger.Error("error on update auto_play_channel", "error", err.Error())
		return err
	}
	return nil
}
