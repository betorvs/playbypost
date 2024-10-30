package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
	"github.com/lib/pq"
)

func (db *DBX) CreateStageTx(ctx context.Context, text, userid string, storyID int) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CreateStageTx failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check user exist
	queryUser := "SELECT id FROM users WHERE userid = $1" // dev:finder+query
	stmtQueryUser, err := db.Conn.PrepareContext(ctx, queryUser)
	if err != nil {
		db.Logger.Error("tx prepare on queryUser failed", "error", err.Error())
		return -1, err
	}
	defer stmtQueryUser.Close()
	var userID int
	err = tx.StmtContext(ctx, stmtQueryUser).QueryRow(userid).Scan(&userID)
	if err != nil {
		db.Logger.Debug("user not found", "return", err.Error())

	}
	// insert user if it does not exist
	if userID == 0 {
		id, err := db.createUser(ctx, userid, tx)
		if err != nil {
			db.Logger.Error("insert into users failed", "error", err.Error())
			return -1, err
		}
		userID = id
	}

	queryStoryKeys := "SELECT k.encoding_key FROM story AS s JOIN story_keys AS k ON s.id = k.story_id WHERE s.id = $1" // dev:finder+query
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

	// create stage
	queryInsertStage := "INSERT INTO stage(display_text, encoding_key, finished, storyteller_id, story_id) VALUES($1, $2, $3, $4, $5) RETURNING id" // dev:finder+query
	stmtInsertStage, err := db.Conn.PrepareContext(ctx, queryInsertStage)
	if err != nil {
		db.Logger.Error("tx prepare on stmtInsertStage failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsertStage.Close()
	var stageID int
	err = tx.StmtContext(ctx, stmtInsertStage).QueryRow(text, encodingKey, false, userID, storyID).Scan(&stageID)
	if err != nil {
		db.Logger.Error("query row insert into stage failed", "error", err.Error())
		return -1, err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateStageTx failed", "error", err.Error())
		return -1, err
	}

	return stageID, nil
}

func (db *DBX) AddChannelToStage(ctx context.Context, channel string, id int) (int, error) {
	query := "INSERT INTO stage_channel(channel, upstream_id, active) VALUES($1, $2, $3) RETURNING id" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(channel, id, true).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into stage_channel failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) AddEncounterToStage(ctx context.Context, text string, stage_id, storyteller_id, encounter_id int) (int, error) {
	query := "INSERT INTO stage_encounters(display_text, stage_id, storyteller_id, encounter_id) VALUES($1, $2, $3, $4) RETURNING id" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into stage_encounters failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(text, stage_id, storyteller_id, encounter_id).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into stage_encounters failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) UpdatePhase(ctx context.Context, id, phase int) error {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on UpdatePhase failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check if phase is valid
	if phase == int(types.Running) {
		// check if on this stage we have any stage_encounter already in running stage
		query := "SELECT (s.phase >= 0) FROM stage_encounters AS s JOIN stage_encounters AS se ON se.stage_id = s.stage_id WHERE se.id = $1 AND s.phase = $2" // dev:finder+query
		var running bool
		p := int(types.Running)
		if err = tx.QueryRowContext(ctx, query, id, p).Scan(&running); err != nil {
			if err != sql.ErrNoRows {
				db.Logger.Error("no rows passed", "err", err.Error())
				return err
			}
		}
		if running {
			return fmt.Errorf("stage_id already have a stage_encounters in running state")
		}
	}

	// change state of stage_encounter
	queryUpdateStageEncounter := "UPDATE stage_encounters SET phase = $1, updated_at = NOW() WHERE id = $2 RETURNING id" // dev:finder+query
	stmtUpdateStageEncounter, err := db.Conn.PrepareContext(ctx, queryUpdateStageEncounter)
	if err != nil {
		db.Logger.Error("tx prepare on stmtUpdateStageEncounter failed", "error", err.Error())
		return err
	}
	defer stmtUpdateStageEncounter.Close()
	var ID int
	err = tx.StmtContext(ctx, stmtUpdateStageEncounter).QueryRow(phase, id).Scan(&ID)
	if err != nil {
		db.Logger.Error("query row insert into StageEncounter failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on UpdatePhase failed", "error", err.Error())
		return err
	}
	db.Logger.Debug("stage_encounter changed", "id", id)
	return nil
}

// AddParticipants func stage_encounters_participants_players
func (db *DBX) AddParticipants(ctx context.Context, encounterID int, npc bool, players []int) error {

	query := "INSERT INTO stage_encounters_participants_players (players_id, stage_encounter_id) VALUES ($1, $2)" // dev:finder+query
	if npc {
		query = "INSERT INTO stage_encounters_participants_non_players (non_players_id, stage_encounter_id) VALUES ($1, $2)" // dev:finder+query
	}
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
func (db *DBX) AddNextEncounter(ctx context.Context, next []types.Next) error {

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
	for _, n := range next {
		query := "INSERT INTO stage_next_encounter (display_text, upstream_id, current_encounter_id, next_encounter_id) VALUES ($1, $2, $3, $4) RETURNING id" // dev:finder+query
		stmt, err := db.Conn.PrepareContext(ctx, query)
		if err != nil {
			db.Logger.Error("prepare insert into stage_next_encounter failed", "error", err.Error())
			return err
		}
		defer stmt.Close()
		var nextEncounterIDDB int
		err = tx.StmtContext(ctx, stmt).QueryRow(n.Text, n.UpstreamID, n.EncounterID, n.NextEncounterID).Scan(&nextEncounterIDDB)
		if err != nil {
			db.Logger.Error("error on insert into stage_next_encounter", "error", err.Error())
			return err
		}
		// insert into stage_next_objectives
		queryObjectives := "INSERT INTO stage_next_objectives (upstream_id, kind, values) VALUES ($1, $2, $3) RETURNING id" // dev:finder+query
		stmtObjectives, err := db.Conn.PrepareContext(ctx, queryObjectives)
		if err != nil {
			db.Logger.Error("prepare insert into stage_next_objectives failed", "error", err.Error())
			return err
		}
		defer stmtObjectives.Close()
		var objectiveID int
		err = tx.StmtContext(ctx, stmtObjectives).QueryRow(nextEncounterIDDB, n.Objective.Kind, pq.Array(n.Objective.Values)).Scan(&objectiveID)
		if err != nil {
			db.Logger.Error("error on insert into stage_next_objectives", "error", err.Error())
			return err
		}
		db.Logger.Debug("stage next objective added", "next_id", nextEncounterIDDB, "objective_id", objectiveID)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit stage_next_encounter", "error", err.Error())
		return err
	}

	return nil
}

// stage_running_tasks
func (db *DBX) AddRunningTask(ctx context.Context, text string, stageID, taskID, StorytellerID, encounterID int) error {
	query := "INSERT INTO stage_running_tasks (display_text, stage_id, task_id, storyteller_id, stage_encounter_id) VALUES ($1, $2, $3, $4, $5)" // dev:finder+query
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
func (db *DBX) RegisterActivities(ctx context.Context, stageID, encounterID int, actions types.Actions) error {
	db.Logger.Debug("RegisterActivities", "stageID", stageID, "encounterID", encounterID, "actions", actions)
	query := "INSERT INTO stage_encounter_activities (actions, upstream_id, encounter_id) VALUES ($1, $2, $3)" // dev:finder+query
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
	query := "UPDATE stage_encounter_activities SET processed = $1, actions = $2 WHERE id = $3" // dev:finder+query
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

func (db *DBX) CloseStage(ctx context.Context, id int) error {
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
	// verify if all encounters are finished
	// types.Running = 3
	query := "SELECT (se.phase = 3), se.id, sc.channel, s.display_text FROM stage_encounters AS se JOIN stage_channel AS sc ON sc.upstream_id = se.stage_id JOIN stage AS s ON s.id = se.stage_id WHERE se.stage_id = $1" // dev:finder+query
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	var lastEncounterID int
	var channel, displayText string
	for rows.Next() {
		var finished bool
		if err := rows.Scan(&finished, &lastEncounterID, &channel, &displayText); err != nil {
			return err
		}
		if !finished {
			return fmt.Errorf("not all encounters are finished")
		}
	}

	queryUpdate := "UPDATE stage SET finished = $1 WHERE id = $2"
	_, err = tx.ExecContext(ctx, queryUpdate, true, id)
	if err != nil {
		return err
	}

	// insert stage_encounter_activities
	queryInsert := "INSERT INTO stage_encounter_activities (actions, upstream_id, encounter_id) VALUES ($1, $2, $3)" // dev:finder+query
	actions := types.NewActions()
	actions["channel"] = channel
	actions["text"] = "Stage is finished"
	actions["command"] = "close-stage"
	actions["display_text"] = displayText
	_, err = tx.ExecContext(ctx, queryInsert, actions, id, lastEncounterID)
	if err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (db *DBX) DeleteStageNextEncounter(ctx context.Context, id int) error {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on DeleteStagegNextEncounter failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// SELECT ids
	query := "SELECT s.id, o.id FROM stage_next_encounter AS s JOIN stage_next_objectives AS o ON s.id = o.upstream_id WHERE s.id = $1" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stage_next_encounter failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	var nextID, objectiveID int
	err = tx.StmtContext(ctx, stmt).QueryRow(id).Scan(&nextID, &objectiveID)
	if err != nil {
		db.Logger.Error("tx query on stage_next_encounter failed", "error", err.Error())
		return err
	}
	// delete FROM stage_next_objectives
	query = "DELETE FROM stage_next_objectives WHERE id = $1" // dev:finder+query
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stage_next_objectives failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, objectiveID)
	if err != nil {
		db.Logger.Error("tx exec on stage_next_objectives failed", "error", err.Error())
		return err
	}
	// delete FROM stage_next_encounter
	query = "DELETE FROM stage_next_encounter WHERE id = $1" // dev:finder+query
	stmt, err = db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stage_next_encounter failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, nextID)
	if err != nil {
		db.Logger.Error("tx exec on stage_next_encounter failed", "error", err.Error())
		return err
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit stage_next_encounter", "error", err.Error())
		return err
	}
	return nil
}

func (db *DBX) DeleteStageEncounterByID(ctx context.Context, id int) error {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on DeleteStageEncounterByID failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// delete FROM stage_encounters
	query := "DELETE FROM stage_encounters WHERE id = $1" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on stage_encounters failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, id)
	if err != nil {
		db.Logger.Error("tx exec on stage_encounters failed", "error", err.Error())
		return err
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("error on commit stage_encounters", "error", err.Error())
		return err
	}
	return nil
}
