package pg

import (
	"context"
	"fmt"

	"github.com/betorvs/playbypost/core/initiative"
)

func (db *DBX) UpdateNextPlayer(ctx context.Context, id, nextPlayer int) error {
	query := "Update initiative SET next_player = $1 WHERE id = $2 RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on initiative failed", "error", err.Error())
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(nextPlayer, id).Scan(&res)
	if err != nil {
		db.Logger.Error("tx statement on initiative failed", "error", err.Error())
	}
	return nil
}

func (db *DBX) SaveInitiativeTx(ctx context.Context, i initiative.Initiative, encounterID int) (int, error) {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on initiative failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()

	query := "INSERT INTO initiative(title, stage_encounters_id, next_player) VALUES($1, $2, $3) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on initiative failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	// ExecContext(ctx, query, "test01", 1)
	var id int
	err = tx.StmtContext(ctx, stmt).QueryRow(i.Name, encounterID, i.Position).Scan(&id)
	if err != nil {
		db.Logger.Error("tx statement on initiative failed", "error", err.Error())
		return -1, err
	}

	queryParticipants := "INSERT INTO initiative_participants(initiative_id, participant_name, participant_bonus, participant_result, active) VALUES($1, $2, $3, $4, $5) RETURNING id"

	for _, p := range i.Participants {
		_, err = tx.ExecContext(ctx, queryParticipants, id, p.Name, p.Bonus, p.Result, true)
		if err != nil {
			db.Logger.Error("tx participants on initiative failed", "error", err.Error())
			return -1, err
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on initiative failed", "error", err.Error())
		return -1, err
	}
	return id, nil
}

func (db *DBX) GetInitiativeByID(ctx context.Context, id int) (initiative.Initiative, error) {
	obj := initiative.Initiative{}
	rows, err := db.Conn.QueryContext(ctx, "SELECT i.id, i.title, i.next_player, p.participant_name, p.participant_result FROM initiative AS i JOIN initiative_participants AS p ON i.id = p.initiative_id WHERE p.active = true AND i.id = $1", id)
	if err != nil {
		db.Logger.Error("query on users failed", "error", err.Error())
		return obj, err
	}
	defer rows.Close()
	var nextPlayer, result int
	var title, name string
	party := initiative.Participants{}
	for rows.Next() {
		if err := rows.Scan(&id, &title, &nextPlayer, &name, &result); err != nil {
			db.Logger.Error("scan on users error", "error", err.Error())
		}
		p := initiative.Participant{}
		p.Name = name
		p.Result = result
		party = append(party, p)
	}
	obj.Name = title
	obj.Position = nextPlayer
	obj.Participants = party
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows on users error", "error", err.Error())
	}

	return obj, nil
}

func (db *DBX) GetRunningInitiativeByEncounterID(ctx context.Context, encounterID int) (initiative.Initiative, int, error) {
	initiativeID := -1
	obj := initiative.Initiative{}
	rows, err := db.Conn.QueryContext(ctx, "SELECT i.id, i.title, i.next_player, p.participant_name, p.participant_result FROM initiative AS i JOIN initiative_participants AS p ON i.id = p.initiative_id JOIN stage_encounters AS se ON se.id = i.stage_encounters_id WHERE p.active = TRUE AND se.phase = 2 AND se.id = $1", encounterID)
	if err != nil {
		db.Logger.Error("query on users failed", "error", err.Error())
		return obj, initiativeID, err
	}
	defer rows.Close()
	var nextPlayer, result int
	var title, name string
	party := initiative.Participants{}
	for rows.Next() {
		if err := rows.Scan(&initiativeID, &title, &nextPlayer, &name, &result); err != nil {
			db.Logger.Error("scan on users error", "error", err.Error())
		}
		p := initiative.Participant{}
		p.Name = name
		p.Result = result
		party = append(party, p)
	}
	obj.Name = title
	obj.Position = nextPlayer
	obj.Participants = party
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows on users error", "error", err.Error())
	}

	return obj, initiativeID, nil
}

// deactivate part
func (db *DBX) DeactivateParticipant(ctx context.Context, id int, name string) (int, error) {
	query := "Update initiative_participants SET active = FALSE WHERE initiative_id = $1 AND participant_name = $2 RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("tx prepare on initiative_participants failed", "error", err.Error())
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(id, name).Scan(&res)
	if err != nil {
		db.Logger.Error("tx statement on initiative_participants failed", "error", err.Error())
	}
	return res, nil
}
