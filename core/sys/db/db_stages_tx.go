package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateStageTx(ctx context.Context, text, userid string, storyID int) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.logger.Error("tx begin on CreateStageTx failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	// check user exist
	queryUser := "SELECT id FROM users WHERE userid = $1"
	stmtQueryUser, err := db.Conn.PrepareContext(ctx, queryUser)
	if err != nil {
		db.logger.Error("tx prepare on queryUser failed", "error", err.Error())
		return -1, err
	}
	defer stmtQueryUser.Close()
	var userID int
	err = tx.StmtContext(ctx, stmtQueryUser).QueryRow(userid).Scan(&userID)
	if err != nil {
		db.logger.Info("user not found", "return", err.Error())
		// just log this error
		// return -1, err

	}
	// insert user if it does not exist
	if userID == 0 {
		// queryInsertUser := "INSERT INTO users(userid) VALUES($1) RETURNING id"
		// stmtInsertUser, err := db.Conn.PrepareContext(ctx, queryInsertUser)
		// if err != nil {
		// 	db.logger.Error("tx prepare on story_keys failed", "error", err.Error())
		// 	return -1, err
		// }
		// defer stmtInsertUser.Close()
		// err = tx.StmtContext(ctx, stmtInsertUser).QueryRow(userid).Scan(&userID)
		// if err != nil {
		// 	db.logger.Error("query row insert into users failed", "error", err.Error())
		// 	return -1, err
		// }
		id, err := db.createUser(ctx, userid, tx)
		if err != nil {
			db.logger.Error("insert into users failed", "error", err.Error())
			return -1, err
		}
		userID = id
	}

	queryStoryKeys := "select k.encoding_key from story AS s JOIN story_keys AS k ON s.id = k.story_id WHERE s.id = $1"
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryStoryKeys)
	if err != nil {
		db.logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtStoryKeys.Close()
	var encodingKey string
	err = tx.StmtContext(ctx, stmtStoryKeys).QueryRow(storyID).Scan(&encodingKey)
	if err != nil {
		db.logger.Error("query row select story_keys and story failed", "error", err.Error())
		return -1, err
	}

	// create stage
	queryInsertStage := "INSERT INTO stage(display_text, encoding_key, finished, storyteller_id, story_id) VALUES($1, $2, $3, $4, $5) RETURNING id"
	stmtInsertStage, err := db.Conn.PrepareContext(ctx, queryInsertStage)
	if err != nil {
		db.logger.Error("tx prepare on stmtInsertStage failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsertStage.Close()
	var stageID int
	err = tx.StmtContext(ctx, stmtInsertStage).QueryRow(text, encodingKey, false, userID, storyID).Scan(&stageID)
	if err != nil {
		db.logger.Error("query row insert into stage failed", "error", err.Error())
		return -1, err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.logger.Error("tx commit on CreateStageTx failed", "error", err.Error())
		return -1, err
	}

	return stageID, nil
}

func (db *DBX) AddChannelToStage(ctx context.Context, channel string, id int) (int, error) {
	query := "INSERT INTO stage_channel(channel, stage_id, active) VALUES($1, $2, $3) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(channel, id, true).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error) {
	query := "INSERT INTO stage_encounters(display_text, stage_id, storyteller_id, encounters_id) VALUES($1, $2, $3, $4) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(text, stage_id, storyteller_id, encounter_id).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) UpdatePhase(ctx context.Context, id, phase int) error {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.logger.Error("tx begin on UpdatePhase failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	if phase == int(types.Running) {
		// check if on this stage we have any stage_encounter already in running stage
		query := "SELECT (s.phase >= 0) from stage_encounters AS s JOIN stage_encounters AS se ON se.stage_id = s.stage_id WHERE se.id = $1 AND s.phase = $2"
		var running bool
		p := int(types.Running)
		if err = tx.QueryRowContext(ctx, query, id, p).Scan(&running); err != nil {
			if err != sql.ErrNoRows {
				db.logger.Error("no rows passed", "err", err.Error())
				return err
			}
		}
		if running {
			return fmt.Errorf("stage_id already have a stage_encounters in running state")
		}
	}

	// change state of stage_encounter
	queryUpdateStageEncounter := "UPDATE stage_encounters SET phase = $1, updated_at = NOW() WHERE id = $2 RETURNING id"
	stmtUpdateStageEncounter, err := db.Conn.PrepareContext(ctx, queryUpdateStageEncounter)
	if err != nil {
		db.logger.Error("tx prepare on stmtUpdateStageEncounter failed", "error", err.Error())
		return err
	}
	defer stmtUpdateStageEncounter.Close()
	var ID int
	err = tx.StmtContext(ctx, stmtUpdateStageEncounter).QueryRow(phase, id).Scan(&ID)
	if err != nil {
		db.logger.Error("query row insert into StageEncounter failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.logger.Error("tx commit on UpdatePhase failed", "error", err.Error())
		return err
	}
	db.logger.Info("stage_encounter changed", "id", id)
	return nil
}

func (db *DBX) AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error {

	query := "INSERT INTO stage_encounters_participants_players (players_id, stage_encounters_id) VALUES ($1, $2)"
	if npc {
		query = "INSERT INTO stage_encounters_participants_non_players (non_players_id, stage_encounters_id) VALUES ($1, $2)"
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
	return nil
}

// stage_next_encounter
func (db *DBX) AddNextEncounter(ctx context.Context, text string, stageID, encounterID, nextEncounterID int) error {
	query := "INSERT INTO stage_next_encounter (display_text, stage_id, current_encounter_id, next_encounter_id) VALUES ($1, $2, $3, $4)"
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, text, stageID, encounterID, nextEncounterID)
	if err != nil {
		db.logger.Error("error on insert into stage_next_encounter", "error", err.Error())
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.logger.Error("error on commit stage_next_encounter", "error", err.Error())
		return err
	}
	return nil
}

// stage_running_tasks
func (db *DBX) AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error {
	query := "INSERT INTO stage_running_tasks (display_text, stage_id, task_id, storyteller_id, stage_encounters_id) VALUES ($1, $2, $3, $4, $5)"
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, text, stageID, taskID, StorytellerID, encounterID)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// stage_encounter_activities
func (db *DBX) AddEncounterActivities(ctx context.Context, text string, stageID, encounterID int) error {
	query := "INSERT INTO stage_encounter_activities (notes, stage_id, encounter_id) VALUES ($1, $2, $3)"
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, text, stageID, encounterID)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// func (db *DBX) GetParticipants(ctx context.Context, encounterID int, npc bool, players []int) ([]types.Players, error) {

// 	enc := types.StageEncounter{}
// 	query := "select se.ID, se.display_text, e.title, se.storyteller_id, e.notes, e.announcement, s.encoding_key from stage_encounters AS se JOIN encounters AS e ON se.encounters_id = e.id JOIN stage AS s ON se.stage_id = s.id WHERE se.ID = $1"
// 	rows, err := db.Conn.QueryContext(ctx, query, id)
// 	if err != nil {
// 		db.logger.Error("query on stage_encounters by stage_encounters.ID failed", "error", err.Error())
// 		return enc, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var s types.StageEncounter
// 		var notes, announce, encodingKey string
// 		if err := rows.Scan(&s.ID, &s.Text, &s.Title, &s.StorytellerID, &notes, &announce, &encodingKey); err != nil {
// 			db.logger.Error("scan error on stage_encounters by stage_encounters.ID ", "error", err.Error())
// 		}
// 		s.Notes, _ = utils.DecryptText(notes, encodingKey)
// 		s.Announcement, _ = utils.DecryptText(announce, encodingKey)
// 		enc = s
// 	}
// 	// Check for errors from iterating over rows.
// 	if err := rows.Err(); err != nil {
// 		db.logger.Error("rows err on stage_encounters by stage_encounters.ID", "error", err.Error())
// 	}
// 	return enc, nil
// }

// stage_encounter_activities
func (db *DBX) RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error {
	db.logger.Info("RegisterActivities", "stageID", stageID, "encounterID", encounterID, "actions", actions)
	query := "INSERT INTO stage_encounter_activities (actions, stage_id, encounter_id) VALUES ($1, $2, $3)"
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, actions, stageID, encounterID)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *DBX) UpdateProcessedActivities(ctx context.Context, id int, processed bool, actions types.Actions) error {
	query := "UPDATE stage_encounter_activities SET processed = $1, actions = $2 WHERE id = $3"
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, query, processed, actions, id)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
